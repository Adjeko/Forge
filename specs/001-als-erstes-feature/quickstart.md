# Quickstart: Forge TUI Task Runner (MVP)

## Prerequisites
- Go 1.23+

## Build & Run
```
go build -o forge .
./forge
```

## Try Single Command
1. Focus command list (press `c`)
2. Select `git status` (arrows)
3. Press `Enter` to run
4. Observe colored output & footer progress (single step)

## Try Workflow (after implementation)
1. Define workflow in UI (sequence selection UI TBD)
2. Run with `Enter`
3. Introduce failing command to verify fail-fast

## View Monitors
- Right panel shows monitors updating (e.g., ping 127.0.0.1 vs invalid host)

## Hotkeys
- Press `?` (after help overlay implemented) for legend
- `q` quit

## Tests
```
go test ./...
```
Golden snapshot tests will live under `tests/ui/`.
