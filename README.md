# LazyCPH

A terminal UI for competitive programming.

## Features

- Run test cases against source files
- Edit input and expected output inline
- Real-time verdicts: AC (Accepted), WA (Wrong Answer), TLE (Time Limit Exceeded), RE (Runtime Error)
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

## Supported Languages

| Extension | Language | Mode |
|-----------|----------|------|
| .c | C (C17) | compile |
| .cpp | C++ (C++17) | compile |
| .go | Go | compile |
| .py | Python 3 | interpret |
| .rs | Rust | compile |
| .zig | Zig | compile |

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

## Building from source

```sh
go build -o lazycph .
```

## License

MIT
