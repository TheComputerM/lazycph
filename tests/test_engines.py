import subprocess
from pathlib import Path
from tempfile import NamedTemporaryFile

import pytest

from lazycph.engines import CompilationError, execute


class TestPython:
    def test_basic(self):
        with NamedTemporaryFile(suffix=".py", mode="w+") as file:
            file.write("print(input())")
            file.flush()
            output = execute(Path(file.name), "hello world")
            assert output == "hello world"

    def test_timeout(self):
        with NamedTemporaryFile(suffix=".py", mode="w+") as file:
            file.write("import time\ntime.sleep(100)")
            file.flush()
            with pytest.raises(subprocess.TimeoutExpired):
                execute(Path(file.name), "")

    def test_runtime_error(self):
        with NamedTemporaryFile(suffix=".py", mode="w+") as file:
            file.write("print(1/0)")
            file.flush()
            with pytest.raises(subprocess.CalledProcessError) as exc_info:
                execute(Path(file.name), "")
            assert "ZeroDivisionError" in exc_info.value.stdout


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
            with pytest.raises(CompilationError):
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
            with pytest.raises(CompilationError):
                execute(Path(file.name), "")

    def test_runtime_error(self):
        with NamedTemporaryFile(suffix=".c", mode="w+") as file:
            file.write("""int main() {return -1;} """)
            file.flush()
            with pytest.raises(subprocess.CalledProcessError):
                execute(Path(file.name), "")


class TestRust:
    def test_basic(self):
        with NamedTemporaryFile(suffix=".rs", mode="w") as file:
            file.write("""fn main() {println!("hello world");}""")
            file.flush()
            output = execute(Path(file.name), "")
            assert output == "hello world"


class TestZig:
    def test_basic(self):
        with NamedTemporaryFile(suffix=".zig", mode="w") as file:
            file.write("""const std = @import("std");
pub fn main() !void {
    std.debug.print("hello world", .{});
}""")
            file.flush()
            output = execute(Path(file.name), "")
            assert output == "hello world"
