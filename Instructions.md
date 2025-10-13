# Instructions

These guidelines describe how to extend and maintain this TUI application consistently.

## Purpose and scope
- Application type: Go TUI using Bubble Tea + Lip Gloss; viewport for scrollable terminal output.
- Primary features:
  - Title bar with FORGE art and version label.
  - Middle area showing terminal output, scrollable.
  - Search line (input) at bottom; pressing Enter runs a shell command.

## Architecture overview
- `main.go`: Program entry; creates the Bubble Tea program with alt screen.
- `model.go`: Core state machine (Bubble Tea model). Handles keys, window size, running shell commands, and manages viewport content.
- `dashboard.go`: Layout renderer for the middle section (output area + input line). Uses the model’s viewport for output.
- `titlebar.go`: Title bar art rendering, gradients, and version label.
- `go.mod`: Module and dependencies.

## UI/UX guidelines
- UI-Framework-Vorgabe: UI-/Oberflächen-Implementierungen sollen, soweit möglich,
  mit Charm Bubble Tea (State/Update/View), Bubbles (UI-Komponenten wie `viewport`)
  und Lip Gloss (Styling) umgesetzt werden. Keine Parallel-Frameworks mischen.
- No frame/borders around the output area.
- Output area is a scrollable viewport:
  - Arrow keys, PageUp/PageDown, and mouse wheel scroll.
  - On new output, auto-scroll to bottom.
- Input line:
  - Shows a light “>> ” label and caret.
  - Does not echo the submitted text into the output area.
- Footer shows concise hotkey hints.
- Keep all user-facing text in German.

## Input behavior and hotkeys
- `Enter` in the search line executes a command (see policy below) and clears the input.
- `q` or `Ctrl+C`: exit.
- `?` or `h`: toggle help.
- `ESC`: closes help or modal.
- `m`: toggle modal (if used).

## Command execution policy (shell, not client libraries)
- Always execute system commands via `os/exec` (no Git client SDKs).
- For the current use case, `Enter` runs `git status` in `C:\ADO\CS2`:
  - Use `exec.Command("git", "status")` and set `cmd.Dir = "C:/ADO/CS2"`.
  - Capture combined stdout/stderr with `CombinedOutput()`.
  - Normalize line endings (`\r\n` → `\n`) and split into lines.
  - Append lines to `termOutput` and refresh the viewport.
- Display command errors as a simple `[error] ...` line appended to the output.

## Dependencies
- Core: `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/lipgloss`.
- Viewport: `github.com/charmbracelet/bubbles/viewport` (module `github.com/charmbracelet/bubbles`).
- Keep `go.mod` tidy and pinned to known-good versions.

Optional (for maintainers):
```pwsh
# Add bubbles (provides viewport)
go get github.com/charmbracelet/bubbles@v0.18.0
# Clean up
go mod tidy
```

## Coding conventions
- Go ≥ the version declared in `go.mod`.
- Avoid global state; keep state in `Model`.
- Keep rendering pure (derive from model state) and avoid side effects in `View()`.
- Isolate shell execution into small functions returning Bubble Tea messages.
- Keep styles centralized (e.g., shared `lipgloss` styles in `model.go`).

## Internationalization
- Maintain German text for UI labels and help.
- Keep user-facing text centralized where reasonable for easier future i18n.

## Versioning
- The title bar shows a version like `v0.0.1` (see `titlebar.go`).
- When bumping versions, update the constant/string there and changelog notes (if added later).

## Testing strategy
- Favor unit tests for pure helpers:
  - `padToLines`, `insertAt`, `flattenOutput`, `compress5to3`, `normalizeLetter`.
- Table-driven tests with clear inputs/outputs.
- For visual functions (gradients), test length and non-empty content rather than exact styling.
- For command execution, prefer wrapping `exec.Command` behind a small function to make it mockable if needed.

## Error handling and diagnostics
- Surface command errors inline in the output area (`[error] ...`).
- Avoid panics in runtime paths; fail fast only at startup if configuration is invalid.

## Contribution workflow
- Branch naming: `feat/...`, `fix/...`, `chore/...`.
- Keep PRs focused and small; include a short summary of UI/behavior changes.
- When changing user-visible behavior, update these Instructions if needed.

## Release checklist
- Update version label in the title bar.
- Ensure dependencies are tidy and pinned.
- Verify key workflows: window resize, scrolling, input handling, Enter command.

## Assistant usage policy
- The assistant may edit files, but terminal actions are opt-in via your request.

## Ideas / backlog
- Configurable command directory (instead of hardcoded `C:\ADO\CS2`).
- Command history/navigation for the input line.
- Additional commands (e.g., `git log -n 10`, `git branch`) with simple routing based on input.
- Git Status Dialog: ^+G öffnet einen Dialog zur Eingabe eines Repository-Pfads (Validierung auf .git).
- Persist output history to a log file.
- Config file (YAML/TOML) for keybindings, styles, and command directory.
