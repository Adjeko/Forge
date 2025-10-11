# Research: Forge TUI Task Runner Core Interface

**Feature**: 001-als-erstes-feature  
**Date**: 2025-10-10  

## R1 Buffer Rendering Strategy

- Decision: Use incremental diff rendering with Bubbletea's Update/View cycle; maintain ring buffer of line metadata + virtual viewport.
- Rationale: Avoid full string concatenation for large outputs; keeps allocations bounded while enabling scrollback.
- Alternatives Considered:
  - Full in-memory concatenated string (simple, but expensive for >50k lines)
  - External pager integration (adds complexity; breaks integrated color-coding)

Planned Implementation Note: `output/buffer.go` will store slices of styled line fragments; viewport rendering only maps visible range.

## R2 BubbleZone + Focus Traversal

- Decision: Central focus manager enumerating logical regions (CommandList, Output, Monitors, HelpOverlay). BubbleZone registers per element; Tab cycles list; Shift+Tab reverse.
- Rationale: Central registry enforces parity and discoverability; simplifies help overlay generation (iterate actions registry).
- Alternatives:
  - Decentralized per-model focus state (risk of inconsistent navigation)
  - Mouse-only zone activation fallback (violates Principle VI)

Help Overlay Generation: Build from hotkey action registry (struct with Action, Keys, ZoneID, Description).

## R3 Structured Logging Library

- Decision: Use Go `slog` (stdlib) with custom text handler tuned for minimal flicker; later optional JSON handler for debugging.
- Rationale: Reduces dependencies, future-proofs with standard API, sufficient performance.
- Alternatives:
  - `zerolog` (very fast, but extra dependency not yet justified)
  - `logrus` (legacy, slower formatting, unnecessary)

Logged Events (initial): `command.start`, `command.end`, `workflow.start`, `workflow.end`, `monitor.poll.start`, `monitor.poll.result`, `ui.focus.change`, `ui.help.toggle`.

## R4 Memory Mitigation Strategy

- Decision: MVP accepts unlimited logical buffer but implements soft advisory threshold (e.g., warn after 100k lines via footer indicator) without truncation.
- Rationale: Keeps scope minimal while providing user feedback for extreme sessions.
- Alternatives:
  - Hard truncation oldest lines (risks losing diagnostic history) 
  - Configurable max lines (future enhancement)

Future Hook: Expose optional `--max-lines` flag once configuration system exists.

## Additional Outcome: New Functional Requirement Proposal

Add FR-022: System MUST provide a help overlay toggled by `?` listing all active hotkeys and descriptions, appearing/disappearing within ≤150 ms.

## Open Questions (Deferred)

- Cancellation semantics (Esc) for workflows: not in MVP.
- Monitor detail pane on click: post-MVP.

## Summary Table

| Research ID | Decision Summary | Status |
|-------------|------------------|--------|
| R1 | Incremental diff + virtual viewport | Final |
| R2 | Central focus + action registry | Final |
| R3 | Use slog for structured logging | Final |
| R4 | Unlimited buffer + advisory warning | Final |
| FR-022 | Add help overlay requirement | Pending spec update |
