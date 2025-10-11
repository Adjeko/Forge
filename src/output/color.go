package output

import "github.com/charmbracelet/lipgloss"

var palette = []lipgloss.Color{
	lipgloss.Color("63"),  // blue
	lipgloss.Color("64"),  // green
	lipgloss.Color("65"),  // teal
	lipgloss.Color("99"),  // purple
	lipgloss.Color("203"), // orange/red
	lipgloss.Color("214"), // gold
}

// ColorForID returns a deterministic color from the palette based on id hash.
func ColorForID(id string) lipgloss.Color {
	if len(palette) == 0 {
		return lipgloss.Color("15")
	}
	// djb2 variant with final xor to reduce clustering on similar suffixes
	var h uint32 = 5381
	for i := 0; i < len(id); i++ {
		h = ((h << 5) + h) + uint32(id[i]) // h*33 + c
	}
	h ^= (h >> 13)
	idx := int(h) % len(palette)
	return palette[idx]
}

// ErrorColor overrides normal color for errors.
var ErrorColor = lipgloss.Color("196")

// Collision Handling Strategy (T043): If caller detects adjacent duplicate colors,
// it can shift the later one by +1 modulo palette length. We keep core hash pure
// and expose a helper for workflow steps.

// ColorForStep returns the color for a step; if error is true returns ErrorColor.
// Palette Exhaustion (T044): When steps > len(palette), modulo reuse occurs by design.
// This ensures deterministic repetition without expansion of palette.
func ColorForStep(id string, error bool) lipgloss.Color {
	if error {
		return ErrorColor
	}
	return ColorForID(id)
}
