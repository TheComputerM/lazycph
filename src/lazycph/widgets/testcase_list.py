from textual.binding import Binding
from textual.reactive import reactive
from textual.widgets import ListView

from lazycph.widgets.testcase_item import TestcaseItem


class TestcaseList(ListView):
    DEFAULT_CSS = "TestcaseList {row-span: 2;}"

    BINDING_GROUP_TITLE = "Testcases"
    BINDINGS = [
        Binding("c", "create", "create", tooltip="create new testcase"),
        Binding("d", "delete", "delete", tooltip="delete selected testcase"),
        Binding("right", "app.focus('input')", "Focus input", show=False),
    ]

    # replace with custom reactive to always update even if the same index is set
    index: reactive[int] = reactive(0, always_update=True, init=False)

    def __init__(self, *children: TestcaseItem) -> None:
        super().__init__(*children, id="testcase-list")

    def action_create(self) -> None:
        """Create a new testcase and add it to the list."""
        self.append(TestcaseItem())
        self.refresh_bindings()

    async def action_delete(self) -> None:
        """Delete the selected testcase from the list."""
        index = self.index
        assert index is not None
        await self.pop(index)
        for item in self.children[index:]:
            item.refresh()  # Update indices
        self.refresh_bindings()

    def check_action(self, action: str, parameters: tuple[object, ...]) -> bool | None:
        if action == "delete" and len(self.children) == 1:
            # Do not delete the last testcase
            return None
        return super().check_action(action, parameters)
