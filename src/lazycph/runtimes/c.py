import subprocess
import tempfile
from pathlib import Path

from lazycph.runtimes import utils


def c(file: Path, stdin: str) -> str:
    """
    Compile a C source file using gcc and execute it with the given stdin.
    """
    # Create a temporary file for the executable
    with tempfile.NamedTemporaryFile(suffix="", delete=False) as temp_exe:
        exe_path = temp_exe.name

    try:
        compile_result = subprocess.run(
            ["gcc", str(file), "-o", exe_path, "-std=gnu23"],
            text=True,
            capture_output=True,
            timeout=utils.COMPILATION_TIMEOUT,
        )

        if compile_result.returncode != 0:
            raise utils.CompilationError(compile_result.stderr)

        run_result = subprocess.run(
            [exe_path],
            check=True,
            text=True,
            capture_output=True,
            input=stdin,
            timeout=utils.TIMEOUT,
        )

        return run_result.stdout.strip()
    finally:
        # Clean up the temporary executable
        Path(exe_path).unlink(missing_ok=True)
