from pathlib import Path

from textual.app import ComposeResult
from textual.containers import CenterMiddle
from textual.screen import Screen
from textual.widgets import Header, Label, ListItem, ListView

from lazycph import workspace
from lazycph.engines import available

_VALID_EXTENSIONS = list(available.keys())


class CompanionScreen(Screen[Path]):
    DEFAULT_CSS = """
    CompanionScreen {
        align: center middle;

        & > CenterMiddle {
            width: 60;
            height: auto;
            padding: 1;
            background: $surface;
            background-tint: $foreground 5%;

            border: tab $secondary;
            border-title-align: center;

            & > ListView {
                height: auto;
            }
        }
    }
    """

    TITLE = "Companion Mode"
    SUB_TITLE = "Select Language"

    BINDINGS = [("escape", "app.pop_screen", "Close modal")]

    def __init__(self, data: dict, base: Path) -> None:
        super().__init__()
        self.data = data
        self.base = base

    def compose(self) -> ComposeResult:
        yield Header()
        with CenterMiddle() as container:
            container.border_title = self.data["group"]
            yield ListView(
                *[
                    ListItem(Label(f"{self.data['name']}{suffix}"))
                    for suffix in _VALID_EXTENSIONS
                ]
            )

    def on_list_view_selected(self, event: ListView.Selected) -> None:
        event.stop()
        suffix = _VALID_EXTENSIONS[event.index]

        group = (
            self.base
            if self.base.name == self.data["group"]
            else self.base.joinpath(self.data["group"])
        )
        group.mkdir(exist_ok=True)

        file = group.joinpath(f"{self.data['name']}{suffix}")

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

    def on_mount(self) -> None:
        self.query_one(ListView).focus()
