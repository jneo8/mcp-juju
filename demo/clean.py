#!/usr/bin/env python3
"""Destroy the mcp-juju demo environment.

Destroys the demo model (and its storage) through the MCP server's
`destroy-model` tool and confirms with Jubilant that it is gone. Does
nothing when the model does not exist.
"""

from __future__ import annotations

import shutil
import sys
import time

import common
from common import BUILD_DIR, DESTROY_TIMEOUT, banner, call, step


def main() -> int:
    p = common.parser(__doc__)
    p.add_argument(
        '--remove-build',
        action='store_true',
        help=f'also delete the binary built under {BUILD_DIR.relative_to(common.REPO_ROOT)}',
    )
    args = p.parse_args()
    target = common.target_from(args)

    banner(f'Model {target.qualified}')
    if not common.model_exists(target):
        step('model does not exist, nothing to clean')
    else:
        binary = common.build_binary(args.binary)
        with common.mcp_server(binary) as mcp:
            call(
                mcp,
                'destroy-model',
                {
                    'args': [target.qualified],
                    'no-prompt': True,
                    'destroy-storage': True,
                    'force': True,
                    'timeout': f'{DESTROY_TIMEOUT}s',
                },
            )
        deadline = time.monotonic() + DESTROY_TIMEOUT
        while common.model_exists(target):
            if time.monotonic() > deadline:
                sys.exit(f'model {target.qualified} still exists after {DESTROY_TIMEOUT}s')
            time.sleep(2)
        step('model destroyed')

    if args.remove_build and BUILD_DIR.exists():
        shutil.rmtree(BUILD_DIR)
        step(f'removed {BUILD_DIR.relative_to(common.REPO_ROOT)}')
    return 0


if __name__ == '__main__':
    sys.exit(main())
