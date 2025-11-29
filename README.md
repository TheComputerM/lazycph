# LazyCPH

![Workspace](./assets/workspace.svg)

A beautiful Terminal User Interface (TUI) for competitive programming that helps you test your solutions quickly and efficiently.

## ✨ Features

- 🎨 **Beautiful TUI** - Modern terminal interface built with [Textual](https://github.com/Textualize/textual)
- 🧪 **Multiple Test Cases** - Create, manage, and run multiple test cases for your solution
- ⚡  **Fast** - Edit run your testcases quickly

## 🚀 Installation

### From Releases

```bash
TODO
```

### From Source

1. Clone the repository:
```bash
git clone https://github.com/TheComputerM/lazycph.git
cd lazycph
```

2. Install using uv:
```bash
uv sync --all-groups
uv run task install
# binary is at ./dist/lazycph
```

## 📖 Usage

### Basic Usage

Launch LazyCPH in the current directory:
```bash
lazycph
```

Launch with a specific directory:
```bash
lazycph /path/to/your/code
```

Launch with a specific file:
```bash
lazycph solution.py
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `f` | Open file picker |
| `Ctrl+R` | Run current test case |
| `c` | Create new test case |
| `d` | Delete current test case |
| `↑/↓` | Navigate between test cases |
| `Tab` | Switch between input/output fields |
| `Ctrl+Q` | Quit application |

## 🎮 How It Works

1. **Choose Your File**: Press `f` or click "Choose File" to select your solution file
2. **Create Test Cases**: Press `c` to create new test cases
3. **Add Input**: Type your test input in the STDIN field
4. **Add Expected Output**: Type the expected output in the Expected STDOUT field
5. **Run & Compare**: Press `Ctrl+R` to run your solution and see the output
6. **Navigate**: Use arrow keys to switch between test cases

## 🛠️ Development

Want to contribute? Check out the [CONTRIBUTING.md](CONTRIBUTING.md) for setup instructions and development guidelines.

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 📝 License

This project is open source and available under the MIT License.

## 💡 Inspiration

Built for competitive programmers who want the features of [CPH VSCode extension](https://marketplace.visualstudio.com/items?itemName=DivyanshuAgrawal.competitive-programming-helper) in any IDE or terminal.
