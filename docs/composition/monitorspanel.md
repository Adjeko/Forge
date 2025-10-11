# Composition Attempt: MonitorsPanel

Component: MonitorsPanel
Intended Behavior: Grid/column of LED-like indicators reflecting monitor states.
Attempted Charm Primitives:
- lipgloss colored blocks with Unicode circles.
- bubbletea model representing slice of monitor statuses.
Gaps Observed:
- Need jitter timestamp display & dynamic color mapping not in basic examples.
Decision Rationale:
- Custom panel rendering each monitor row with status-specific style.
Parity Considerations (Hotkey vs Zone):
- Zone per LED for future detail on click vs hotkey cycle.
Performance Considerations:
- ≤5 monitors; trivial rendering cost.
Accessibility Considerations:
- Provide text label + status word, not only color.
Next Step:
- Implement in `src/ui/components/monitors_panel.go` after tests.
