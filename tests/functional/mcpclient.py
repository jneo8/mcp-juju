"""Synchronous wrapper around the MCP Python client for use in pytest.

Jubilant and pytest-jubilant are synchronous, so the tests are too. The MCP
SDK is async; this module runs it on a private event loop thread and keeps the
whole client lifecycle inside one task, which anyio requires for its cancel
scopes.
"""

from __future__ import annotations

import asyncio
import concurrent.futures
import contextlib
import dataclasses
import secrets
import socket
import subprocess
import threading
import time
from collections.abc import Awaitable, Callable, Generator
from typing import Any

import mcp
import mcp.types
from mcp.client.stdio import StdioServerParameters
from mcp.client.streamable_http import create_mcp_http_client, streamable_http_client

_STOP = object()


class McpJujuClient:
    """Drive an mcp-juju server from synchronous code.

    ``server`` is anything ``mcp.Client`` accepts: ``StdioServerParameters``
    (spawn the binary and talk over stdio), a URL string, or a transport.
    Use as a context manager. Each public method blocks until the server has
    replied.
    """

    def __init__(self, server: Any, *, call_timeout: float = 600.0):
        self._server = server
        self._call_timeout = call_timeout
        self._loop = asyncio.new_event_loop()
        self._thread = threading.Thread(target=self._loop.run_forever, daemon=True)
        self._queue: asyncio.Queue[Any] | None = None
        self._ready: concurrent.futures.Future[None] = concurrent.futures.Future()
        self._done: concurrent.futures.Future[None] | None = None
        self.instructions: str | None = None

    @classmethod
    def stdio(
        cls,
        binary: str,
        *,
        args: list[str] | None = None,
        env: dict[str, str] | None = None,
        call_timeout: float = 600.0,
    ) -> McpJujuClient:
        """Spawn ``binary`` in stdio mode and talk to it over its pipes."""
        params = StdioServerParameters(command=binary, args=args or [], env=env)
        return cls(params, call_timeout=call_timeout)

    @classmethod
    def http(
        cls, url: str, *, token: str | None = None, call_timeout: float = 600.0
    ) -> McpJujuClient:
        """Connect to a running Streamable HTTP server, optionally with a bearer token."""
        headers = {'Authorization': f'Bearer {token}'} if token else None
        transport = streamable_http_client(url, http_client=create_mcp_http_client(headers=headers))
        return cls(transport, call_timeout=call_timeout)

    def __enter__(self) -> McpJujuClient:
        self._thread.start()
        self._done = asyncio.run_coroutine_threadsafe(self._main(), self._loop)
        self._ready.result(timeout=60)
        return self

    def __exit__(self, *exc: object) -> None:
        if self._queue is not None:
            self._loop.call_soon_threadsafe(self._queue.put_nowait, _STOP)
        if self._done is not None:
            self._done.result(timeout=30)
        self._loop.call_soon_threadsafe(self._loop.stop)
        self._thread.join(timeout=10)
        self._loop.close()

    async def _main(self) -> None:
        self._queue = asyncio.Queue()
        try:
            async with mcp.Client(self._server) as client:
                self.instructions = client.instructions
                self._ready.set_result(None)
                while True:
                    item = await self._queue.get()
                    if item is _STOP:
                        return
                    future, fn = item
                    try:
                        future.set_result(await fn(client))
                    except BaseException as e:  # noqa: BLE001 - forwarded to the caller
                        future.set_exception(e)
        except BaseException as e:
            if not self._ready.done():
                self._ready.set_exception(e)
            raise

    def _submit(self, fn: Callable[[mcp.Client], Awaitable[Any]]) -> Any:
        assert self._queue is not None, 'client not started; use it as a context manager'
        future: concurrent.futures.Future[Any] = concurrent.futures.Future()
        self._loop.call_soon_threadsafe(self._queue.put_nowait, (future, fn))
        return future.result(timeout=self._call_timeout)

    def list_tools(self) -> list[mcp.types.Tool]:
        return self._submit(lambda c: c.list_tools()).tools

    def list_resources(self) -> list[mcp.types.Resource]:
        return self._submit(lambda c: c.list_resources()).resources

    def list_resource_templates(self) -> list[mcp.types.ResourceTemplate]:
        return self._submit(lambda c: c.list_resource_templates()).resource_templates

    def call_tool(self, name: str, arguments: dict[str, Any] | None = None) -> mcp.types.CallToolResult:
        return self._submit(lambda c: c.call_tool(name, arguments or {}))

    def read_resource(self, uri: str) -> mcp.types.ReadResourceResult:
        return self._submit(lambda c: c.read_resource(uri))


def result_text(result: mcp.types.CallToolResult) -> str:
    """Return the concatenated text blocks of a tool result."""
    return '\n'.join(
        block.text for block in result.content if isinstance(block, mcp.types.TextContent)
    )


@dataclasses.dataclass(frozen=True)
class HttpServer:
    """A running mcp-juju Streamable HTTP server."""

    url: str
    token: str
    process: subprocess.Popen[bytes]


def free_port() -> int:
    """Ask the kernel for an unused loopback port."""
    with socket.socket() as s:
        s.bind(('127.0.0.1', 0))
        return s.getsockname()[1]


@contextlib.contextmanager
def http_server(
    binary: str, *, extra_args: list[str] | None = None, timeout: float = 30.0
) -> Generator[HttpServer]:
    """Run ``binary`` in HTTP mode on a free loopback port with a random bearer token."""
    port = free_port()
    token = secrets.token_hex(16)
    args = [
        binary,
        '--server-type',
        'http',
        '--host',
        '127.0.0.1',
        '--port',
        str(port),
        '--auth-token',
        token,
        *(extra_args or []),
    ]
    process = subprocess.Popen(args, stdout=subprocess.DEVNULL, stderr=subprocess.PIPE)
    try:
        deadline = time.monotonic() + timeout
        while True:
            with contextlib.suppress(OSError), socket.create_connection(('127.0.0.1', port), 0.5):
                break
            if process.poll() is not None:
                stderr = process.stderr.read().decode() if process.stderr else ''
                raise RuntimeError(f'mcp-juju exited with {process.returncode}: {stderr}')
            if time.monotonic() > deadline:
                raise TimeoutError(f'mcp-juju did not listen on port {port} within {timeout}s')
            time.sleep(0.1)
        yield HttpServer(url=f'http://127.0.0.1:{port}/mcp', token=token, process=process)
    finally:
        process.terminate()
        with contextlib.suppress(subprocess.TimeoutExpired):
            process.wait(timeout=10)
        if process.poll() is None:
            process.kill()
