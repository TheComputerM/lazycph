from pathlib import Path
from typing import Optional

from textual import on, work
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
from textual.worker import get_current_worker

from lazycph.screens.companion import CompanionScreen
from lazycph.screens.file_picker import FilePicker
from lazycph.widgets.editor import Editor


class LazyCPH(App):
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
    TITLE = "LazyCPH"
    SUB_TITLE = "Competitive Programming Helper"
    BINDINGS = [
        Binding("ctrl+f", "choose_file", "select file"),
    ]

    base: Path
    file: reactive[Optional[Path]] = reactive(None, recompose=True, always_update=True)

    companion_mode: bool

    def __init__(
        self,
        base: Path,
        selected: Optional[Path] = None,
        companion: bool = False,
    ):
        super().__init__()
        self.base = base
        self.set_reactive(LazyCPH.file, selected)
        self.companion_mode = companion

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
            yield Editor(file=self.file)
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

    def on_load(self) -> None:
        self.theme = "tokyo-night"

    def on_mount(self) -> None:
        if self.companion_mode:
            self.run_companion_server()
        if self.file is None:
            self.action_choose_file()

    @work(thread=True, exclusive=True)
    def run_companion_server(self) -> None:
        import json
        from http.server import BaseHTTPRequestHandler, HTTPServer

        def set_file(file: Path | None) -> None:
            self.file = file

        app = self

        class CompanionHandler(BaseHTTPRequestHandler):
            def do_POST(self):
                content_length = int(self.headers.get("Content-Length", 0))
                raw = self.rfile.read(content_length)
                data = json.loads(raw)

                app.call_from_thread(
                    app.push_screen, CompanionScreen(data, app.base), set_file
                )

        httpd = HTTPServer(("localhost", 27121), CompanionHandler)
        httpd.timeout = 0.5
        worker = get_current_worker()
        while not worker.is_cancelled:
            httpd.handle_request()
        httpd.server_close()
