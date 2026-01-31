# AGENTS.md - Coding Agent Guidelines for LazyCPH

## Project Overview

LazyCPH is a Terminal User Interface (TUI) application for competitive programming, built with Python 3.13+ and the Textual framework. It helps competitive programmers test their solutions by running test cases against source files in multiple languages (Python, C, C++, Rust, Zig).

## Build/Lint/Test Commands

### Package Manager
This project uses **uv** (Astral's fast Python package manager). Do not use pip directly.

### Setup
```bash
uv sync --all-groups   # Install all dependencies including dev
```

### Common Tasks
| Command | Description |
|---------|-------------|
| `uv run task start` | Start the application in dev mode |
| `uv run task test` | Run the full test suite |
| `uv run task console` | Run Textual dev console for debugging |
| `uv run task web` | Open the app in web browser mode |
| `uv run task build` | Build standalone binary with PyInstaller |

### Running Tests

**Full test suite:**
```bash
uv run task test
```

**Single test file:**
```bash
uv run pytest tests/test_engines.py -v
```

**Single test class:**
```bash
uv run pytest tests/test_engines.py::TestPython -v
```

**Single test method:**
```bash
uv run pytest tests/test_engines.py::TestPython::test_basic -v
```

### Linting and Formatting
Ruff is used for linting and formatting (uses default configuration).

```bash
uv run ruff check .        # Lint
uv run ruff format .       # Format
uv run ruff check --fix .  # Auto-fix lint issues
```

## Project Structure

```
src/lazycph/
├── __init__.py          # Version info
├── __main__.py          # CLI entry point
├── app.py               # Main Textual App class
├── engines.py           # Language execution engines
├── workspace.py         # Testcase persistence (JSON)
├── screens/             # Textual screens
│   ├── companion.py     # Competitive Companion integration
│   └── file_picker.py   # File selection modal
└── widgets/             # Textual widgets
    ├── editor.py        # Main workspace editor
    ├── testcase_item.py # Individual testcase widget
    └── testcase_list.py # Testcase list management

tests/
└── test_engines.py      # Tests for language engines
```

## Code Style Guidelines

### Formatting
- **Indentation:** 4 spaces
- **Quotes:** Double quotes (`"`) for strings
- **Line length:** ~88 characters (Ruff default)
- **Trailing commas:** Use in multi-line structures

### Import Ordering
Imports are organized in three groups, each alphabetized:

```python
# 1. Standard library imports
import socket
from pathlib import Path
from typing import Optional

# 2. Third-party imports
from textual import on, work
from textual.app import App, ComposeResult
from textual.binding import Binding

# 3. Local application imports
from lazycph.screens.companion import CompanionScreen
from lazycph.widgets.editor import Editor
```

### Type Annotations
Always use type hints for function parameters and return types:

```python
def validate_target_path(path_str: str) -> Path:
    ...

def read_save(source_file: Path) -> Optional[list[dict]]:
    ...
```

- Use `Path` from `pathlib` for file paths (not strings)
- Use modern type syntax: `list[dict]` not `List[Dict]`
- Use `Optional[T]` or `T | None` for nullable types
- Textual reactives should be typed: `reactive[Optional[Path]]`

### Naming Conventions
| Type | Convention | Examples |
|------|------------|----------|
| Classes | PascalCase | `LazyCPH`, `TestcaseItem`, `CompilationError` |
| Functions/Methods | snake_case | `validate_target_path`, `read_save` |
| Variables | snake_case | `source_file`, `save_dir` |
| Constants | SCREAMING_SNAKE_CASE | `DEFAULT_CSS`, `BINDINGS` |
| Private | Leading underscore | `_is_expected_output_correct` |

### Docstrings
Use triple double-quotes with brief, descriptive content:

```python
class CompilationError(Exception):
    """
    Exception raised for errors during the compilation process of runtimes.
    """
```

## Textual Framework Patterns

### Widget Composition
Use the `compose` method to yield child widgets:

```python
def compose(self) -> ComposeResult:
    yield Header()
    yield Container(
        TestcaseList(),
        id="main",
    )
    yield Footer()
```

### CSS Styling
Define CSS as a class attribute:

```python
class MyWidget(Widget):
    DEFAULT_CSS = """
    MyWidget {
        height: auto;
        padding: 1;
    }
    """
```

### Event Handling
Use the `@on` decorator with optional CSS selectors:

```python
@on(Button.Pressed, "#run-button")
def handle_run(self, event: Button.Pressed) -> None:
    ...
```

### Reactive State
- Use `reactive[]` for state that should trigger UI updates
- Use `var[]` for state that doesn't trigger updates

```python
file: reactive[Optional[Path]] = reactive(None)
status: var[Status] = var(Status.PENDING)
```

## Testing Patterns

Tests are organized by language engine. Each test class covers:
- Basic execution
- Timeout handling
- Runtime errors
- Compilation errors (for compiled languages)

```python
class TestPython:
    def test_basic(self):
        with tempfile.NamedTemporaryFile(suffix=".py", delete=False) as f:
            f.write(b"print('hello')")
            f.flush()
            result = execute(Path(f.name), "")
            assert result == "hello\n"

    def test_timeout(self):
        with pytest.raises(TimeoutExpired):
            ...
```

## Key Files to Understand

- `src/lazycph/app.py` - Main application entry and screen management
- `src/lazycph/engines.py` - Language execution engines and registry
- `src/lazycph/widgets/testcase_item.py` - Core test execution logic
- `src/lazycph/workspace.py` - JSON persistence for test cases
