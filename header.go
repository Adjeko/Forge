package main

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// gradientStyles berechnet eine Liste von Lipgloss Styles, die einen linearen Farbverlauf
// zwischen zwei Hex-Farben (z.B. "#de0000" und "#ffbcbc") über die angegebene Breite darstellen.
// Die Anzahl der Segmente wird adaptiv gesteuert, um nicht zu viele Styles zu erzeugen.
func gradientStyles(startHex, endHex string, breite int, segmenteWunsch int) []lipgloss.Style {
	// Sicherstellen, dass Hex mit '#'
	if len(startHex) != 7 || startHex[0] != '#' || len(endHex) != 7 || endHex[0] != '#' {
		// Fallback auf einfaches Rot wenn Format falsch (Fehler hier nicht fatal)
		return []lipgloss.Style{lipgloss.NewStyle().Foreground(lipgloss.Color("#e100ff"))}
	}
	parse := func(h string) (int, int, int) {
		r, _ := strconv.ParseInt(h[1:3], 16, 64)
		g, _ := strconv.ParseInt(h[3:5], 16, 64)
		b, _ := strconv.ParseInt(h[5:7], 16, 64)
		return int(r), int(g), int(b)
	}
	startR, startG, startB := parse(startHex)
	endR, endG, endB := parse(endHex)
	// Segmentanzahl: wenn segmenteWunsch > 0 genutzt, sonst heuristisch (~ alle 12 Zeichen).
	minSeg := 2  // Minimal 2 für sichtbaren Verlauf
	maxSeg := 48 // Obergrenze um Performance & Lesbarkeit zu sichern
	var segmente int
	if segmenteWunsch > 0 {
		segmente = segmenteWunsch
	} else {
		segmente = breite / 12
	}
	if segmente < minSeg {
		segmente = minSeg
	}
	if segmente > maxSeg {
		segmente = maxSeg
	}
	styles := make([]lipgloss.Style, segmente)
	for i := 0; i < segmente; i++ {
		t := float64(i) / float64(segmente-1)
		r := int(float64(startR) + (float64(endR)-float64(startR))*t)
		g := int(float64(startG) + (float64(endG)-float64(startG))*t)
		b := int(float64(startB) + (float64(endB)-float64(startB))*t)
		hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)
		styles[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(hex))
	}
	return styles
}

// headerModel repräsentiert den separaten Kopfbereich (Header) der Anwendung.
// Er ist genau 4 Zeilen hoch und besteht in jeder Zeile aus der Wiederholung des Musters "//////".
// In der ersten Zeile wird zusätzlich " SEW" hinter das Muster gesetzt.
type headerModel struct {
	breite int  // aktuelle verfügbare Breite im Terminal
	fertig bool // Flag ob eine Breite gesetzt wurde
	stufen int  // Anzahl gewünschter Farb-Stufen im Verlauf (0 = automatisch)
}

// newHeaderModel erzeugt ein leeres Header-Modell.
func newHeaderModel() headerModel { return headerModel{stufen: 30} }

// Init erfüllt das Bubble Tea Interface, kein Start-Command nötig.
func (h headerModel) Init() tea.Cmd { return nil }

// Update verarbeitet relevante Nachrichten; hier nur Terminalgrößen.
func (h headerModel) Update(msg tea.Msg) (headerModel, tea.Cmd) {
	switch nachricht := msg.(type) {
	case tea.WindowSizeMsg:
		// Speichert Breite und markiert als fertig.
		h.breite = nachricht.Width
		h.fertig = true
	}
	return h, nil
}

// View rendert den 4-zeiligen Header.
func (h headerModel) View() string {
	if !h.fertig {
		return "" // Keine Ausgabe bis Breite bekannt
	}
	// Stil für das Muster "///" in gewünschter Farbe (#E30018)
	musterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E30018")).Bold(true)
	unit := "///"
	// Hilfsfunktion zum Erzeugen eines gefärbten Musters mit exakt gewünschter sichtbarer Länge.
	patternFill := func(l int) string {
		if l <= 0 {
			return ""
		}
		repeat := l/len(unit) + 2 // +2 Sicherheitszugabe
		roh := strings.Repeat(unit, repeat)
		roh = roh[:l]
		return musterStyle.Render(roh)
	}

	// Sichtbare Länge einer Zeichenkette (ANSI-fähig)
	visible := func(s string) int { return lipgloss.Width(s) }

	// Konfigurierbare linke Abstände für SEW und FORGE separat, um Leerraum zu minimieren.
	linkerBlockSew := 10   // Abstand vor SEW (wieder auf 10 gesetzt)
	linkerBlockForge := 10 // Abstand vor FORGE-Zeilen (wieder auf 10 gesetzt)

	// Breite des FORGE-Artworks bestimmen um das rechte Muster exakt ab dieser Spalte zu starten.
	// Dafür eine Probe der Roh-FORGE-Zeilen ohne Farben erstellen wie in RenderForgeLines.
	F5 := []string{"████████", "█      ", "███████ ", "█      ", "█      "}
	O5 := []string{" █████ ", "█     █", "█     █", "█     █", " █████ "}
	R5 := []string{"███████ ", "█     █", "███████ ", "█    █ ", "█     █"}
	G5 := []string{" █████ ", "█      ", "█  ███ ", "█    █ ", " █████ "}
	E5 := []string{"████████", "█      ", "███████ ", "█      ", "████████"}
	compress5to3Local := func(rows []string) []string {
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
	normalizeLetterLocal := func(rows []string) []string {
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
	F3 := normalizeLetterLocal(compress5to3Local(F5))
	O3 := normalizeLetterLocal(compress5to3Local(O5))
	R3 := normalizeLetterLocal(compress5to3Local(R5))
	G3 := normalizeLetterLocal(compress5to3Local(G5))
	E3 := normalizeLetterLocal(compress5to3Local(E5))
	spacer := " "
	artWidth := 0
	for i := 0; i < 3; i++ {
		r := F3[i] + spacer + O3[i] + spacer + R3[i] + spacer + G3[i] + spacer + E3[i]
		if l := len([]rune(r)); l > artWidth {
			artWidth = l
		}
	}

	// Layout erste Zeile:
	// Linkes Muster (wie gehabt linkerBlockSew), dann ein Leerzeichen, dann SEW, dann ein Raum bis zur Spalte artWidth,
	// in diesem Raum rechtsbündig Versionsstring, danach startet das rechte Muster.
	sewText := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#E30018")).Render("SEW")
	versionPlain := "v0.0.1"
	versionStyled := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#E30018")).Render(versionPlain)
	sewLen := visible(sewText)
	versionLen := len([]rune(versionPlain))
	// Mindestbreite prüfen
	minNeeded := linkerBlockSew + 1 + sewLen + 1 + versionLen // grob für Raum + rechte Muster-Start
	if h.breite < minNeeded {
		return sewText
	} // zu schmal, fallback
	leftPattern := patternFill(linkerBlockSew)
	// Bereich zwischen Ende SEW und artWidth
	spaceBetween := artWidth - sewLen
	if spaceBetween < versionLen+1 { // +1 für Mindestabstand
		spaceBetween = versionLen + 1
	}
	// Version rechtsbündig in spaceBetween
	padLen := spaceBetween - versionLen
	zwischenBlock := strings.Repeat(" ", padLen) + versionStyled
	// Rechtes Muster beginnt ab Gesamtlänge: linkerPattern + Leerzeichen + SEW + zwischenBlock
	// Ein zusätzliches Leerzeichen nach der Version vor dem rechten Muster einkalkulieren.
	used := visible(leftPattern) + 1 + sewLen + visible(zwischenBlock) + 1
	rightAvail := h.breite - used
	if rightAvail < 0 {
		rightAvail = 0
	}
	rightPattern := patternFill(rightAvail)
	zeile1 := leftPattern + " " + sewText + zwischenBlock + " " + rightPattern

	// FORGE-Zeilen generieren (bestehende Farbverläufe beibehalten) und links/rechts mit Muster füllen.
	forgeRaw := RenderForgeLines(h.breite)
	for len(forgeRaw) < 3 {
		forgeRaw = append(forgeRaw, "")
	}
	forgeOut := make([]string, 3)
	for i := 0; i < 3; i++ {
		line := forgeRaw[i]
		// Führende Einrückung entfernen bis zur gewünschten linken Musterbreite.
		cut := 0
		for _, r := range line {
			if r == ' ' && cut < linkerBlockForge {
				cut++
				continue
			}
			break
		}
		lineContent := line[cut:]
		contentWidth := visible(lineContent)
		// Benötigte Breite: linkes Muster + 2 Leerzeichen (eins nach linkem Muster, eins vor rechtem Muster) + Inhalt + rechtes Muster.
		leftPattern := patternFill(linkerBlockForge)
		// Restbreite für rechtes Muster berechnen.
		rest := h.breite - linkerBlockForge - contentWidth - 2 // -2 für die beiden Leerzeichen
		if rest < 0 {
			rest = 0
		}
		rightPattern := patternFill(rest)
		forgeOut[i] = leftPattern + " " + lineContent + " " + rightPattern
	}

	// Zusätzliche fünfte Zeile: komplettes Muster über gesamte Breite
	extraLine := patternFill(h.breite)

	return strings.Join([]string{zeile1, forgeOut[0], forgeOut[1], forgeOut[2], extraLine}, "\n")
}

// Lokale forgeLines Implementierung entfernt; stattdessen RenderForgeLines (forgeTitle.go) genutzt.
