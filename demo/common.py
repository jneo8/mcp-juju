"""Shared pieces of the mcp-juju demo scripts.

The demo mirrors the functional tests: the MCP server does the work through
its tools, and Jubilant (a thin wrapper over the `juju` CLI) observes the
result. Everything lives in one fixed model so the three scripts can be run
independently, in any order, and repeatedly.
"""

from __future__ import annotations

import argparse
import contextlib
import json
import os
import pathlib
import subprocess
import sys
import time
from collections.abc import Callable, Generator
from dataclasses import dataclass

import jubilant
from mcp.types import CallToolResult
from mcpclient import McpJujuClient, result_text

REPO_ROOT = pathlib.Path(__file__).resolve().parents[1]
BUILD_DIR = pathlib.Path(__file__).resolve().parent / '.build'

DEFAULT_MODEL = os.environ.get('MCP_JUJU_DEMO_MODEL', 'mcp-juju-demo')
DEFAULT_CHARM = os.environ.get('MCP_JUJU_DEMO_CHARM', 'postgresql')
APP = os.environ.get('MCP_JUJU_DEMO_APP', 'db')

DEPLOY_TIMEOUT = 20 * 60
DESTROY_TIMEOUT = 10 * 60


@dataclass(frozen=True)
class Target:
    """Where the demo lives: a model, optionally on a named controller."""

    model: str
    controller: str | None

    @property
    def qualified(self) -> str:
        """The `[controller:]model` form accepted by `--model` and by Jubilant."""
        return f'{self.controller}:{self.model}' if self.controller else self.model


def parser(description: str) -> argparse.ArgumentParser:
    p = argparse.ArgumentParser(description=description)
    p.add_argument(
        '--model',
        default=DEFAULT_MODEL,
        help=f'demo model name (default: {DEFAULT_MODEL}, or MCP_JUJU_DEMO_MODEL)',
    )
    p.add_argument(
        '--controller',
        default=os.environ.get('MCP_JUJU_DEMO_CONTROLLER'),
        help='controller to use (default: the current controller, or MCP_JUJU_DEMO_CONTROLLER)',
    )
    p.add_argument(
        '--binary',
        default=os.environ.get('MCP_JUJU_BINARY'),
        help='prebuilt mcp-juju binary (default: build from source, or MCP_JUJU_BINARY)',
    )
    return p


def target_from(args: argparse.Namespace) -> Target:
    return Target(model=args.model, controller=args.controller)


def juju_for(target: Target) -> jubilant.Juju:
    """A Jubilant handle bound to the demo model (it may not exist yet)."""
    return jubilant.Juju(model=target.qualified)


def model_exists(target: Target) -> bool:
    juju = jubilant.Juju()
    try:
        juju.show_model(target.qualified)
    except jubilant.CLIError:
        return False
    return True


def build_binary(override: str | None) -> pathlib.Path:
    """Return the mcp-juju binary to demo, building it from the repository if needed."""
    if override:
        path = pathlib.Path(override).resolve()
        if not path.exists():
            sys.exit(f'mcp-juju binary not found: {path}')
        return path
    BUILD_DIR.mkdir(exist_ok=True)
    path = BUILD_DIR / 'mcp-juju'
    step(f'building {path.relative_to(REPO_ROOT)}')
    subprocess.run(['go', 'build', '-o', str(path), '.'], cwd=REPO_ROOT, check=True)
    return path


@contextlib.contextmanager
def mcp_server(binary: pathlib.Path, *extra_args: str) -> Generator[McpJujuClient]:
    """Start the mcp-juju binary over stdio and yield a synchronous client."""
    args = ['--server-type', 'stdio', *extra_args]
    with McpJujuClient.stdio(str(binary), args=args) as client:
        yield client


def call(client: McpJujuClient, tool: str, arguments: dict | None = None) -> CallToolResult:
    """Call a tool, echo the call, and exit on an isError result."""
    step(f'mcp: {tool} {arguments or {}}')
    result = client.call_tool(tool, arguments)
    if result.is_error:
        sys.exit(f'tool {tool} failed:\n{result_text(result)}')
    return result


# --- Output helpers ---------------------------------------------------------

_START = time.monotonic()


def step(message: str) -> None:
    elapsed = time.monotonic() - _START
    print(f'[{elapsed:6.1f}s] {message}', flush=True)


def banner(title: str) -> None:
    print(f'\n=== {title} ===', flush=True)


# Set by verify.py --full: show tool output in full instead of an excerpt.
FULL_OUTPUT = False


def show(label: str, text: str, *, limit: int = 12) -> None:
    """Print a tool's output indented under the current check.

    Long output is cut after *limit* lines unless FULL_OUTPUT is set.
    """
    lines = text.rstrip().splitlines() or ['(empty)']
    if not FULL_OUTPUT and len(lines) > limit:
        hidden = len(lines) - limit
        lines = lines[:limit] + [f'... ({hidden} more lines, use --full to see everything)']
    print(f'        {label}:')
    for line in lines:
        print(f'          {line}')
    print(flush=True)


def show_json(label: str, data: object, *, limit: int = 12) -> None:
    """Print JSON data (for example a tool's structuredContent) with show()."""
    show(label, json.dumps(data, indent=2, sort_keys=True), limit=limit)


class Checks:
    """Collects verification results and prints a summary."""

    def __init__(self) -> None:
        self.passed: list[str] = []
        self.failed: list[tuple[str, str]] = []

    def run(self, name: str, fn: Callable[[], None]) -> None:
        """Run one check: print its name, the output it shows, then the verdict."""
        print(f'  * {name}', flush=True)
        try:
            fn()
        except Exception as e:  # noqa: BLE001 - reported in the summary
            self.failed.append((name, f'{type(e).__name__}: {e}'))
            print(f'    FAIL  {type(e).__name__}: {e}\n', flush=True)
        else:
            self.passed.append(name)
            print('    ok\n', flush=True)

    def summary(self) -> int:
        banner('Summary')
        print(f'{len(self.passed)} passed, {len(self.failed)} failed')
        for name, reason in self.failed:
            print(f'  FAIL  {name}: {reason}')
        return 1 if self.failed else 0
