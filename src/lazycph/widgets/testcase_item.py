from enum import Enum
from pathlib import Path
from subprocess import CalledProcessError, TimeoutExpired

from textual import on, work
from textual.app import RenderResult
from textual.message import Message
from textual.reactive import reactive, var
from textual.widgets import ListItem

from lazycph.engines import CompilationError, execute


class Status(Enum):
    INITIAL = None
    CORRECT = "CA"
    WRONG = "WA"
    COMPILATION_ERROR = "CE"
    RUNTIME_ERROR = "RE"
    TIME_LIMIT_EXCEEDED = "TLE"
    UNKNOWN_ERROR = "XX"


class TestcaseItem(ListItem):
    class StateChanged(Message):
        """Message sent when the output and status of the testcase are updated."""

        output: str
        status: Status

        def __init__(self, output: str, status: Status) -> None:
            self.output = output
            self.status = status
            super().__init__()

    DEFAULT_CSS = "TestcaseItem {padding: 1;}"

    # These variables are reactive so they can be observed and updated by
    # the Editor widget
    input: var[str] = var("")
    output: var[str] = var("")
    expected_output: var[str] = var("")

    status: reactive[Status] = reactive(Status.INITIAL)

    @property
    def index(self) -> int:
        """Dynamically compute the current index in parent's children."""
        assert self.parent is not None
        return self.parent.children.index(self)

    def render(self) -> RenderResult:
        output = f"Testcase {self.index}"
        if self.status is not Status.INITIAL:
            if self.status is Status.CORRECT:
                color = "$text-success"
            elif self.status is Status.WRONG:
                color = "$text-error"
            else:
                color = "$text-warning"
            output += f" [{color}]({self.status.value})[/]"
        return output

    @on(StateChanged)
    def handle_state_changed(self, update: StateChanged):
        self.output = update.output
        self.status = update.status

    @work(thread=True)
    def run(self, file: Path):
        """
        Runs the given testcase using the specified file and updates the output and status accordingly.
        """
        self.post_message(self.StateChanged("Running...", Status.INITIAL))
        try:
            output = execute(file, self.input)
            status = (
                Status.CORRECT
                if output.split() == self.expected_output.split()
                else Status.WRONG
            )
            self.post_message(self.StateChanged(output, status))
        except CompilationError as e:
            self.post_message(self.StateChanged(str(e), Status.COMPILATION_ERROR))
        except TimeoutExpired:
            self.post_message(
                self.StateChanged("Time Limit Exceeded", Status.TIME_LIMIT_EXCEEDED)
            )
        except CalledProcessError as e:
            self.post_message(
                self.StateChanged(f"{e}\n{e.stderr}", Status.RUNTIME_ERROR)
            )
        except Exception as e:
            self.post_message(
                self.StateChanged(f"Unexpected Error: {e}", Status.UNKNOWN_ERROR)
            )

    def to_json(self) -> dict:
        """
        Converts the TestcaseItem to a JSON object that is able to be stored.
        """
        return {
            "input": self.input,
            "expected_output": self.expected_output,
            "output": self.output,
            "status": self.status.value,
        }

    @staticmethod
    def from_json(data: dict) -> "TestcaseItem":
        """
        Creates a TestcaseItem from the given JSON data.
        """
        testcase = TestcaseItem()
        testcase.set_reactive(TestcaseItem.input, data["input"])
        testcase.set_reactive(TestcaseItem.expected_output, data["expected_output"])
        testcase.set_reactive(TestcaseItem.output, data["output"])
        testcase.set_reactive(TestcaseItem.status, Status(data["status"]))
        return testcase
