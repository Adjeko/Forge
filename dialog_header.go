package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderDialogHeader erstellt eine einzeilige Kopfzeile für Dialoge.
// title steht links in SEW-Rot fett, gefolgt von zwei Leerzeichen, danach wird mit
// dem Muster "///" bis zur gewünschten Breite aufgefüllt. Die Slashes erhalten
// denselben Farbverlauf wie das FORGE Logo (Start #B50013 -> Ende #CCCCCC).
func RenderDialogHeader(title string, width int) string {
	if width <= 0 {
		width = len(title) + 2
	}
	if width < len(title)+2 {
		width = len(title) + 2
	}
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E30018")).Bold(true)
	left := titleStyle.Render(title)
	remaining := width - visibleLen(left) - 2 // zwei Leerzeichen
	if remaining < 0 {
		remaining = 0
	}
	// Roh-String aus Slashes als Füller aufbauen
	fillerRunes := make([]rune, 0, remaining)
	pattern := []rune("///")
	for len(fillerRunes) < remaining {
		fillerRunes = append(fillerRunes, pattern...)
	}
	if len(fillerRunes) > remaining {
		fillerRunes = fillerRunes[:remaining]
	}
	// Apply gradient across filler runes
	gradient := gradientRunes(fillerRunes, [3]int{181, 0, 19}, [3]int{204, 204, 204})
	out := left + "  " + gradient
	// Falls durch Rundung / ANSI-Längen Differenz besteht -> mit letzten Farbton auffüllen
	if diff := width - visibleLen(out); diff > 0 {
		tail := strings.Repeat("/", diff)
		tailColored := gradientRunes([]rune(tail), [3]int{204, 204, 204}, [3]int{204, 204, 204})
		out += tailColored
	}
	return out
}

// visibleLen (rough) – wir gehen davon aus, dass lipgloss keine Zero-Width Runen injiziert außer ANSI.
func visibleLen(s string) int { return len([]rune(stripANSI(s))) }

// gradientRunes färbt jeden Rune einzeln.
func gradientRunes(runes []rune, startRGB, endRGB [3]int) string {
	n := len(runes)
	if n == 0 {
		return ""
	}
	if n == 1 {
		hex := fmt.Sprintf("#%02X%02X%02X", startRGB[0], startRGB[1], startRGB[2])
		return lipgloss.NewStyle().Foreground(lipgloss.Color(hex)).Render(string(runes[0]))
	}
	var b strings.Builder
	b.Grow(n * 10)
	for i, r := range runes {
		t := float64(i) / float64(n-1)
		cr := int(float64(startRGB[0])*(1-t) + float64(endRGB[0])*t)
		cg := int(float64(startRGB[1])*(1-t) + float64(endRGB[1])*t)
		cb := int(float64(startRGB[2])*(1-t) + float64(endRGB[2])*t)
		hex := fmt.Sprintf("#%02X%02X%02X", cr, cg, cb)
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(hex)).Bold(false).Render(string(r)))
	}
	return b.String()
}
