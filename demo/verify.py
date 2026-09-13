#!/usr/bin/env python3
"""Verify the mcp-juju server against the demo environment.

Starts the server three ways (stdio with every tool, stdio with --read-only,
Streamable HTTP with a bearer token), shows what each read-only tool returns,
and checks the results against what Jubilant sees through the `juju` CLI.
Only read-only operations are performed, so the script can be rerun at any
time. Run deploy.py first.
"""

from __future__ import annotations

import json
import pathlib
import sys
import urllib.error
import urllib.request

import common
from common import APP, Checks, banner, show, show_json
from mcp.types import TextContent, TextResourceContents
from mcpclient import McpJujuClient, http_server, result_text


def ok(result):
    """Assert a tool call succeeded and return it."""
    assert not result.is_error, result_text(result)
    return result


def main() -> int:
    p = common.parser(__doc__)
    p.add_argument('--full', action='store_true', help='print tool output in full, not excerpts')
    args = p.parse_args()
    common.FULL_OUTPUT = args.full
    target = common.target_from(args)
    if not common.model_exists(target):
        sys.exit(f'model {target.qualified} not found; run `just demo-deploy` first')
    binary = common.build_binary(args.binary)
    juju = common.juju_for(target)
    model = target.qualified
    checks = Checks()
    stdio_tools: set[str] = set()

    banner('stdio server, every tool')
    with common.mcp_server(binary) as mcp:

        def instructions() -> None:
            assert mcp.instructions and 'juju' in mcp.instructions.lower(), mcp.instructions
            show('server instructions', mcp.instructions, limit=6)

        def tool_list() -> None:
            tools = mcp.list_tools()
            stdio_tools.update(t.name for t in tools)
            missing = {'status', 'deploy', 'integrate', 'destroy-model'} - stdio_tools
            assert not missing, f'missing tools: {missing}'
            assert len(stdio_tools) > 100, f'only {len(stdio_tools)} tools'
            names = sorted(stdio_tools)
            show(f'{len(names)} tools', ', '.join(names[:20]) + f', ... ({len(names) - 20} more)')

        def status_tool_definition() -> None:
            tool = next(t for t in mcp.list_tools() if t.name == 'status')
            assert tool.annotations and tool.annotations.read_only_hint is True
            assert tool.annotations.destructive_hint is False
            fmt = tool.input_schema['properties']['format']
            assert fmt['default'] == 'json' and 'json' in fmt['enum'], fmt
            show_json('status tool annotations', tool.annotations.model_dump(exclude_none=True))
            show_json('status tool "format" parameter', fmt)

        def doc_resource() -> None:
            contents = mcp.read_resource('juju://status-doc').contents[0]
            assert isinstance(contents, TextResourceContents)
            assert contents.text.startswith('# status'), contents.text[:80]
            show('juju://status-doc', contents.text, limit=8)

        def status_matches_juju() -> None:
            result = ok(mcp.call_tool('status', {'model': model}))
            structured = result.structured_content
            assert structured is not None
            live = juju.status()
            assert structured['model']['name'] == live.model.name
            assert set(structured['applications']) == set(live.apps)
            first = result.content[0]
            assert isinstance(first, TextContent) and json.loads(first.text) == structured
            app = structured['applications'][APP]
            unit = next(iter(app['units'].values()))
            assert unit['workload-status']['current'] == 'active', unit['workload-status']
            show_json(
                f'status {model} (structuredContent, excerpt)',
                {
                    'model': {k: structured['model'][k] for k in ('name', 'controller', 'cloud', 'version')},
                    'applications': {
                        APP: {
                            'charm': app['charm'],
                            'application-status': app['application-status'],
                            'units': {
                                name: {
                                    'machine': u['machine'],
                                    'public-address': u.get('public-address'),
                                    'workload-status': u['workload-status'],
                                }
                                for name, u in app['units'].items()
                            },
                        }
                    },
                },
                limit=40,
            )

        def models_lists_demo() -> None:
            arguments = {'controller': target.controller} if target.controller else {}
            result = ok(mcp.call_tool('models', arguments))
            models = result.structured_content['models']  # type: ignore[index]
            names = {m['short-name'] for m in models}
            assert target.model in names, names
            show(
                'models',
                '\n'.join(
                    f'{m["short-name"]:<24} {m.get("cloud", "")}/{m.get("region", "")}'
                    f'  {m["status"]["current"]}'
                    for m in models
                ),
            )

        def show_application() -> None:
            result = ok(mcp.call_tool('show-application', {'model': model, 'args': [APP]}))
            assert result.structured_content and APP in result.structured_content
            show_json(f'show-application {APP}', result.structured_content, limit=16)

        def config_matches_juju() -> None:
            result = ok(mcp.call_tool('config', {'model': model, 'args': [APP]}))
            structured = result.structured_content
            assert structured and structured['application'] == APP
            expected = json.loads(juju.cli('config', APP, '--format', 'json'))
            assert set(structured['settings']) == set(expected['settings'])
            settings = structured['settings']
            show(
                f'config {APP} ({len(settings)} settings, excerpt)',
                '\n'.join(
                    f'{name:<32} {s.get("type", ""):<8} value={s.get("value", "<unset>")!r}'
                    for name, s in sorted(settings.items())
                ),
            )

        checks.run('server instructions describe Juju', instructions)
        checks.run('tool list covers the Juju CLI', tool_list)
        checks.run('status tool is annotated read-only with JSON output', status_tool_definition)
        checks.run('juju://status-doc resource is the Juju help', doc_resource)
        checks.run(f'status of {model} matches `juju status`', status_matches_juju)
        checks.run(f'models lists {target.model}', models_lists_demo)
        checks.run(f'show-application {APP}', show_application)
        checks.run(f'config {APP} matches `juju config`', config_matches_juju)

    banner('stdio server, --read-only')
    with common.mcp_server(binary, '--read-only') as ro:

        def readonly_tool_list() -> None:
            tools = {t.name for t in ro.list_tools()}
            present = {'status', 'models', 'config', 'remove-unit'}
            absent = {'deploy', 'add-unit', 'integrate', 'destroy-model'}
            assert present <= tools, present - tools
            assert not (absent & tools), absent & tools
            assert all(t.annotations and t.annotations.read_only_hint for t in ro.list_tools())
            dropped = sorted(stdio_tools - tools)
            show(f'{len(tools)} tools registered, {len(dropped)} dropped', ', '.join(dropped), limit=6)
            policed = [t for t in ro.list_tools() if 'Read-only mode:' in (t.description or '')]
            show(
                'tools kept behind a read-only policy',
                '\n'.join(
                    f'{t.name}: ' + t.description.split('Read-only mode:', 1)[1].strip()  # type: ignore[union-attr]
                    for t in policed
                ),
            )

        def readonly_query() -> None:
            result = ok(ro.call_tool('config', {'model': model, 'args': [APP]}))
            assert result.structured_content['application'] == APP  # type: ignore[index]
            settings = result.structured_content['settings']  # type: ignore[index]
            show(f'config {APP}', f'application={APP}, {len(settings)} settings returned')

        def readonly_rejects_write() -> None:
            before = juju.config(APP)
            result = ro.call_tool('config', {'model': model, 'args': [APP, 'profile=testing']})
            assert result.is_error and 'read-only mode' in result_text(result), result_text(result)
            assert juju.config(APP) == before, 'rejected write reached Juju'
            show(f'config {APP} profile=testing -> isError', result_text(result))

        def readonly_rejects_host_write() -> None:
            out = pathlib.Path('/tmp/mcp-juju-demo-status.json')
            result = ro.call_tool('status', {'model': model, 'output': str(out)})
            assert result.is_error and 'writes to the host' in result_text(result)
            assert not out.exists()
            show(f'status --output {out} -> isError', result_text(result))

        def readonly_dry_run() -> None:
            units = set(juju.status().apps[APP].units)
            unit = sorted(units)[0]
            refused = ro.call_tool('remove-unit', {'model': model, 'args': [unit]})
            assert refused.is_error and 'read-only mode' in result_text(refused)
            show(f'remove-unit {unit} -> isError', result_text(refused))
            dry = ok(ro.call_tool('remove-unit', {'model': model, 'args': [unit], 'dry-run': True}))
            assert unit in result_text(dry), result_text(dry)
            assert set(juju.status().apps[APP].units) == units, 'dry run removed a unit'
            show(f'remove-unit {unit} --dry-run -> ok', result_text(dry))

        checks.run('only read-only tools are registered', readonly_tool_list)
        checks.run('config query is accepted', readonly_query)
        checks.run('config write is rejected before reaching Juju', readonly_rejects_write)
        checks.run('status --output is rejected (host write)', readonly_rejects_host_write)
        checks.run('remove-unit needs --dry-run and changes nothing', readonly_dry_run)

    banner('Streamable HTTP server with bearer token')
    with http_server(str(binary)) as srv:
        show('server', f'{srv.url}  (token {srv.token[:6]}...)', limit=1)

        def http_requires_token() -> None:
            body = json.dumps(
                {
                    'jsonrpc': '2.0',
                    'id': 1,
                    'method': 'initialize',
                    'params': {
                        'protocolVersion': '2025-06-18',
                        'capabilities': {},
                        'clientInfo': {'name': 'demo', 'version': '0'},
                    },
                }
            ).encode()
            headers = {
                'Content-Type': 'application/json',
                'Accept': 'application/json, text/event-stream',
            }
            req = urllib.request.Request(srv.url, data=body, headers=headers, method='POST')
            try:
                urllib.request.urlopen(req, timeout=10)
            except urllib.error.HTTPError as e:
                assert e.code == 401, e.code
                www = e.headers.get('WWW-Authenticate', '')
                assert 'Bearer' in www
                show('POST initialize without Authorization', f'HTTP {e.code}\nWWW-Authenticate: {www}')
            else:
                raise AssertionError('request without a token was accepted')

        def http_matches_stdio() -> None:
            with McpJujuClient.http(srv.url, token=srv.token) as http:
                tools = {t.name for t in http.list_tools()}
                assert tools == stdio_tools
                result = ok(http.call_tool('status', {'model': model}))
                apps = result.structured_content['applications']  # type: ignore[index]
                assert set(apps) == set(juju.status().apps)
                show(
                    'with Authorization: Bearer <token>',
                    f'{len(tools)} tools (same set as stdio)\n'
                    f'status {model}: applications={sorted(apps)}, '
                    f'{APP} is {apps[APP]["application-status"]["current"]}',
                )

        checks.run('request without a token gets 401', http_requires_token)
        checks.run('tool list and status over HTTP match stdio', http_matches_stdio)

    return checks.summary()


if __name__ == '__main__':
    sys.exit(main())
