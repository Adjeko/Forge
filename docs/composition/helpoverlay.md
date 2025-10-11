# Composition Attempt: HelpOverlay

Component: HelpOverlay
Intended Behavior: Modal overlay listing all hotkeys with descriptions.
Attempted Charm Primitives:
- lipgloss box styling with border + title.
- bubbletea conditional view layering.
Gaps Observed:
- Need auto-generation from registry; stock examples static.
Decision Rationale:
- Custom overlay assembling table from actions registry slice.
Parity Considerations (Hotkey vs Zone):
- Toggle via '?' hotkey; potential future clickable close zone.
Performance Considerations:
- Generated on toggle; small list.
Accessibility Considerations:
- Focus moves into overlay; ESC or '?' hides.
Next Step:
- Implement in `src/ui/help/overlay.go` after tests.
