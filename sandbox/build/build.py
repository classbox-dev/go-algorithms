#!/usr/bin/env python3

import argparse
import dataclasses
import enum
import functools as ft
import json
import os
import pathlib
import shutil
import subprocess
import tempfile
import time
import traceback
import typing
import urllib.request

import yaml


# ------------------------------------------------------------------------------


class Status(enum.Enum):
    OK = 'ok'
    FAILURE = 'failure'


@dataclasses.dataclass
class Stage:
    status: Status
    name: str
    output: str = ''

    def json(self):
        data = dataclasses.asdict(self)
        data['status'] = self.status.value
        return data


def command(func: typing.Callable[..., typing.List[Stage]]):
    def inner(*args, **kwargs) -> typing.List[dict]:
        # noinspection PyBroadException
        try:
            data = func(*args, **kwargs)
        except Exception:
            data = [Stage(
                status=Status.FAILURE,
                name=f"command::{func.__name__}",
                output=f"Unexpected error:\n{traceback.format_exc()}"
            )]
        return [x.json() for x in data]

    return inner


class Stager:

    def __init__(self):
        self.stages = []
        self.__name = None

    def __enter__(self):
        return self

    def __call__(self, name) -> 'Stager':
        self.__name = name
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        if not exc_type:
            return
        if exc_type is subprocess.CalledProcessError:
            self.fail(exc_val.stdout)
        elif issubclass(exc_type, Exception):
            self.fail(f"Unexpected error:\n{traceback.format_tb(exc_tb)}")
        return True

    def ok(self):
        self.stages.append(Stage(Status.OK, self.__name))

    def fail(self, output):
        self.stages.append(Stage(Status.FAILURE, self.__name, output=output))


# ------------------------------------------------------------------------------


def _run(args: typing.List[str], **kwargs):
    return subprocess.run(
        args,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        encoding='utf-8',
        errors='replace',
        text=True,
        check=True,
        **kwargs
    )


def _all_tests() -> typing.List[str]:
    path = pathlib.Path('/stdlib-tests')
    with (path / '.tests.yaml').open() as f:
        tests = yaml.safe_load(f)
        return [test['id'] for test in tests if (path / test['id']).is_dir()]


def _compile(tests: typing.List[str], output_path: pathlib.Path):
    results: typing.List[Stage] = []

    stager = Stager()

    with tempfile.TemporaryDirectory() as tmp_dir:
        tmp_path = pathlib.Path(tmp_dir).absolute()
        tests_dir = tmp_path / 'stdlib-tests'

        shutil.copytree('/stdlib-tests', str(tests_dir))
        shutil.copytree('/in', str(tmp_path / 'stdlib'))

        for test_name in tests:
            with stager(f"build::{test_name}") as st:
                _run(
                    [
                        'go', 'test', '-c', '-o',
                        str(output_path / f'{test_name}.test'),
                        f'hsecode.com/stdlib-tests/{test_name}'
                    ],
                    cwd=str(tests_dir)
                )
                st.ok()

    return stager.stages


# ------------------------------------------------------------------------------

@command
def build_tests(args: argparse.Namespace) -> typing.List[Stage]:
    stages: typing.List[Stage] = []

    # Format
    name = 'format'

    c = _run(['gofmt', '-l', str(args.input_path)])
    output = c.stdout.strip()
    if output:
        errfmt_files = ', '.join(output.split('\n'))
        msg = (
            f'Formatting error: {errfmt_files}. '
            f'Use `gofmt` to format these files'
        )
        result = Stage(Status.FAILURE, name, output=msg)
    else:
        result = Stage(Status.OK, name)

    stages.append(result)

    # Lint

    #
    #
    # name = 'build'
    # env = os.environ.copy()
    # env['GO111MODULE'] = 'off'
    # c = run_cmd(
    #     ['go', 'test', '-c', '-o', str(args.output_path / 'tests')],
    #     env=env, cwd=str(args.input_path)
    # )
    # if c.returncode != 0:
    #     output = c.stdout.strip()
    #     result = Stage(
    #         status=Status.FAILURE,
    #         name=name, message='Build error', output=output
    #     )
    # else:
    #     result = Stage(status=Status.OK, name=name)
    #
    # stages.append(result)
    #
    # json_result = [x.json() for x in stages]
    # with (args.output_path / 'result.json').open('w') as f:
    #     json.dump(json_result, f)

    return []


@command
def build_baseline(args: argparse.Namespace) -> typing.List[Stage]:
    tests = _all_tests()
    return _compile(tests, args.output_path)


@command
def build_docs(args: argparse.Namespace) -> typing.List[Stage]:
    godoc = subprocess.Popen(
        ['godoc', '-http=:6060', '-links=false', '-templates=/opt/static'],
        stdout=subprocess.DEVNULL,
        stderr=subprocess.STDOUT,
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
    r = _run([
        'wget', '-r', '-np', '-N', '-nH', '--cut-dirs=3',
        '-E', '-p', '-k', '-e', 'robots=off',
        'http://127.0.0.1:6060/pkg/hsecode.com/stdlib'
    ], cwd=str(docs_path))

    if r.returncode:
        print(r.stdout)

    godoc.kill()

    # Postprocess
    (docs_path / 'stdlib.html').rename(docs_path / 'index.html')

    return [Stage(Status.OK, "docs")]


# ------------------------------------------------------------------------------

def main():
    def arg_path(x, write=False):
        path = pathlib.Path(x).absolute()
        if not path.is_dir():
            raise argparse.ArgumentTypeError(f'invalid path: {x}')
        if not os.access(path, os.R_OK):
            raise argparse.ArgumentTypeError(f'unreadable path: {x}')
        if write and not os.access(path, os.W_OK):
            raise argparse.ArgumentTypeError(f'non-writable path: {x}')
        return path

    parser = argparse.ArgumentParser()
    parser.add_argument('-i', dest='input_path', default='/in', type=arg_path)
    parser.add_argument(
        '-o', dest='output_path', default='/out',
        type=ft.partial(arg_path, write=True)
    )

    subparsers = parser.add_subparsers(dest='command', required=True)

    parser_tests = subparsers.add_parser('tests')
    parser_tests.set_defaults(func=build_tests)

    parser_tests = subparsers.add_parser('baseline')
    parser_tests.set_defaults(func=build_baseline)

    parser_docs = subparsers.add_parser('docs')
    parser_docs.set_defaults(func=build_docs)

    args_ = parser.parse_args()

    output = args_.func(args_)
    print(json.dumps(output, indent=2))


if __name__ == '__main__':
    main()
