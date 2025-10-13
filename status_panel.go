package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderStatusColumn baut die Zeilen der Status-Spalte (rechte Seite) inkl. Kopf mit Trennmuster und Icons.
func renderStatusColumn(m Model, width int, lines int) []string {
	if lines <= 0 || width <= 0 {
		return nil
	}
	colStyle := lipgloss.NewStyle()
	headStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#E30018"))
	iconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#BBBBBB"))
	contentLines := make([]string, 0, lines)
	// Kopfzeile mit Schrägstrich-Gradient wie andere Header
	statusWord := headStyle.Render("Status")
	rem := width - visibleLen(statusWord) - 2 // zwei Leerzeichen
	if rem < 0 {
		rem = 0
	}
	// Füllmuster aufbauen
	fillerRunes := make([]rune, 0, rem)
	pattern := []rune("///")
	for len(fillerRunes) < rem {
		fillerRunes = append(fillerRunes, pattern...)
	}
	if len(fillerRunes) > rem {
		fillerRunes = fillerRunes[:rem]
	}
	gradient := ""
	if rem > 0 {
		gradient = gradientRunes(fillerRunes, [3]int{181, 0, 19}, [3]int{204, 204, 204})
	}
	headerLine := statusWord
	if rem >= 0 {
		headerLine = statusWord + "  " + gradient
	}
	// Breite sicherstellen (auffüllen falls kleine Differenz durch ANSI)
	if diff := width - visibleLen(headerLine); diff > 0 {
		tail := strings.Repeat("/", diff)
		headerLine += gradientRunes([]rune(tail), [3]int{204, 204, 204}, [3]int{204, 204, 204})
	}
	contentLines = append(contentLines, headerLine)
	contentLines = append(contentLines, padToWidth("", width))
	for _, ic := range m.statusIcons {
		contentLines = append(contentLines, iconStyle.Render(padToWidth(ic+" ", width)))
	}
	for len(contentLines) < lines {
		contentLines = append(contentLines, padToWidth("", width))
	}
	if len(contentLines) > lines {
		contentLines = contentLines[:lines]
	}
	for i, l := range contentLines {
		contentLines[i] = colStyle.Width(width).Render(l)
	}
	return contentLines
}

// padToWidth stellt sicher, dass ein String exakt die gewünschte Rune-Breite hat (abschneiden oder mit Leerzeichen auffüllen).
func padToWidth(s string, w int) string {
	r := []rune(s)
	if len(r) > w {
		return string(r[:w])
	}
	if len(r) < w {
		return s + strings.Repeat(" ", w-len(r))
	}
	return s
}
