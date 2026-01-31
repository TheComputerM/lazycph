import subprocess
import tempfile
from pathlib import Path
from typing import Literal

COMPILE_TIMEOUT = 10.0
EXECUTION_TIMEOUT = 5.0


def _run_command(
    cmd: str | list[str],
    stdin: str,
    timeout: float,
    check: bool = True,
    **kwargs,
) -> subprocess.CompletedProcess[str]:
    """
    Run a subprocess with common options: text mode, check for errors,
    and merge stderr into stdout. Additional kwargs are passed to subprocess.run.
    """
    return subprocess.run(
        cmd,
        check=check,
        text=True,
        input=stdin,
        timeout=timeout,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        **kwargs,
    )


class CompilationError(Exception):
    """
    Exception raised for errors during the compilation process of runtimes.
    """

    stderr: str

    def __init__(self, stderr: str) -> None:
        self.stderr = stderr
        super().__init__(f"CompilationError:\n{stderr}")


class Engine:
    command: str
    mode: Literal["compile", "interpret"]

    def __init__(self, command: str, mode: Literal["compile", "interpret"]) -> None:
        self.command = command
        self.mode = mode

    @staticmethod
    def execute_interpreted(file: Path, stdin: str, command: str) -> str:
        result = _run_command(
            command.format(file=f'"{file.resolve()}"'),
            stdin=stdin,
            timeout=EXECUTION_TIMEOUT,
            shell=True,
            cwd=file.parent,
        )
        return result.stdout.strip()

    @staticmethod
    def execute_compiled(file: Path, stdin: str, command: str) -> str:
        from random import choices
        from string import ascii_letters

        random_name = "".join(choices(ascii_letters, k=8))
        exe_path = Path(tempfile.gettempdir()) / f"lazycph-{random_name}"

        try:
            compile_result = _run_command(
                command.format(file=f'"{file.resolve()}"', temp=exe_path),
                stdin="",
                timeout=COMPILE_TIMEOUT,
                shell=True,
                check=False,
            )

            if compile_result.returncode != 0:
                raise CompilationError(compile_result.stdout)

            run_result = _run_command(
                [str(exe_path)],
                stdin=stdin,
                timeout=EXECUTION_TIMEOUT,
                cwd=file.parent,
            )

            return run_result.stdout.strip()
        finally:
            # Clean up the temporary executable
            Path(exe_path).unlink(missing_ok=True)

    def execute(self, file: Path, stdin: str) -> str:
        assert file.exists(), "The provided file does not exist."
        assert file.is_file(), "The provided path is not a file."
        if self.mode == "compile":
            return self.execute_compiled(file, stdin, self.command)
        return self.execute_interpreted(file, stdin, self.command)


available: dict[str, Engine] = {
    ".py": Engine("python3 {file}", mode="interpret"),
    ".cpp": Engine("g++ {file} -o {temp} -std=c++17", mode="compile"),
    ".c": Engine("gcc {file} -o {temp} -std=gnu23", mode="compile"),
    ".rs": Engine("rustc {file} -o {temp}", mode="compile"),
    ".zig": Engine("zig build-exe {file} -femit-bin={temp}", mode="compile"),
}


def execute(file: Path, stdin: str) -> str:
    if file.suffix not in available:
        return "Unsupported file type"
    command = available[file.suffix]
    return command.execute(file, stdin)
