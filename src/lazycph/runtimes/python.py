import subprocess
from pathlib import Path

from lazycph.runtimes import utils


def python(file: Path, stdin: str) -> str:
    result = subprocess.run(
        ["python3", file.absolute()],
        check=True,
        capture_output=True,
        text=True,
        input=stdin,
        timeout=utils.TIMEOUT,
    )
    return result.stdout.strip()
