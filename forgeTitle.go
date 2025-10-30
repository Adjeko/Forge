package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Stile für die Titelzeile.
var (
	titleTextStyle = lipgloss.NewStyle().Bold(true) // nur fett; Farbe über Gradient
	SlashColor     = lipgloss.Color("#9A0010")      // exportiert für Wiederverwendung (Modal-Rand, etc.)
)

// RenderTitleBar rendert die Titelzeile über die angegebene Breite.
// Enthält eine zusätzliche letzte Zeile mit den Seitennamen (Dashboard / Pull Requests).
func RenderTitleBar(w int, currentPage string) string {
	if w <= 0 {
		w = 80
	}
	// FORGE als 3 Zeilen aus Halbblöcken (▀ oben, ▄ unten, █ voll) aufbauen.
	// Ausgangspunkt: 5-Zeilen-Glyphen zu 3 Zeilen komprimieren.
	F5 := []string{
		"████████",
		"█      ",
		"███████ ",
		"█      ",
		"█      ",
	}
	O5 := []string{
		" █████ ",
		"█     █",
		"█     █",
		"█     █",
		" █████ ",
	}
	R5 := []string{
		"███████ ",
		"█     █",
		"███████ ",
		"█    █ ",
		"█     █",
	}
	G5 := []string{
		" █████ ",
		"█      ",
		"█  ███ ",
		"█    █ ",
		" █████ ",
	}
	E5 := []string{
		"████████",
		"█      ",
		"███████ ",
		"█      ",
		"████████",
	}

	F3 := normalizeLetter(compress5to3(F5))
	O3 := normalizeLetter(compress5to3(O5))
	R3 := normalizeLetter(compress5to3(R5))
	G3 := normalizeLetter(compress5to3(G5))
	E3 := normalizeLetter(compress5to3(E5))

	spacer := " "
	// Drei FORGE-Reihen (roh, ohne Farbe) vorberechnen um Breite zu erhalten.
	forgeRows := make([]string, 3)
	artWidth := 0
	for i := 0; i < 3; i++ {
		r := F3[i] + spacer + O3[i] + spacer + R3[i] + spacer + G3[i] + spacer + E3[i]
		forgeRows[i] = r
		if l := len([]rune(r)); l > artWidth {
			artWidth = l
		}
	}

	lines := make([]string, 0, 5)
	// Prefix-Muster "/" mit Breite 8
	prefix := strings.Repeat("/", 8)
	prefixLen := len([]rune(prefix))
	slashColor := SlashColor
	prefixColored := titleTextStyle.Foreground(SlashColor).Render(prefix)
	beforePad := " "
	afterPad := " "
	beforePadLen := len([]rune(beforePad))
	afterPadLen := len([]rune(afterPad))

	// Rote "SEW" Zeile über FORGE; links ausgerichtet sodass 'S' mit 'F' bündig ist.
	red := lipgloss.Color("#E30018")
	sewPlain := "SEW"
	sewColored := titleTextStyle.Foreground(red).Render(sewPlain)
	verPlain := " v0.0.1" // führendes Leerzeichen
	verColored := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(verPlain)
	gapBetween := artWidth - len([]rune(sewPlain)) - len([]rune(verPlain))
	if gapBetween < 0 {
		gapBetween = 0
	}
	sewLine := prefixColored + beforePad + sewColored + strings.Repeat(" ", gapBetween) + verColored + afterPad
	if fill := w - (prefixLen + beforePadLen + len([]rune(sewPlain)) + gapBetween + len([]rune(verPlain)) + afterPadLen); fill > 0 {
		sewLine += titleTextStyle.Foreground(slashColor).Render(strings.Repeat("/", fill))
	}
	lines = append(lines, sewLine)

	for i := 0; i < 3; i++ {
		row := forgeRows[i]
		colored := gradientText(row, [3]int{181, 0, 19}, [3]int{204, 204, 204})
		line := prefixColored + beforePad + colored + afterPad
		if fill := w - (prefixLen + beforePadLen + len([]rune(row)) + afterPadLen); fill > 0 {
			line += titleTextStyle.Foreground(SlashColor).Render(strings.Repeat("/", fill))
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// RenderForgeLines liefert ausschließlich die drei FORGE-Zeilen (ohne SEW/Version/Seitenzeile) mit festem Einzug (10 Leerzeichen).
// Dies dient der Wiederverwendung im Header ohne Zusatzlogik.
func RenderForgeLines(w int) []string {
	if w <= 0 {
		w = 80
	}
	F5 := []string{"████████", "█      ", "███████ ", "█      ", "█      "}
	O5 := []string{" █████ ", "█     █", "█     █", "█     █", " █████ "}
	R5 := []string{"███████ ", "█     █", "███████ ", "█    █ ", "█     █"}
	G5 := []string{" █████ ", "█      ", "█  ███ ", "█    █ ", " █████ "}
	E5 := []string{"████████", "█      ", "███████ ", "█      ", "████████"}
	F3 := normalizeLetter(compress5to3(F5))
	O3 := normalizeLetter(compress5to3(O5))
	R3 := normalizeLetter(compress5to3(R5))
	G3 := normalizeLetter(compress5to3(G5))
	E3 := normalizeLetter(compress5to3(E5))
	spacer := " "
	rows := make([]string, 3)
	for i := 0; i < 3; i++ {
		rows[i] = F3[i] + spacer + O3[i] + spacer + R3[i] + spacer + G3[i] + spacer + E3[i]
	}
	einzug := strings.Repeat(" ", 10)
	for i, r := range rows {
		full := einzug + r
		// Kürzen falls zu breit
		runes := []rune(full)
		if len(runes) > w {
			full = string(runes[:w])
		}
		// Farbverlauf anwenden (ähnlich RenderTitleBar)
		rows[i] = gradientText(full, [3]int{181, 0, 19}, [3]int{204, 204, 204})
	}
	return rows
}

// compress5to3 nimmt 5 Textzeilen (volle Blöcke für gefüllt, Leerzeichen für leer) und erzeugt 3 Zeilen
// mit ▀ für obere Hälfte, ▄ für untere Hälfte, █ für voll und Leerzeichen für leer.
func compress5to3(rows []string) []string {
	width := 0
	for _, r := range rows {
		if len(r) > width {
			width = len(r)
		}
	}
	padded := make([][]rune, 5)
	for i := 0; i < 5; i++ {
		rr := []rune(rows[i])
		if len(rr) < width {
			// mit Leerzeichen auffüllen
			tmp := make([]rune, width)
			copy(tmp, rr)
			for j := len(rr); j < width; j++ {
				tmp[j] = ' '
			}
			rr = tmp
		}
		padded[i] = rr
	}

	out := make([]string, 3)
	pairs := [][2]int{{0, 1}, {2, 3}, {4, -1}}
	for idx, pair := range pairs {
		var b strings.Builder
		topIdx, botIdx := pair[0], pair[1]
		for c := 0; c < width; c++ {
			top := padded[topIdx][c] != ' '
			bottom := false
			if botIdx >= 0 {
				bottom = padded[botIdx][c] != ' '
			}
			var ch rune
			switch {
			case top && bottom:
				ch = '█'
			case top && !bottom:
				ch = '▀'
			case !top && bottom:
				ch = '▄'
			default:
				ch = ' '
			}
			b.WriteRune(ch)
		}
		out[idx] = b.String()
	}
	return out
}

// gradientText färbt jede Rune mit horizontalem Farbverlauf von startRGB nach endRGB.
func gradientText(s string, startRGB, endRGB [3]int) string {
	runes := []rune(s)
	n := len(runes)
	if n == 0 {
		return s
	}
	if n == 1 {
		hex := fmt.Sprintf("#%02X%02X%02X", startRGB[0], startRGB[1], startRGB[2])
		return titleTextStyle.Foreground(lipgloss.Color(hex)).Render(string(runes[0]))
	}
	var b strings.Builder
	b.Grow(len(s) * 10)
	for i, r := range runes {
		// linearer Interpolationsfaktor in [0,1]
		t := float64(i) / float64(n-1)
		cr := int(float64(startRGB[0])*(1-t) + float64(endRGB[0])*t)
		cg := int(float64(startRGB[1])*(1-t) + float64(endRGB[1])*t)
		cb := int(float64(startRGB[2])*(1-t) + float64(endRGB[2])*t)
		hex := fmt.Sprintf("#%02X%02X%02X", cr, cg, cb)
		b.WriteString(titleTextStyle.Foreground(lipgloss.Color(hex)).Render(string(r)))
	}
	return b.String()
}

// normalizeLetter trimmt jede Zeile rechts und füllt anschließend alle Zeilen auf die maximale Breite auf, um das Glyph auszurichten.
func normalizeLetter(rows []string) []string {
	trimmed := make([]string, len(rows))
	maxw := 0
	for i, r := range rows {
		tr := strings.TrimRight(r, " ")
		trimmed[i] = tr
		if lw := len([]rune(tr)); lw > maxw {
			maxw = lw
		}
	}
	out := make([]string, len(rows))
	for i, tr := range trimmed {
		lw := len([]rune(tr))
		if lw < maxw {
			tr = tr + strings.Repeat(" ", maxw-lw)
		}
		out[i] = tr
	}
	return out
}
