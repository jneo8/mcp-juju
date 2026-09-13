#!/usr/bin/env python3
"""Deploy the mcp-juju demo environment.

Creates the demo model and deploys the charm through the MCP server's own
`add-model` and `deploy` tools, then waits with Jubilant until the
application is active. Safe to rerun: existing pieces are reused.
"""

from __future__ import annotations

import sys

import common
import jubilant
from common import APP, DEFAULT_CHARM, DEPLOY_TIMEOUT, banner, call, step


def main() -> int:
    p = common.parser(__doc__)
    p.add_argument(
        '--charm',
        default=DEFAULT_CHARM,
        help=f'charm to deploy (default: {DEFAULT_CHARM}, or MCP_JUJU_DEMO_CHARM)',
    )
    args = p.parse_args()
    target = common.target_from(args)
    binary = common.build_binary(args.binary)
    juju = common.juju_for(target)

    with common.mcp_server(binary) as mcp:
        banner(f'Model {target.qualified}')
        if common.model_exists(target):
            step('model already exists, reusing it')
        else:
            arguments: dict = {'args': [target.model], 'no-switch': True}
            if target.controller:
                arguments['controller'] = target.controller
            call(mcp, 'add-model', arguments)

        banner(f'Application {APP} ({args.charm})')
        if APP in juju.status().apps:
            step('application already deployed, reusing it')
        else:
            call(mcp, 'deploy', {'model': target.qualified, 'args': [args.charm, APP]})

        step(f'waiting up to {DEPLOY_TIMEOUT // 60} min for {APP} to become active')
        status = juju.wait(
            lambda s: jubilant.all_active(s, APP),
            error=jubilant.any_error,
            timeout=DEPLOY_TIMEOUT,
        )

    banner('Ready')
    app = status.apps[APP]
    print(f'model:   {target.qualified}')
    print(f'charm:   {app.charm} (rev {app.charm_rev})')
    for name, unit in sorted(app.units.items()):
        print(f'unit:    {name} on machine {unit.machine}, {unit.workload_status.current}')
    print()
    print('Connect an MCP client to the same Juju client store, for example:')
    print(f'  claude mcp add juju -- {binary} --server-type stdio')
    print(f'  claude mcp add juju-ro -- {binary} --server-type stdio --read-only')
    print()
    print('Then:')
    print('  just demo-verify   # check the server against this deployment')
    print('  just demo-clean    # destroy the demo model')
    return 0


if __name__ == '__main__':
    sys.exit(main())
