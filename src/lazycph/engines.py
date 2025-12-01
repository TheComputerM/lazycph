import subprocess
import tempfile
from pathlib import Path


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
    compiled: bool

    def __init__(self, command: str, compiled: bool) -> None:
        self.command = command
        self.compiled = compiled

    @staticmethod
    def execute_interpreted(file: Path, stdin: str, command: str) -> str:
        result = subprocess.run(
            command.format(file=f'"{file.resolve()}"'),
            shell=True,
            check=True,
            capture_output=True,
            text=True,
            input=stdin,
            timeout=5.0,
        )
        return result.stdout.strip()

    @staticmethod
    def execute_compiled(file: Path, stdin: str, command: str) -> str:
        from random import choices
        from string import ascii_letters

        random_name = "".join(choices(ascii_letters, k=8))
        exe_path = Path(tempfile.gettempdir()) / f"lazycph-{random_name}"

        try:
            compile_result = subprocess.run(
                command.format(file=f'"{file.resolve()}"', temp=exe_path),
                shell=True,
                text=True,
                capture_output=True,
                timeout=10.0,
            )

            if compile_result.returncode != 0:
                raise CompilationError(compile_result.stderr)

            run_result = subprocess.run(
                [exe_path],
                check=True,
                text=True,
                input=stdin,
                timeout=5.0,
                cwd=file.parent,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
            )

            return run_result.stdout.strip()
        finally:
            # Clean up the temporary executable
            Path(exe_path).unlink(missing_ok=True)

    def execute(self, file: Path, stdin: str) -> str:
        assert file.is_file(), "The provided path is not a file."
        assert file.exists(), "The provided file does not exist."
        if self.compiled:
            return self.execute_compiled(file, stdin, self.command)
        return self.execute_interpreted(file, stdin, self.command)


available: dict[str, Engine] = {
    ".py": Engine("python3 {file}", compiled=False),
    ".cpp": Engine("g++ {file} -o {temp} -std=c++17", compiled=True),
    ".c": Engine("g++ {file} -o {temp} -std=gnu23", compiled=True),
    ".rs": Engine("rustc {file} -o {temp}", compiled=True),
}


def execute(file: Path, stdin: str) -> str:
    if file.suffix not in available:
        return "Unsupported file type"
    command = available[file.suffix]
    return command.execute(file, stdin)
