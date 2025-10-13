# CRUSH.md

## Build, Lint, and Test Commands
- Build: `go build .`
- Lint: Use `golangci-lint run` if installed (not pinned)
- Run all tests: `go test ./...`
- Run a single test file: `go test -v -run TestName ./path/to/file.go`

## Code Style Guidelines
- Language: Go (`go.mod` specifies version, Bubble Tea stack)
- Für TUI-Implementierung bevorzugt die Charmland-Libraries Bubbletea, Bubbles, Lipgloss, Glamour, bubblezone, huh etc. verwenden.
- Imports: Standard library first, external grouped, sorted alphabetically within groups
- Formatting: Use `gofmt` before committing
- Types: Prefer explicit typing for structs, interfaces, and method receivers
- Naming:
  - Branches: `feat/...`, `fix/...`, `chore/...`
  - Variables/camelCase except exported (PascalCase)
  - Functions: simple, imperative
- Error Handling:
  - Use `error` returns, surface errors in terminal as `[error] ...`
  - Avoid panics in runtime; only panic at startup for invalid config
- Never use uninitialized globals; state goes in `Model`
- UI pure rendering: `View()` should be side-effect free
- All user-facing text in German; centralize for easy i18n
- Keep styles in `model.go` when possible
- Tests: Prefer table-driven for helpers, mock shell execution
- Use pinned, tidy dependencies (`go mod tidy`)

## Assistant Usage Policy
- Only perform file edits, commands on explicit user request.

## Architektur
- Befehle sind in Go implementierte, wiederverwendbare Bausteine, die Nutzer zu eigenen Abläufen (Flows) kombinieren.

## Qualitätsregeln
- Vor Abschluss einer Aufgabe müssen alle LSP/Compiler Fehler behoben sein (keine roten Diagnostics im Projekt).

