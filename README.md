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

## Integration with Zed

You can easily integrate LazyCPH with the [Zed](https://zed.dev) editor for a seamless competitive programming experience. You just need to configure a task and bind it to a shortcut key.

```jsonc
// tasks.json
[
  {
    "label": "lazycph",
    "command": "{path_to_lazycph_executable}",
    "args": ["$ZED_FILE"],
    "use_new_terminal": true,
  }
]
```

```jsonc
// keymap.json
{
  "context": "Workspace",
  "bindings": {
    "alt-g": [
      "task::Spawn",
      { "task_name": "lazycph", "reveal_target": "center" }
    ]
  }
}
```

So now whenever you press your keybind (alt+g in this case), a new terminal window opens with LazyCPH running on the current file. [See how it looks](./assets/zed.mp4).

<video src="https://github.com/user-attachments/assets/a18089c0-594b-4bf6-8053-92a924c2af91"></video>

## 🛠️ Development

Want to contribute? Check out the [CONTRIBUTING.md](CONTRIBUTING.md) for setup instructions and development guidelines.

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 📝 License

This project is open source and available under the MIT License.

## 💡 Inspiration

Built for competitive programmers who want the features of [CPH VSCode extension](https://marketplace.visualstudio.com/items?itemName=DivyanshuAgrawal.competitive-programming-helper) in any IDE or terminal.
