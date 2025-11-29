from enum import Enum
from pathlib import Path
from subprocess import CalledProcessError, TimeoutExpired

from textual import work
from textual.app import RenderResult
from textual.reactive import reactive, var
from textual.widgets import ListItem

from lazycph.runtimes import executor, utils


class Status(Enum):
    PENDING = "PE"
    CORRECT = "CA"
    WRONG = "WA"
    COMPILATION_ERROR = "CE"
    RUNTIME_ERROR = "RE"
    TIME_LIMIT_EXCEEDED = "TLE"
    UNKNOWN_ERROR = "XX"


class TestcaseItem(ListItem):
    DEFAULT_CSS = """
    TestcaseItem {
        padding: 1;
    }
    """

    input: var[str] = var("")
    output: var[str] = var("")
    expected_output: var[str] = var("")

    status: reactive[Status] = reactive(Status.PENDING)

    @property
    def index(self) -> int:
        """Dynamically compute the current index in parent's children."""
        assert self.parent is not None
        return self.parent.children.index(self)

    def render(self) -> RenderResult:
        output = f"Testcase {self.index}"
        if self.status is not Status.PENDING:
            color = "$text-success"
            if self.status is not Status.CORRECT:
                color = "$text-error"
            output = f"{output} [{color}]({self.status.value})[/]"
        return output

    def _is_expected_output_correct(self) -> bool:
        """Compare the actual output with the expected output."""
        return self.output.split() == self.expected_output.split()

    @work(thread=True)
    def run(self, file: Path):
        """
        Runs the given testcase using the specified file and updates the output and status accordingly.
        """
        self.output = "Running..."
        try:
            self.output = executor.execute(file, self.input)
        except utils.CompilationError as e:
            self.output = str(e)
            self.status = Status.COMPILATION_ERROR
        except TimeoutExpired:
            self.output = "Time Limit Exceeded"
            self.status = Status.TIME_LIMIT_EXCEEDED
        except CalledProcessError as e:
            self.output = f"{e}\n{e.stderr}"
            self.status = Status.RUNTIME_ERROR
        except Exception as e:
            self.output = f"Unexpected Error: {e}"
            self.status = Status.UNKNOWN_ERROR
        else:
            self.status = (
                Status.CORRECT if self._is_expected_output_correct() else Status.WRONG
            )

    def to_json(self) -> dict:
        return {
            "input": self.input,
            "expected_output": self.expected_output,
            "output": self.output,
            "status": self.status.value,
        }

    @staticmethod
    def from_json(data: dict) -> "TestcaseItem":
        testcase = TestcaseItem()
        testcase.set_reactive(TestcaseItem.input, data["input"])
        testcase.set_reactive(TestcaseItem.expected_output, data["expected_output"])
        testcase.set_reactive(TestcaseItem.output, data["output"])
        testcase.set_reactive(TestcaseItem.status, Status(data["status"]))
        return testcase
