from pathlib import Path
from typing import Iterable

from textual.app import ComposeResult
from textual.reactive import var
from textual.screen import Screen
from textual.widgets import DirectoryTree, Footer, Header, Input

from lazycph.engines import engines


class SourceTree(DirectoryTree):
    search: var[str] = var("")

    DEFAULT_CSS = """
    SourceTree {
        padding: 1 1;
        margin: 0 1;
    }
    """

    async def watch_search(self, _: str) -> None:
        self.run_worker(self.reload(), name="reload_tree_search")

    def filter_paths(self, paths: Iterable[Path]) -> Iterable[Path]:
        return [
            path
            for path in paths
            if (
                path.is_dir()
                or (
                    path.suffix in engines.keys()
                    and self.search.lower() in path.name.lower()
                )
            )
        ]


class FilePicker(Screen[Path]):
    TITLE = "Select a file"
    BINDINGS = [("escape", "app.pop_screen", "Close file picker")]

    base: Path
    search: var[str] = var("")

    def __init__(self, base: Path) -> None:
        super().__init__()
        self.base = base

    def compose(self) -> ComposeResult:
        yield Header()
        yield Input(placeholder="Search for file...")
        yield SourceTree(self.base).data_bind(search=FilePicker.search)
        yield Footer()

    def on_directory_tree_file_selected(self, event: SourceTree.FileSelected) -> None:
        self.dismiss(event.path)

    def on_input_changed(self, event: Input.Changed) -> None:
        self.search = event.value
