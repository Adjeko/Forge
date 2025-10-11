# Quickstart: Forge TUI Task Runner (MVP)

## Prerequisites
- Go 1.23+

## Build & Run
```
go build -o forge .
./forge
```

## Run Workflow (Current Default)
The default model initializes a single-step workflow using the first whitelisted command.
1. Press `Enter` to execute the workflow
2. Footer shows `exit=0` and progress reaches 100 upon success
3. Any failure would set `exit=-1` and halt progress

## Help Overlay
Press `?` to toggle the help overlay (latency target <150ms). Press `?` again to hide it.
Overlay lists registered hotkey actions and their associated zones for accessibility parity.

## View Monitors
- Right panel shows monitors updating (e.g., ping 127.0.0.1 vs invalid host)

## Hotkeys
| Key    | Action              |
|--------|---------------------|
| Enter  | Run workflow        |
| Tab    | Cycle focus regions |
| PgUp   | Scroll output up    |
| PgDown | Scroll output down  |
| ?      | Toggle help overlay |
| q      | Quit                |

## Tests
Run the full suite (includes perf & accessibility latency checks):
```
go test ./...
```
Golden snapshot tests will live under `tests/ui/`.
