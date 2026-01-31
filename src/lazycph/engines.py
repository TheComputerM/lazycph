import subprocess
import tempfile
from dataclasses import dataclass
from pathlib import Path
from shlex import quote
from typing import Any, Literal


def _run_command(
    cmd: str | list[str],
    stdin: str | None,
    timeout: float,
    check: bool = True,
    **kwargs: Any,
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

    output: str
    command: str
    returncode: int

    def __init__(self, output: str, command: str, returncode: int) -> None:
        self.output = output
        self.command = command
        self.returncode = returncode
        super().__init__(str(self))

    def __str__(self) -> str:
        return f"Compilation failed (exit code {self.returncode}):\n{self.output}"


@dataclass
class Engine:
    command: str
    mode: Literal["compile", "interpret"]
    compile_timeout: float = 10.0
    execution_timeout: float = 5.0

    def _interpret(self, file: Path, stdin: str) -> str:
        result = _run_command(
            self.command.format(file=quote(str(file.resolve()))),
            stdin=stdin,
            timeout=self.execution_timeout,
            shell=True,
            cwd=file.parent,
        )
        return result.stdout.strip()

    def _compile(self, file: Path, stdin: str) -> str:
        from random import choices
        from string import ascii_letters

        random_name = "".join(choices(ascii_letters, k=8))
        exe_path = Path(tempfile.gettempdir()) / f"lazycph-{random_name}"

        try:
            compile_command = self.command.format(
                file=quote(str(file.resolve())), temp=exe_path
            )
            compile_result = _run_command(
                compile_command,
                stdin=None,
                timeout=self.compile_timeout,
                shell=True,
                check=False,
                cwd=exe_path.parent,
            )

            if compile_result.returncode != 0:
                raise CompilationError(
                    compile_result.stdout,
                    compile_command,
                    compile_result.returncode,
                )

            run_result = _run_command(
                [str(exe_path)],
                stdin=stdin,
                timeout=self.execution_timeout,
                cwd=file.parent,
            )

            return run_result.stdout.strip()
        finally:
            # Clean up the temporary executable
            exe_path.unlink(missing_ok=True)

    def execute(self, file: Path, stdin: str) -> str:
        assert file.exists(), "The provided file does not exist."
        assert file.is_file(), "The provided path is not a file."
        if self.mode == "compile":
            return self._compile(file, stdin)
        return self._interpret(file, stdin)


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
    engine = available[file.suffix]
    return engine.execute(file, stdin)
