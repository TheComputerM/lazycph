from pathlib import Path

from textual.app import ComposeResult
from textual.screen import Screen
from textual.widgets import Header, Label, ListItem, ListView

from lazycph import workspace
from lazycph.engines import engines

_AVAILABLE_ENGINES = list(engines.keys())


class CompanionScreen(Screen[Path]):
    TITLE = "Companion Mode"
    SUB_TITLE = "Select Runtime"

    def __init__(self, data: dict, base: Path) -> None:
        super().__init__()
        self.data = data
        self.base = base

    def compose(self) -> ComposeResult:
        yield Header()
        yield ListView(*[ListItem(Label(name)) for name in _AVAILABLE_ENGINES])

    def on_list_view_selected(self, event: ListView.Selected) -> None:
        event.stop()
        extension = _AVAILABLE_ENGINES[event.index]

        group = self.base.joinpath(self.data["group"])
        group.mkdir(exist_ok=True)

        file = group.joinpath(self.data["name"]).with_suffix(extension)

        testcases = [
            {
                "input": item["input"],
                "expected_output": item["output"],
                "output": "",
                "status": None,
            }
            for item in self.data["tests"]
        ]

        workspace.save_file(file, testcases)

        file.write_text("")

        self.dismiss(file)
