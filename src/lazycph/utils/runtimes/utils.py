from enum import Enum

TIMEOUT = 5.0
COMPILATION_TIMEOUT = 10.0


class RuntimeExtension(Enum):
    PYTHON = ".py"
    CPP = ".cpp"
    C = ".c"


class CompilationError(Exception):
    """
    Exception raised for errors during the compilation process of runtimes.
    """

    def __init__(
        self,
        stderr: str,
    ) -> None:
        super().__init__("CompilationError")
        self.stderr = stderr

    def __str__(self):
        return f"CompilationError: {self.stderr}"
