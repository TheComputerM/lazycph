# AGENTS.md

## Overview

LazyCPH is a terminal UI for competitive programming. It runs source files against test cases, shows verdicts, and integrates with the [Competitive Companion](https://github.com/jmerle/competitive-companion) browser extension.

Built with **Go 1.26**, **Bubble Tea v2**, **Bubbles v2**, **Lip Gloss v2**, and **bubblezone v2** for mouse support.

## Build & Run

```sh
go build -o lazycph .        # build
go run .                     # run in cwd (dir → filepicker, file → workspace)
go test ./...                # all tests
go test ./internal/core/...  # engine tests only (needs c++ and python3 in PATH)
```

No Makefile, no CI, no linter config. `go vet` and `go test` are the only verification steps.

## Architecture

```
main.go                          # entry: tea.NewProgram(app.New(cwd))
internal/
  app/model.go                   # root Elm model; routes by path type (dir/file)
  core/
    config.go                    # Engine registry; defaults (.cpp, .py); user overrides from ~/.config/lazycph.json
    engine.go                    # compile/interpret runner; uses {file}/{temp} placeholders
    testcase.go                  # TestCase, TestCaseList, Execute, verdict logic
    store.go                     # JSON persistence in .lazycph/ sibling dir (auto-gitignored)
    messages.go                  # NavigateMsg (path-based screen routing)
  screens/
    filepicker/                  # directory → file selection screen
    workspace/                   # main workspace: test list + input/expected/output panels
    companion/                   # Competitive Companion HTTP listener + file creation dialog
  ui/
    list/                        # test case list component
    textarea/                    # editable input/expected with pointer binding
    output/                      # read-only scrollable output viewport
```

### Key patterns

- **Elm architecture**: every component is a `tea.Model` with `Init/Update/View`. State flows via messages, never direct mutation across boundaries.
- **Screen routing**: `app.Model.active` holds the current screen. `NavigateMsg{Path}` switches screens; empty path = "go back" to state-derived screen.
- **Focus system**: workspace has 4 focusable panels (list, input, expected, output) cycled with tab/shift-tab. Each implements `Focus()/Blur()` interface.
- **Mouse zones**: `bubblezone` marks regions; `tea.MouseReleaseMsg` dispatches focus/clicks by zone ID.
- **File convention per component**: `model.go` (state + Update), `keymap.go` (bindings), `style.go` (lipgloss styles), `layout.go` (size calculations).

### Companion integration

`companion/server.go` listens on well-known ports (1327, 4244, 6174, 10042, 10043, 10045, 27121) for POST payloads from the Competitive Companion browser extension. Received data triggers the companion screen to create source files with pre-populated test cases.

### Test data persistence

Test cases are stored as JSON in `.lazycph/<filename-without-ext>.json` next to the source file. The directory is auto-created with a `.gitignore` containing `*`.
