# Composition Attempt: OutputViewport

Component: OutputViewport
Intended Behavior: Scrollable region showing colored output lines with high retention.
Attempted Charm Primitives:
- bubbletea default text area not suitable for large dynamic streaming lines.
- lipgloss styling for borders & header line.
- bubblezone: zones may wrap scroll buttons later.
Gaps Observed:
- Need efficient slicing of large buffer without re-concatenating full text each frame.
Decision Rationale:
- Custom viewport logic referencing OutputBuffer lines slice + offset.
Parity Considerations (Hotkey vs Zone):
- PgUp/PgDn hotkeys vs future clickable scroll indicators.
Performance Considerations:
- Avoid building one giant string; incremental join of visible window only.
Accessibility Considerations:
- Focus state influences scroll hotkey capture.
Next Step:
- Implement in `src/ui/components/output_view.go` after tests.
