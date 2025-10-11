# Data Model: Forge TUI Task Runner Core

## Entities

### PrimitiveCommand
- Fields: id (string), label (string), command (string), defaultColor (color), lastExitCode (int), lastRunAt (time)
- Invariants: command must be whitelisted.

### Workflow
- Fields: id (string), name (string), steps ([]ExecutionStep), status (enum: Pending|Running|Success|Failed), startedAt (time), endedAt (time), failedStepIndex (int? nullable)
- Transitions: Pending -> Running -> (Success|Failed)

### ExecutionStep
- Fields: workflowId (string), index (int), primitiveId (string), startTime (time), endTime (time), exitCode (int), outputLines ([]LineFragmentRef), color (color), status (enum: Pending|Running|Success|Failed|Hanging)
- Invariants: status Hanging only after >60s no new output while Running.

### StatusMonitor
- Fields: id (string), label (string), type (enum: Ping|Script|Custom), target (string), interval (duration), lastCheck (time), status (enum: Checking|OK|Failed|Disabled), lastError (string?)
- Invariants: interval >= 1s; scheduler jitter ≤500ms optional.

### ProgressBar
- Fields: workflowId (string), total (int), completed (int), lastUpdate (time)
- Derived: percent = completed/total

### HotkeyAction
- Fields: actionId (string), description (string), keys ([]string), zoneId (string), global (bool)
- Invariants: keys non-empty; zoneId registered if not global.

### OutputBuffer
- Fields: lines ([]StyledLine), warningIssued (bool)
- Invariants: append-only; may emit advisory warning after threshold.

### StyledLine
- Fields: text (string), styleTokens ([]Token)

### Token
- Fields: content (string), fg (color?), bg (color?), attr (bitmask?)

## Relationships
- Workflow has many ExecutionSteps.
- ExecutionStep references PrimitiveCommand.
- StatusMonitor independent of Workflow (parallel model tree branch).
- HotkeyAction references BubbleZone ID (string identifier in UI layer).

## Validation Rules
- PrimitiveCommand.command must match whitelist regex list (exact equality for MVP).
- Workflow must have ≥1 step to start.
- Monitor interval default 5s if unspecified.

## Derived / Computed
- Hanging detection timer per ExecutionStep.
- LED color from StatusMonitor.status.
- Palette assignment deterministic (e.g., hash primitiveId mod palette length, with override for errors -> red).

## Notes
- No persistence layer; all state ephemeral.
- Potential future persistence: snapshot of workflows & outputs.
