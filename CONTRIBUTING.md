## Getting Started

### Prerequisites

- Python 3.13 or higher
- [uv](https://github.com/astral-sh/uv) package manager
- Git

1. Fork and clone the repository:

  ```bash
  git clone https://github.com/YOUR_USERNAME/lazycph.git
  cd lazycph
  ```

2. Install dependencies using uv:

  ```bash
  uv sync --all-groups
  ```

3. Verify installation:

  ```bash
  uv run task test
  ```

### Scripts

We use taskipy for writing and running scripts.

- `uv run task start`: Start the application.
- `uv run task console`: Run the textual dev console.
- `uv run task web`: Open the app on the web.
- `uv run task test`: Run the test suite.
- `uv run task build`: Build the app into a single binary using pyinstaller.

### Releasing

Make the version changes to pyproject.toml and lazycph/__init__.py

```sh
git tag v{VERSION}
git push origin v{VERSION}
```
