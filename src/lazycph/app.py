from pathlib import Path
from typing import Optional

from textual import on
from textual.app import App, ComposeResult
from textual.binding import Binding
from textual.containers import CenterMiddle
from textual.reactive import reactive
from textual.widgets import (
    Button,
    Footer,
    Header,
    Label,
)

from lazycph.screens import FilePicker
from lazycph.widgets.workspace import Workspace


class LazyCPH(App):
    TITLE = "LazyCPH"
    SUB_TITLE = "Competitive Programming Helper"
    DEFAULT_CSS = """
    #main {
        grid-size: 2;
        grid-columns: 28 1fr;
    }

    #btn-choose-file {
        width: 100%;
        padding: 1 0;
    }
    """
    BINDINGS = [
        Binding("ctrl+f", "choose_file", "select file"),
    ]

    base: Path
    file: reactive[Optional[Path]] = reactive(None, recompose=True)

    def __init__(
        self,
        base: Path,
        selected: Optional[Path] = None,
    ):
        super().__init__()
        self.base = base
        self.set_reactive(LazyCPH.file, selected)

    def compose(self) -> ComposeResult:
        yield Header()
        _btn = Button(
            label=str(self.file.absolute()) if self.file else "Choose File",
            compact=True,
            id="btn-choose-file",
        )
        _btn.can_focus = False
        yield _btn
        if self.file:
            yield Workspace(file=self.file)
        else:
            # When no file is chosen, show a message in the center
            yield CenterMiddle(Label("Select your file (^f) to begin."))

        yield Footer()

    @on(Button.Pressed, "#btn-choose-file")
    def handle_choose_file(self, _: Button.Pressed) -> None:
        self.action_choose_file()

    def action_choose_file(self) -> None:
        def set_file(file: Path | None) -> None:
            self.file = file

        self.push_screen(FilePicker(self.base), set_file)

    def on_mount(self) -> None:
        self.theme = "tokyo-night"
        if self.file is None:
            self.action_choose_file()
