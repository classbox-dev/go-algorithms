#!/usr/bin/env python3

import argparse
import dataclasses
import enum
import functools as ft
import io
import json
import os
import pathlib
import shutil
import subprocess
import sys
import tempfile
import time
import traceback
import typing
import urllib.request
import zipfile

import requests
import yaml


# ------------------------------------------------------------------------------


class Status(enum.Enum):
    SUCCESS = 'success'
    FAILURE = 'failure'
    EXCEPTION = 'exception'


@dataclasses.dataclass
class Stage:
    status: Status
    name: str
    output: str = ''
    test: str = ''

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
                status=Status.EXCEPTION,
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
            self.failure(exc_val.stdout)
        elif issubclass(exc_type, Exception):
            self.exception(f"Unexpected error:\n{traceback.format_exc()}")
        return True

    def success(self, output="", test=""):
        self.stages.append(Stage(Status.SUCCESS, self.__name, output=output, test=test))

    def failure(self, output, test=""):
        self.stages.append(Stage(Status.FAILURE, self.__name, output=output, test=test))

    def exception(self, output, test=""):
        self.stages.append(Stage(Status.EXCEPTION, self.__name, output=output, test=test))

    def is_success(self):
        return all(s.status == Status.SUCCESS for s in self.stages)


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
    with (path / '.stdlib.yaml').open() as f:
        tests = yaml.safe_load(f)['tests']
        return [test['name'] for test in tests if (path / test['name']).is_dir()]


def _compile(tests: typing.List[str],
             input_path: pathlib.Path,
             output_path: pathlib.Path):
    stager = Stager()

    with tempfile.TemporaryDirectory() as tmp_dir:
        tmp_path = pathlib.Path(tmp_dir).absolute()
        tests_dir = tmp_path / 'stdlib-tests'

        shutil.copytree('/stdlib-tests', str(tests_dir))
        shutil.copytree(str(input_path), str(tmp_path / 'stdlib'))

        for test_name in tests:
            with stager(f"build::{test_name}") as st:
                _run(
                    [
                        'go', 'test', '-c', '-ldflags', '-s -w', '-trimpath',
                        '-o', str(output_path / f'{test_name}.test'),
                        f'hsecode.com/stdlib-tests/{test_name}'
                    ],
                    cwd=str(tests_dir)
                )
                st.success(test=test_name)

    return stager.stages


def _test(tests: typing.List[str]) -> typing.List[Stage]:
    stager = Stager()

    for test_name in tests:
        with stager(f"build::{test_name}") as st:
            _run(
                ['go', 'test', f'hsecode.com/stdlib-tests/{test_name}'],
                cwd=str("/stdlib-tests")
            )
            st.success(test=test_name)

    return stager.stages


# ------------------------------------------------------------------------------

@command
def build_tests(args: argparse.Namespace) -> typing.List[Stage]:
    stager = Stager()

    tmp_dir = tempfile.mkdtemp()
    tmp_path = pathlib.Path(tmp_dir).absolute()

    with stager("fetch") as st:
        # noinspection PyBroadException
        try:
            r = requests.get(args.url)
            r.raise_for_status()
        except Exception:
            st.exception(f'Could not download source code')
            return stager.stages

        # noinspection PyBroadException
        try:
            with zipfile.ZipFile(io.BytesIO(r.content)) as f:
                f.extractall(tmp_path)
        except Exception:
            st.exception(f'Could not unpack source code')
            return stager.stages

        nested = list(tmp_path.glob("*"))
        if len(nested) != 1:
            st.exception(f'Could not read unpacked source code')
            return stager.stages
        src_dir = nested[0]

        st.success()

    with stager("format") as st:
        r = _run(['gofmt', '-l', '.'], cwd=str(src_dir))
        output = r.stdout.strip()
        if output:
            errfmt_files = ', '.join(output.split('\n'))
            st.failure(
                f'Formatting error: {errfmt_files}.\n'
                f'Use `gofmt` to format these files'
            )
            return stager.stages
        else:
            st.success()

    with stager("lint") as st:
        _run(['stdlib-linter', '-c', '/linter_config.yaml', '.'], cwd=str(src_dir))
        st.success()

    if not stager.is_success():
        return stager.stages

    with stager("configure") as st:
        config_file = src_dir / '.stdlib.yaml'
        try:
            with config_file.open() as f:
                try:
                    tests = yaml.safe_load(f)["tests"]
                except yaml.YAMLError as exc:
                    st.failure(f'Could not parse `{config_file.name}`\n{exc}')
                    return stager.stages
        except Exception as exc:
            st.failure(f'Could not read `{config_file.name}`\n{exc}')
            return stager.stages

        tests = tests or []  # initial config is parsed as `None`

        if type(tests) is not list or any(type(x) is not str for x in tests):
            st.failure(
                f'`{config_file.name}` is invalid: '
                f'it must contain a list of test names only'
            )
            return stager.stages

        all_tests = set(_all_tests())
        invalid_tests = sorted(set(tests) - all_tests)
        if invalid_tests:
            st.failure(
                f'`{config_file.name}` contains invalid test names: '
                f'{", ".join(invalid_tests)}'
            )
            return stager.stages

        tests = sorted(set(_all_tests()) & set(tests))
        if tests:
            st.success()
        else:
            st.success("No tests selected")
            return stager.stages

    stager.stages.extend(_compile(tests, src_dir, args.output_path))
    return stager.stages


@command
def build_baseline(args: argparse.Namespace) -> typing.List[Stage]:
    tests = _all_tests()
    return _compile(tests, pathlib.Path("/stdlib"), args.output_path)


def test_all(args: argparse.Namespace):
    results = _test(_all_tests())

    if all(r.status == Status.SUCCESS for r in results):
        print("OK")
        sys.exit(0)

    for stage in results:
        if stage.status != Status.SUCCESS:
            print(f"{stage.name}: {stage.status.name}\n")
            print(stage.output)
            print()
    sys.exit(1)


@command
def build_docs(args: argparse.Namespace) -> typing.List[Stage]:
    stager = Stager()

    with stager("docs") as st:

        godoc = subprocess.Popen(
            ['godoc', '-http=:6060', '-links=false', '-templates=/opt/static'],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.STDOUT,
            cwd=str("/stdlib")
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
        _run([
            'wget', '-r', '-np', '-N', '-nH', '--cut-dirs=3',
            '-E', '-p', '-k', '-e', 'robots=off',
            'http://127.0.0.1:6060/pkg/hsecode.com/stdlib'
        ], cwd=str(docs_path))

        godoc.kill()

        # Postprocess
        (docs_path / 'stdlib.html').rename(docs_path / 'index.html')
        shutil.copy2('/linter_config.yaml', str(docs_path / 'linter.yaml'))

        for p in docs_path.glob("**/*.html"):
            src = p.read_text()
            src = src.replace("__HEAD__", f'<a href="{args.web}">hsecode.com/stdlib</a> / <a href="{args.docs}">docs</a>')
            src = src.replace('<a href="http://127.0.0.1:6060/pkg/">GoDoc</a></div>', f'<a href="{args.docs}">docs</a>')
            p.write_text(src)

        st.success()

    return stager.stages


@command
def build_meta(args: argparse.Namespace) -> typing.List[Stage]:
    path = pathlib.Path('/stdlib-tests')
    with (path / '.stdlib.yaml').open() as f:
        tests = yaml.safe_load(f)["tests"]
    return [Stage(Status.SUCCESS, "meta", json.dumps(tests))]


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
    # noinspection PyTypeChecker
    parser.add_argument(
        '-o', dest='output_path', default='/out',
        type=ft.partial(arg_path, write=True)
    )

    subparsers = parser.add_subparsers(dest='command', required=True)

    parser_tests = subparsers.add_parser('tests')
    parser_tests.add_argument('url')
    parser_tests.set_defaults(func=build_tests)

    parser_tests = subparsers.add_parser('baseline')
    parser_tests.set_defaults(func=build_baseline)

    parser_docs = subparsers.add_parser('docs')
    parser_docs.add_argument('--web', required=True)
    parser_docs.add_argument('--docs', required=True)
    parser_docs.set_defaults(func=build_docs)

    parser_meta = subparsers.add_parser('meta')
    parser_meta.set_defaults(func=build_meta)

    parser_meta = subparsers.add_parser('test-all')
    parser_meta.set_defaults(func=test_all)

    args_ = parser.parse_args()

    output = args_.func(args_)
    print(json.dumps(output, indent=2))


if __name__ == '__main__':
    main()
