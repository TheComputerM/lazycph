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

    stderr: str

    def __init__(self, stderr: str) -> None:
        self.stderr = stderr
        super().__init__(f"CompilationError:\n{stderr}")
