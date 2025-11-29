import subprocess
from pathlib import Path
from tempfile import NamedTemporaryFile

import pytest

from lazycph.runtimes import utils
from lazycph.runtimes.executor import execute


class TestPython:
    def test_basic(self):
        with NamedTemporaryFile(suffix=".py", mode="w+") as file:
            file.write("print(input())")
            file.flush()
            output = execute(Path(file.name), "hello world")
            assert output == "hello world"

    def test_timeout(self):
        with NamedTemporaryFile(suffix=".py", mode="w+") as file:
            file.write(f"import time\ntime.sleep({utils.TIMEOUT + 1})")
            file.flush()
            with pytest.raises(subprocess.TimeoutExpired):
                execute(Path(file.name), "")

    def test_runtime_error(self):
        with NamedTemporaryFile(suffix=".py", mode="w+") as file:
            file.write("print(1/0)")
            file.flush()
            with pytest.raises(subprocess.CalledProcessError) as exc_info:
                execute(Path(file.name), "")
            assert "ZeroDivisionError" in exc_info.value.stderr


class TestCPP:
    def test_basic(self):
        with NamedTemporaryFile(suffix=".cpp", mode="w") as file:
            file.write("""
#include <stdio.h>
int main() { printf("hello world"); return 0; }
""")
            file.flush()
            output = execute(Path(file.name), "")
            assert output == "hello world"

    def test_compilation_error(self):
        with NamedTemporaryFile(suffix=".cpp", mode="w+") as file:
            file.write("MAKIMA IS LISTENING")
            file.flush()
            with pytest.raises(utils.CompilationError):
                execute(Path(file.name), "")

    def test_runtime_error(self):
        with NamedTemporaryFile(suffix=".cpp", mode="w+") as file:
            file.write("""int main() {return -1;} """)
            file.flush()
            with pytest.raises(subprocess.CalledProcessError):
                execute(Path(file.name), "")


class TestC:
    def test_basic(self):
        with NamedTemporaryFile(suffix=".c", mode="w") as file:
            file.write("""
#include <stdio.h>
int main() { printf("hello world"); return 0; }
""")
            file.flush()
            output = execute(Path(file.name), "")
            assert output == "hello world"

    def test_compilation_error(self):
        with NamedTemporaryFile(suffix=".c", mode="w+") as file:
            file.write("MAKIMA IS LISTENING")
            file.flush()
            with pytest.raises(utils.CompilationError):
                execute(Path(file.name), "")

    def test_runtime_error(self):
        with NamedTemporaryFile(suffix=".c", mode="w+") as file:
            file.write("""int main() {return -1;} """)
            file.flush()
            with pytest.raises(subprocess.CalledProcessError):
                execute(Path(file.name), "")
