# LazyCPH

A terminal UI for competitive programming.

![LazyCPH Demo](./assets/lazycph.gif)

---

## Features

- Run test cases against source files
- Competitive Companion integration for automatic test case retrieval
- Mouse support for panel selection and scrolling
- Configurable language engines for compilation and execution

## Installation

```sh
go install github.com/thecomputerm/lazycph@main
```

## Usage

Run in the current directory to open the file picker:
```sh
lazycph
```

Open a specific solution file:
```sh
lazycph path/to/solution.cpp
```

Enable Competitive Companion integration:
```sh
lazycph --companion
```

## Configuration

Custom engines can be defined in `~/.config/lazycph.json`. This file merges with the default configuration. Use `{file}` as a placeholder for the source file and `{temp}` for the compiled binary path.

```json
{
  "engines": {
    ".cpp": {
      "mode": "compile",
      "command": ["g++", "{file}", "-o", "{temp}", "-O3", "-std=c++20"]
    },
    ".js": {
      "mode": "interpret",
      "command": ["node", "{file}"]
    }
  }
}
```

## Competitive Companion

LazyCPH integrates with the [Competitive Companion](https://github.com/jmerle/competitive-companion) browser extension. When running with the `--companion` or `-c` flag, LazyCPH listens on standard ports for problem data. Sending a problem from the browser will trigger a dialog to create the source file and automatically populate it with the provided test cases.

## Integration with Zed

You can easily integrate LazyCPH with the [Zed](https://zed.dev) editor for a seamless competitive programming experience. You just need to configure a task and bind it to a shortcut key.

```jsonc
// ~/.config/zed/tasks.json
[
  {
    "label": "lazycph",
    "command": "lazycph",
    "args": ["$ZED_FILE"],
    "use_new_terminal": true,
  }
]
```

```jsonc
// ~/config/zed/keymap.json
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

## Building from source

```sh
go build -o lazycph .
```

## License

MIT
