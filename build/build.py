#!/usr/bin/env python3

import argparse
import dataclasses
import enum
import functools as ft
import json
import os
import pathlib
import subprocess
import sys
import time
import shutil
import typing
import urllib.request


class Status(enum.Enum):
    OK = 'ok'
    FAILURE = 'failure'


@dataclasses.dataclass
class Stage:
    status: Status
    name: str = None
    message: str = None
    output: str = None

    def json(self):
        data = dataclasses.asdict(self)
        data['status'] = self.status.name
        return data


def run_cmd(args: typing.List[str], **kwargs):
    return subprocess.run(
        args,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        encoding='utf-8',
        errors='replace',
        text=True,
        **kwargs
    )


def submission(args) -> int:
    stages: typing.List[Stage] = []

    # Format
    name = 'format'

    c = run_cmd(['gofmt', '-l', str(args.input_path)])
    output = c.stdout.strip()
    if output:
        errfmt_files = ', '.join(output.split('\n'))
        msg = (
            f'Formatting error: {errfmt_files}. '
            f'Use `gofmt` to format these files'
        )
        result = Stage(status=Status.FAILURE, name=name, message=msg)
    else:
        result = Stage(status=Status.OK, name=name)

    stages.append(result)

    # Lint

    #


    name = 'build'
    env = os.environ.copy()
    env['GO111MODULE'] = 'off'
    c = run_cmd(
        ['go', 'test', '-c', '-o', str(args.output_path / 'tests')],
        env=env, cwd=str(args.input_path)
    )
    if c.returncode != 0:
        output = c.stdout.strip()
        result = Stage(
            status=Status.FAILURE,
            name=name, message='Build error', output=output
        )
    else:
        result = Stage(status=Status.OK, name=name)

    stages.append(result)

    json_result = [x.json() for x in stages]
    with (args.output_path / 'result.json').open('w') as f:
        json.dump(json_result, f)

    return 0


def make_doc(args):
    subprocess.Popen(
        ['godoc', '-http=:6060', '-links=false', '-templates=/opt/static'],
        cwd=str(args.input_path)
    )

    deadline = time.monotonic() + 3
    while time.monotonic() < deadline:
        try:
            urllib.request.urlopen('http://127.0.0.1:6060', timeout=0.5)
            break
        except urllib.request.URLError:
            continue
    else:
        raise RuntimeError('could not start documentation server')

    docs_path = args.output_path / 'docs'
    shutil.rmtree(str(docs_path), ignore_errors=True)
    docs_path.mkdir()
    r = run_cmd([
        'wget', '-r', '-np', '-N', '-nH', '--cut-dirs=3',
        '-E', '-p', '-k', '-e', 'robots=off',
        'http://127.0.0.1:6060/pkg/hsecode.com/stdlib'
    ], cwd=str(docs_path))

    if r.returncode:
        print(r.stdout)

    # Postprocess
    (docs_path / 'stdlib.html').rename(docs_path / 'index.html')

    return 0


def arg_path(x, write=False):
    path = pathlib.Path(x).absolute()
    if not path.is_dir():
        raise argparse.ArgumentTypeError(f'invalid path: {x}')
    if not os.access(path, os.R_OK):
        raise argparse.ArgumentTypeError(f'unreadable path: {x}')
    if write and not os.access(path, os.W_OK):
        raise argparse.ArgumentTypeError(f'non-writable path: {x}')
    return path


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument(
        '-i', dest='input_path', default='/in',
        type=arg_path
    )
    parser.add_argument(
        '-o', dest='output_path', default='/out',
        type=ft.partial(arg_path, write=True)
    )

    subparsers = parser.add_subparsers(dest='command', required=True)

    parser_submission = subparsers.add_parser('submission')
    parser_submission.set_defaults(func=submission)

    parser_docs = subparsers.add_parser('docs')
    parser_docs.set_defaults(func=make_doc)

    args_ = parser.parse_args()

    sys.exit(args_.func(args_))
