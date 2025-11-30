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


class Runtime:
    command: str
    compiled: bool

    def __init__(self, command: str, compiled: bool) -> None:
        self.command = command
        self.compiled = compiled

    @staticmethod
    def execute_interpreted(file: Path, stdin: str, command: str) -> str:
        result = subprocess.run(
            command.format(file=file.absolute()),
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
        with tempfile.NamedTemporaryFile(suffix="", delete=False) as temp_exe:
            exe_path = temp_exe.name

        try:
            compile_result = subprocess.run(
                command.format(file=file.absolute(), temp=exe_path),
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
                capture_output=True,
                input=stdin,
                timeout=5.0,
            )

            return run_result.stdout.strip()
        finally:
            # Clean up the temporary executable
            Path(exe_path).unlink(missing_ok=True)

    def execute(self, file: Path, stdin: str) -> str:
        if self.compiled:
            return self.execute_compiled(file, stdin, self.command)
        return self.execute_interpreted(file, stdin, self.command)


runtimes: dict[str, Runtime] = {
    ".py": Runtime("python3 {file}", compiled=False),
    ".cpp": Runtime("g++ {file} -o {temp} -std=c++17", compiled=True),
    ".c": Runtime("g++ {file} -o {temp} -std=gnu23", compiled=True),
    ".rs": Runtime("rustc {file} -o {temp}", compiled=True),
}


def execute(file: Path, stdin: str) -> str:
    if file.suffix not in runtimes:
        return "Unsupported file type"
    command = runtimes[file.suffix]
    return command.execute(file, stdin)
