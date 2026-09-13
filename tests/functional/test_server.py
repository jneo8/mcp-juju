"""Functional tests for the mcp-juju server.

One temporary model is created for this module and one application is
deployed into it through the MCP `deploy` tool; every test then verifies a
different tool or resource against that deployment, cross-checking with
Jubilant. The default charm is `postgresql` (a machine charm, so a machine
cloud such as LXD is expected); pass --test-charm to use another one.
"""

from __future__ import annotations

import contextlib
import json
import pathlib
import subprocess

import httpx2
import jubilant
import pytest
from mcp.types import TextContent, TextResourceContents

from mcpclient import HttpServer, McpJujuClient, result_text

APP = 'mcp-app'
DEPLOY_TIMEOUT = 20 * 60


@pytest.fixture(scope='module', autouse=True)
def deployed(juju: jubilant.Juju, mcp: McpJujuClient, test_charm: str) -> None:
    """Deploy the test charm through the MCP `deploy` tool, once per module.

    Skips the deploy when the application already exists, so the module can be
    rerun against a kept model with --no-juju-setup.
    """
    if APP in juju.status().apps:
        return
    result = mcp.call_tool('deploy', {'model': juju.model, 'args': [test_charm, APP]})
    assert not result.is_error, result_text(result)
    juju.wait(lambda status: jubilant.all_active(status, APP), timeout=DEPLOY_TIMEOUT)


# --- Server metadata -------------------------------------------------------


def test_instructions_and_tool_list(mcp: McpJujuClient):
    assert mcp.instructions and 'juju' in mcp.instructions.lower()

    tools = {tool.name: tool for tool in mcp.list_tools()}
    assert {'status', 'deploy', 'integrate', 'destroy-model'} <= tools.keys()

    status = tools['status']
    assert status.annotations is not None
    assert status.annotations.read_only_hint is True
    assert status.annotations.destructive_hint is False
    assert status.input_schema['properties']['format']['default'] == 'json'
    assert 'json' in status.input_schema['properties']['format']['enum']

    destroy = tools['destroy-model']
    assert destroy.annotations is not None
    assert destroy.annotations.destructive_hint is True


def test_doc_resource(mcp: McpJujuClient):
    resources = {r.name for r in mcp.list_resources()}
    assert 'status-doc' in resources

    result = mcp.read_resource('juju://status-doc')
    contents = result.contents[0]
    assert isinstance(contents, TextResourceContents)
    assert contents.text.startswith('# status')
    assert 'Details' in contents.text


# --- Read-only tools -------------------------------------------------------


@pytest.mark.juju_setup
def test_deploy(juju: jubilant.Juju, test_charm: str):
    # The autouse `deployed` fixture ran the MCP deploy; this checks the outcome.
    app = juju.status().apps[APP]
    assert app.is_active
    assert test_charm in app.charm


def test_status_matches_jubilant(juju: jubilant.Juju, mcp: McpJujuClient):
    result = mcp.call_tool('status', {'model': juju.model})
    assert not result.is_error, result_text(result)

    # Structured content mirrors `juju status --format json`.
    structured = result.structured_content
    assert structured is not None
    assert structured['model']['name'] == juju.model
    assert set(structured['applications']) == set(juju.status().apps) == {APP}

    # The first text block is exactly the JSON, for clients without structured
    # support; Juju's stderr notes (if any) follow in a separate block.
    first = result.content[0]
    assert isinstance(first, TextContent)
    assert json.loads(first.text) == structured


def test_models_lists_test_model(juju: jubilant.Juju, mcp: McpJujuClient):
    result = mcp.call_tool('models')
    assert not result.is_error, result_text(result)
    assert result.structured_content is not None
    names = {m['short-name'] for m in result.structured_content['models']}
    assert juju.model in names


def test_show_application(juju: jubilant.Juju, mcp: McpJujuClient):
    result = mcp.call_tool('show-application', {'model': juju.model, 'args': [APP]})
    assert not result.is_error, result_text(result)
    assert result.structured_content is not None
    assert APP in result.structured_content


def test_config_tool(juju: jubilant.Juju, mcp: McpJujuClient):
    result = mcp.call_tool('config', {'model': juju.model, 'args': [APP]})
    assert not result.is_error, result_text(result)
    structured = result.structured_content
    assert structured is not None
    assert structured['application'] == APP
    # Compare with the raw CLI output: jubilant's juju.config() drops unset
    # options, whereas the tool returns the full `juju config --format json`.
    expected = json.loads(juju.cli('config', APP, '--format', 'json'))
    assert set(structured['settings']) == set(expected['settings'])


def test_config_resource_template(juju: jubilant.Juju, mcp: McpJujuClient):
    templates = {t.name for t in mcp.list_resource_templates()}
    assert 'Juju Application Configuration' in templates

    # The template has no model parameter, so it reads the current model.
    # Switch to the test model and switch back afterwards; the previous
    # current model may no longer exist, so restoring it is best effort.
    previous = juju.cli('switch', include_model=False).strip()
    juju.cli('switch', juju.model, include_model=False)
    try:
        resource = mcp.read_resource(f'juju://config/{APP}')
    finally:
        with contextlib.suppress(jubilant.CLIError):
            juju.cli('switch', previous, include_model=False)
    contents = resource.contents[0]
    assert isinstance(contents, TextResourceContents)
    assert json.loads(contents.text)['application'] == APP


def test_unknown_model_is_tool_error(mcp: McpJujuClient):
    result = mcp.call_tool('status', {'model': 'mcp-juju-no-such-model'})
    assert result.is_error
    assert 'mcp-juju-no-such-model' in result_text(result)


def test_unknown_application_is_tool_error(juju: jubilant.Juju, mcp: McpJujuClient):
    result = mcp.call_tool('show-application', {'model': juju.model, 'args': ['no-such-app']})
    assert result.is_error
    assert 'no-such-app' in result_text(result)


# --- Streamable HTTP transport -----------------------------------------------


def test_http_requires_bearer_token(mcp_http_server: HttpServer):
    initialize = {
        'jsonrpc': '2.0',
        'id': 1,
        'method': 'initialize',
        'params': {
            'protocolVersion': '2025-06-18',
            'capabilities': {},
            'clientInfo': {'name': 'test', 'version': '0'},
        },
    }
    headers = {'Accept': 'application/json, text/event-stream'}
    with httpx2.Client(headers=headers) as http:
        missing = http.post(mcp_http_server.url, json=initialize)
        assert missing.status_code == 401
        assert 'Bearer' in missing.headers.get('WWW-Authenticate', '')

        wrong = http.post(
            mcp_http_server.url, json=initialize, headers={'Authorization': 'Bearer nope'}
        )
        assert wrong.status_code == 401

        ok = http.post(
            mcp_http_server.url,
            json=initialize,
            headers={'Authorization': f'Bearer {mcp_http_server.token}'},
        )
        assert ok.status_code == 200, ok.text


def test_http_matches_stdio(juju: jubilant.Juju, mcp: McpJujuClient, mcp_http: McpJujuClient):
    assert {t.name for t in mcp_http.list_tools()} == {t.name for t in mcp.list_tools()}

    result = mcp_http.call_tool('status', {'model': juju.model})
    assert not result.is_error, result_text(result)
    assert result.structured_content is not None
    assert result.structured_content['model']['name'] == juju.model
    assert set(result.structured_content['applications']) == {APP}


def test_http_refuses_non_loopback_without_token(mcp_juju_binary: pathlib.Path):
    proc = subprocess.run(
        [str(mcp_juju_binary), '--server-type', 'http', '--host', '0.0.0.0', '--port', '0'],
        capture_output=True,
        text=True,
        timeout=30,
    )
    assert proc.returncode != 0
    assert 'auth-token' in proc.stderr + proc.stdout


# --- Mutating tools --------------------------------------------------------


def test_add_unit(juju: jubilant.Juju, mcp: McpJujuClient):
    result = mcp.call_tool('add-unit', {'model': juju.model, 'args': [APP], 'num-units': 1})
    assert not result.is_error, result_text(result)
    juju.wait(
        lambda status: jubilant.all_active(status, APP) and len(status.apps[APP].units) == 2,
        timeout=DEPLOY_TIMEOUT,
    )


@pytest.mark.juju_teardown
def test_remove_unit(juju: jubilant.Juju, mcp: McpJujuClient):
    unit = sorted(juju.status().apps[APP].units)[-1]
    result = mcp.call_tool('remove-unit', {'model': juju.model, 'args': [unit], 'no-prompt': True})
    assert not result.is_error, result_text(result)
    juju.wait(lambda status: len(status.apps[APP].units) == 1, timeout=10 * 60)


@pytest.mark.juju_teardown
def test_remove_application(juju: jubilant.Juju, mcp: McpJujuClient):
    result = mcp.call_tool(
        'remove-application', {'model': juju.model, 'args': [APP], 'no-prompt': True}
    )
    assert not result.is_error, result_text(result)
    juju.wait(lambda status: APP not in status.apps, timeout=10 * 60)
