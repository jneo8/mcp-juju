"""Synchronous wrapper around the MCP Python client for use in pytest.

Jubilant and pytest-jubilant are synchronous, so the tests are too. The MCP
SDK is async; this module runs it on a private event loop thread and keeps the
whole client lifecycle inside one task, which anyio requires for its cancel
scopes.
"""

from __future__ import annotations

import asyncio
import concurrent.futures
import threading
from collections.abc import Awaitable, Callable
from typing import Any

import mcp
import mcp.types
from mcp.client.stdio import StdioServerParameters

_STOP = object()


class McpJujuClient:
    """Drive an mcp-juju server process over stdio from synchronous code.

    Use as a context manager. Each public method blocks until the server has
    replied.
    """

    def __init__(
        self,
        binary: str,
        *,
        args: list[str] | None = None,
        env: dict[str, str] | None = None,
        call_timeout: float = 600.0,
    ):
        self._params = StdioServerParameters(command=binary, args=args or [], env=env)
        self._call_timeout = call_timeout
        self._loop = asyncio.new_event_loop()
        self._thread = threading.Thread(target=self._loop.run_forever, daemon=True)
        self._queue: asyncio.Queue[Any] | None = None
        self._ready: concurrent.futures.Future[None] = concurrent.futures.Future()
        self._done: concurrent.futures.Future[None] | None = None
        self.instructions: str | None = None

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
            async with mcp.Client(self._params) as client:
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
