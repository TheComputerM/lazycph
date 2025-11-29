import subprocess
from pathlib import Path

from lazycph.utils.runtimes import utils


def python(file: Path, stdin: str) -> str:
    result = subprocess.run(
        f"python {file.absolute()}",
        check=True,
        shell=True,
        capture_output=True,
        text=True,
        input=stdin,
        timeout=utils.TIMEOUT,
    )
    return result.stdout.strip()
