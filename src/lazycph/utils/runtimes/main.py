from pathlib import Path

from lazycph.utils.runtimes.utils import RuntimeExtension


def execute(file: Path, stdin: str) -> str:
    match file.suffix:
        case RuntimeExtension.PYTHON.value:
            from .python import python

            return python(file, stdin)
        case RuntimeExtension.CPP.value:
            from .cpp import cpp

            return cpp(file, stdin)
        case _:
            return "Unsupported file type"
