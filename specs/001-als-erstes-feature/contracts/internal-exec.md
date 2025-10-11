# Internal Contract: Execution & Monitoring

## Command Execution Interface
```
Execute(cmd PrimitiveCommand) (ExecutionStep, error)
```
Behavior:
- Spawns subprocess (no shell expansion) with command tokens
- Streams stdout/stderr line fragments to OutputBuffer
- Detects hanging after 60s inactivity -> status Hanging

Errors:
- Non-whitelisted command -> validation error
- Spawn failure -> returns error, step status Failed

## Workflow Orchestrator
```
RunWorkflow(wf *Workflow) error
```
Rules:
- Iterate steps sequentially
- Abort on first Failed unless future continuation flag (not in MVP)
- Update ProgressBar after each step

## Status Monitor Poller
```
Poll(m *StatusMonitor) MonitorResult
```
Return:
- status (OK|Failed)
- latency (duration)
- errorMessage (string?)

Scheduler ensures jitter ≤500ms and non-blocking behavior.

## Hotkey Registry
```
Register(action HotkeyAction) error
Dispatch(key string, focusContext string) (handled bool)
List(context string) []HotkeyAction
```
Guarantees unique (context,key) pairs.
