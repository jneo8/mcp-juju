"""Shared fixtures for the mcp-juju functional tests.

The tests need a bootstrapped Juju controller; pytest-jubilant creates a
temporary model on it for each test module (see its README for the
--juju-model, --juju-controller and --no-juju-teardown options).
"""

from __future__ import annotations

import os
import pathlib
import subprocess
from collections.abc import Generator

import pytest

from mcpclient import HttpServer, McpJujuClient, http_server

REPO_ROOT = pathlib.Path(__file__).resolve().parents[2]


def pytest_addoption(parser: pytest.Parser):
    parser.addoption(
        '--test-charm',
        default=os.environ.get('MCP_JUJU_TEST_CHARM', 'postgresql'),
        help='charm to deploy in the functional tests (default: postgresql, a machine charm)',
    )


@pytest.fixture(scope='session')
def test_charm(request: pytest.FixtureRequest) -> str:
    return str(request.config.getoption('--test-charm'))


@pytest.fixture(scope='session')
def mcp_juju_binary(tmp_path_factory: pytest.TempPathFactory) -> pathlib.Path:
    """Build mcp-juju from the repository, unless MCP_JUJU_BINARY points at one."""
    override = os.environ.get('MCP_JUJU_BINARY')
    if override:
        path = pathlib.Path(override).resolve()
        assert path.exists(), f'MCP_JUJU_BINARY not found: {path}'
        return path
    path = tmp_path_factory.mktemp('bin') / 'mcp-juju'
    subprocess.run(['go', 'build', '-o', str(path), '.'], cwd=REPO_ROOT, check=True)
    return path


@pytest.fixture(scope='module')
def mcp(mcp_juju_binary: pathlib.Path) -> Generator[McpJujuClient]:
    """Module-scoped mcp-juju server exposing every tool, talking over stdio.

    The server inherits the environment, so it uses the same Juju client store
    as the `juju` CLI that Jubilant drives.
    """
    with McpJujuClient.stdio(str(mcp_juju_binary), args=['--server-type', 'stdio']) as client:
        yield client


@pytest.fixture(scope='module')
def mcp_http_server(mcp_juju_binary: pathlib.Path) -> Generator[HttpServer]:
    """Module-scoped mcp-juju Streamable HTTP server on loopback with a bearer token."""
    with http_server(str(mcp_juju_binary)) as srv:
        yield srv


@pytest.fixture(scope='module')
def mcp_http(mcp_http_server: HttpServer) -> Generator[McpJujuClient]:
    """Module-scoped MCP client connected to the HTTP server with the right token."""
    with McpJujuClient.http(mcp_http_server.url, token=mcp_http_server.token) as client:
        yield client
