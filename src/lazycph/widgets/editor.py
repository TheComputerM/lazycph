import threading
from functools import wraps
from pathlib import Path
from typing import Optional

from textual import on
from textual.app import ComposeResult
from textual.binding import Binding
from textual.containers import Grid
from textual.events import Click
from textual.widgets import ListView, TextArea

from lazycph import workspace
from lazycph.widgets.testcase_item import TestcaseItem
from lazycph.widgets.testcase_list import TestcaseList


def debounce(wait_time: float):
    """
    Decorator that will postpone a function's execution until after wait_time seconds
    have elapsed since the last time it was invoked.
    """

    def decorator(function):
        timer: Optional[threading.Timer] = None

        @wraps(function)
        def wrapper(*args, **kwargs):
            nonlocal timer

            # If a timer is already running (meaning a call happened recently),
            # we cancel it. This effectively "resets" the clock.
            if timer is not None:
                timer.cancel()

            def call_function():
                nonlocal timer
                timer = None
                function(*args, **kwargs)

            # Start a new timer for the specified wait_time
            # The function will execute only if this timer completes without being cancelled
            timer = threading.Timer(wait_time, call_function)
            timer.start()

        return wrapper

    return decorator


class Editor(Grid):
    DEFAULT_CSS = """
    Editor {
        grid-size: 3 2;
        grid-columns: 28 1fr 1fr;

        & > TextArea {
            border: thick $boost;
            &:focus {
                border: thick $primary;
            }
        }
    }

    #stdout {
        row-span: 2;
    }
    """

    BINDING_GROUP_TITLE = "Workspace"

    BINDINGS = [
        Binding("ctrl+r", "run", "run", tooltip="run selected testcase"),
        Binding("ctrl+shift+r", "run_all", "run all", tooltip="run all testcases"),
        Binding(
            "escape", "app.focus('testcase-list')", "focus on testcases", show=False
        ),
        Binding(
            "ctrl+left_square_bracket",
            "prev_testcase",
            "select previous testcase",
            show=False,
        ),
        Binding(
            "ctrl+right_square_bracket",
            "next_testcase",
            "select next testcase",
            show=False,
        ),
    ]

    file: Path

    initial_testcases: list[TestcaseItem]

    def __init__(self, file: Path) -> None:
        super().__init__()
        self.file = file

        save_data = workspace.read_save(file)
        if save_data is None:
            # didn't find save file
            self.initial_testcases = [TestcaseItem()]
        else:
            self.initial_testcases = [
                TestcaseItem.from_json(item) for item in save_data
            ]

    def compose(self) -> ComposeResult:
        assert len(self.initial_testcases) > 0
        yield TestcaseList(*self.initial_testcases)
        yield TextArea(
            id="input",
            placeholder="[i]STDIN[/i]",
            show_line_numbers=True,
        )
        yield TextArea(
            id="stdout",
            placeholder="Run [i](^r)[/i] the testcase to see the output",
            show_line_numbers=True,
            read_only=True,
            compact=True,
        )
        yield TextArea(
            id="expected-output",
            placeholder="Expected [i]STDOUT[/i]",
            show_line_numbers=True,
        )

    @property
    def testcase_list(self) -> ListView:
        return self.query_one(TestcaseList)

    @property
    def selected_testcase(self) -> TestcaseItem:
        testcase = self.testcase_list.highlighted_child
        assert isinstance(testcase, TestcaseItem)
        return testcase

    @on(TextArea.Changed, "#input")
    def handle_input_changed(self, event: TextArea.Changed) -> None:
        self.selected_testcase.input = event.control.text

    @on(TextArea.Changed, "#expected-output")
    def handle_expected_output_changed(self, event: TextArea.Changed) -> None:
        self.selected_testcase.expected_output = event.control.text

    @on(Click, "TextArea")
    async def handle_text_area_click(self, event: Click) -> None:
        assert isinstance(event.control, TextArea)
        if event.chain == 2:
            # Select all text in the text area on double click
            await event.control.run_action("select_all")

    def action_run(self) -> None:
        self.selected_testcase.run(self.file)

    def action_run_all(self) -> None:
        for item in self.testcase_list.children:
            assert isinstance(item, TestcaseItem)
            item.run(self.file)

    async def action_prev_testcase(self) -> None:
        await self.testcase_list.run_action("cursor_up")

    async def action_next_testcase(self) -> None:
        await self.testcase_list.run_action("cursor_down")

    @debounce(0.3)
    def action_save_state(self) -> None:
        """
        Save the current state of the workspace (testcases) to a JSON file.
        """
        data = [
            item.to_json()
            for item in self.testcase_list.children
            if isinstance(item, TestcaseItem)
        ]
        workspace.save_file(self.file, data)

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

            # Save workspace state
            self.watch(item, "input", self.action_save_state)
            self.watch(item, "output", self.action_save_state)
            self.watch(item, "expected_output", self.action_save_state)

        self.watch(self.testcase_list, "index", update_selected)
        self.query_one("#stdout", TextArea).can_focus = False
        self.testcase_list.focus()
