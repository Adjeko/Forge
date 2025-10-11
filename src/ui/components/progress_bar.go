package components

import "github.com/charmbracelet/lipgloss"

// ProgressBar simple placeholder
func ProgressBar(percent int) string {
	width := 20
	filled := (percent * width) / 100
	if filled > width {
		filled = width
	}
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "#"
		} else {
			bar += "-"
		}
	}
	return style.Render(bar)
}
