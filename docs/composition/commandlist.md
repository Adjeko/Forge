# Composition Attempt: CommandList

Component: CommandList
Intended Behavior: Display whitelist of primitive commands, allow selection & execution trigger.
Attempted Charm Primitives:
- bubbletea list model (considered) but need custom styling for color assignment preview.
- lipgloss nested styles for selected vs unselected row.
- bubblezone: Zone per command row for click support.
Gaps Observed:
- Standard list lacks integrated color swatch & zone registration convenience.
Decision Rationale:
- Custom lightweight list wrapper around slice + selected index; integrate BubbleZone IDs.
Parity Considerations (Hotkey vs Zone):
- Hotkey navigation (arrows/enter) vs mouse click zone activation.
Performance Considerations:
- O(n) render over ≤10 commands acceptable.
Accessibility Considerations:
- Focus ring style & clear selected highlight.
Next Step:
- Implement in `src/ui/components/command_list.go` after tests.
