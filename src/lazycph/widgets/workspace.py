from pathlib import Path

from textual import on
from textual.app import ComposeResult
from textual.binding import Binding
from textual.containers import Grid
from textual.widgets import ListView, TextArea

from lazycph.widgets import TestcaseItem
from lazycph.widgets.testcase_list import TestcaseList


class Workspace(Grid):
    DEFAULT_CSS = """
    Workspace {
        grid-size: 3 2;
        grid-columns: 28 1fr 1fr;
    }

    Workspace > TextArea {
        border: thick $boost;
    }

    Workspace > TextArea:focus {
        border: thick $primary;
    }

    #stdout {
        row-span: 2;
    }
    """

    BINDINGS = [
        Binding("ctrl+r", "run", "run"),
        Binding("ctrl+shift+r", "run_all", "run all"),
        Binding("escape", "app.focus('testcase-list')", "Focus list", show=False),
    ]

    file: Path

    def __init__(self, file: Path) -> None:
        super().__init__()
        self.file = file

    def compose(self) -> ComposeResult:
        with TestcaseList():
            yield TestcaseItem()
        yield TextArea(
            id="input",
            placeholder="STDIN",
            show_line_numbers=True,
        )
        yield TextArea(
            id="stdout",
            placeholder="Run (^r) the testcase to see the output",
            show_line_numbers=True,
            read_only=True,
            compact=True,
        )
        yield TextArea(
            id="expected-output",
            placeholder="Expected STDOUT",
            show_line_numbers=True,
        )

    @property
    def testcase_list(self) -> ListView:
        return self.query_one(TestcaseList)

    @property
    def selected_testcase(self) -> TestcaseItem:
        selected_index = self.testcase_list.index
        assert selected_index is not None
        testcase = self.testcase_list.children[selected_index]
        assert isinstance(testcase, TestcaseItem)
        return testcase

    @on(TextArea.Changed, "#input")
    def handle_input_changed(self, event: TextArea.Changed) -> None:
        self.selected_testcase.input = event.control.text

    @on(TextArea.Changed, "#expected-output")
    def handle_expected_output_changed(self, event: TextArea.Changed) -> None:
        self.selected_testcase.expected_output = event.control.text

    async def action_run(self) -> None:
        self.selected_testcase.run(self.file)

    async def action_run_all(self) -> None:
        for item in self.testcase_list.children:
            assert isinstance(item, TestcaseItem)
            item.run(self.file)

    def on_mount(self) -> None:
        def update_output(output: str):
            self.query_one("#stdout", TextArea).text = output

        def update_selected(index: int):
            """
            Update textareas when selected testcase changes
            """
            item = self.testcase_list.children[index]
            assert isinstance(item, TestcaseItem)
            self.query_one("#input", TextArea).text = item.input
            self.query_one("#expected-output", TextArea).text = item.expected_output
            self.query_one("#stdout", TextArea).text = item.output
            self.watch(item, "output", update_output)

        self.watch(self.testcase_list, "index", update_selected)
        self.query_one("#stdout", TextArea).can_focus = False
