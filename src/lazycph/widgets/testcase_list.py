from textual.widgets import ListView

from lazycph.widgets import TestcaseItem


class TestcaseList(ListView):
    DEFAULT_CSS = """
    TestcaseList {
        row-span: 2;
    }
    """

    def __init__(self, *children: TestcaseItem) -> None:
        super().__init__(*children, initial_index=0, id="testcase-list")
