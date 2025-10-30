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
	// Basis-Muster.
	muster := "//////"
	// Wie oft muss das Muster wiederholt werden um mindestens die Breite zu erreichen?
	wiederholungen := h.breite/len(muster) + 1
	linieBasis := strings.Repeat(muster, wiederholungen)
	// Erzeuge Gradient Styles (Start -> Ende) mit gewünschter Stufenzahl.
	// h.stufen = 0 bedeutet automatische Heuristik; >0 setzt explizite Segmentanzahl.
	styles := gradientStyles("#e30018", "#dfaeb3", h.breite, h.stufen)

	// Hilfsfunktion zur Anwendung des Gradienten auf eine Zeile.
	// Jede Zeile wird in gleich große Segmente geschnitten; Rest hängt an das letzte Segment.
	// Verarbeitung erfolgt auf Rune-Ebene damit Mehrbyte-Zeichen (▀▄█ usw.) nicht zerstückelt werden.
	var renderGradient = func(roherText string, bold bool, suffix string) string {
		if len(roherText) == 0 {
			return ""
		}
		// In Runen umwandeln für korrekte Zeichenbreiten-Verarbeitung.
		runes := []rune(roherText)
		breite := len(runes)
		segmente := len(styles)
		if segmente == 0 {
			return roherText + suffix
		}
		basisLaenge := breite / segmente
		rest := breite % segmente
		var teile []string
		offset := 0
		for s := 0; s < segmente; s++ {
			laenge := basisLaenge
			if s == segmente-1 { // letztes Segment bekommt Rest
				laenge += rest
			}
			if laenge <= 0 {
				continue
			}
			segmentRunes := runes[offset : offset+laenge]
			segmentText := string(segmentRunes)
			st := styles[s]
			if bold {
				st = st.Bold(true)
			}
			teile = append(teile, st.Render(segmentText))
			offset += laenge
		}
		if suffix != "" {
			// Suffix separat stylen (Bold wenn gefordert, gleiche letzte Farbe für optisch weichen Übergang)
			st := styles[len(styles)-1]
			if bold {
				st = st.Bold(true)
			}
			teile = append(teile, st.Render(suffix))
		}
		return strings.Join(teile, "")
	}

	// Erste Zeile mit Wort "SEW" so ausrichten, dass das 'S' unter dem 'F' des FORGE-Schriftzuges steht (Einzug 10).
	einzug := strings.Repeat(" ", 10)
	sew := "SEW"
	// Verfügbare Breite für Auffüllmuster links und rechts bestimmen.
	// Gesamte Zeile: Einzug + SEW + rechts Rest mit Muster füllen.
	var erste string
	if h.breite <= len(einzug)+len(sew) { // Edge Case: sehr schmale Breite
		basis := einzug + sew
		if len(basis) > h.breite {
			basisRunes := []rune(basis)
			basis = string(basisRunes[:h.breite])
		}
		erste = renderGradient(basis, true, "")
	} else {
		restBreite := h.breite - (len(einzug) + len(sew))
		auffuell := linieBasis
		if len(auffuell) > restBreite {
			auffuell = auffuell[:restBreite]
		}
		roh := einzug + sew + auffuell
		erste = renderGradient(roh, true, "")
	}

	// Untere drei Zeilen durch FORGE-Schriftzug ersetzen.
	forge := h.forgeLines(styles)
	// Absicherung falls kürzer (z.B. extrem kleine Breite).
	for len(forge) < 3 {
		forge = append(forge, "")
	}
	return strings.Join([]string{erste, forge[0], forge[1], forge[2]}, "\n")
}

// forgeLines erzeugt drei Zeilen mit dem Schriftzug "FORGE" in kompakter Blockschrift.
// Ausgangspunkt ist eine 5-Zeilen-Matrix je Buchstabe, die auf 3 Zeilen reduziert wird.
// Verwendete Zeichen: ▀ (nur oberer Halbblock), ▄ (nur unterer Halbblock), █ (voller Block), Leerzeichen.
// Linkseinzug: 10 Leerzeichen. Am Ende werden die Zeilen auf die verfügbare Breite gekürzt.
// Farbverlauf: Es wird der übergebene Style-Gradient segmentweise angewendet.
func (h headerModel) forgeLines(styles []lipgloss.Style) []string {
	// 5-Zeilen Rohdefinition mit 'X' = gesetzt, ' ' = leer; alle Buchstaben gleiche Breite für einfache Verarbeitung.
	// Breite je Buchstabe: 5 Spalten.
	letters := map[rune][]string{
		'F': {"XXXXX", "X    ", "XXXX ", "X    ", "X    "},
		'O': {"XXXXX", "X   X", "X   X", "X   X", "XXXXX"},
		'R': {"XXXX ", "X   X", "XXXX ", "X   X", "X   X"},
		'G': {"XXXXX", "X    ", "X XXX", "X   X", "XXXXX"},
		'E': {"XXXXX", "X    ", "XXXX ", "X    ", "XXXXX"},
	}
	wort := "FORGE"
	// Normalisierung: fehlende Buchstaben (sollte nicht passieren) -> leere Matrix.
	var rohw [][][]rune // Liste von Buchstaben-Matrizen
	for _, r := range wort {
		pattern, ok := letters[r]
		if !ok {
			pattern = []string{"     ", "     ", "     ", "     ", "     "}
		}
		matrix := make([][]rune, len(pattern))
		for i, line := range pattern {
			matrix[i] = []rune(line)
		}
		rohw = append(rohw, matrix)
	}
	// Komprimierung von 5 auf 3 Zeilen: Paare (0,1), (2,3) und Rest (4).
	// Für eine schlankere Unterkante verwenden wir in der dritten Zeile einen unteren Halbblock ('▄') statt eines vollen Blocks.
	compressPair := func(top, bottom rune, last bool) rune {
		oben := top == 'X'
		unten := bottom == 'X'
		if last { // letzte Einzelzeile nur als unterer Halbblock zeichnen
			if oben {
				return '▀'
			}
			return ' '
		}
		if oben && unten {
			return '█'
		}
		if oben && !unten {
			return '▀'
		}
		if !oben && unten {
			return '▄'
		}
		return ' '
	}
	// Erzeuge für jeden Buchstaben seine 3 komprimierten Zeilen.
	var buchstabenZeilen [3][]rune
	for li, matrix := range rohw {
		breiteB := len(matrix[0])
		// Initialisierung bei erstem Buchstaben.
		if li == 0 {
			for i := 0; i < 3; i++ {
				buchstabenZeilen[i] = []rune{}
			}
		}
		for col := 0; col < breiteB; col++ {
			// Zeile 1 aus Matrix[0] & Matrix[1]
			buchstabenZeilen[0] = append(buchstabenZeilen[0], compressPair(matrix[0][col], matrix[1][col], false))
			// Zeile 2 aus Matrix[2] & Matrix[3]
			buchstabenZeilen[1] = append(buchstabenZeilen[1], compressPair(matrix[2][col], matrix[3][col], false))
			// Zeile 3 aus Matrix[4]
			buchstabenZeilen[2] = append(buchstabenZeilen[2], compressPair(matrix[4][col], ' ', true))
		}
		// Abstand zwischen Buchstaben (2 Leerzeichen als optische Trennung) außer letztem.
		if li < len(rohw)-1 {
			for i := 0; i < 3; i++ {
				buchstabenZeilen[i] = append(buchstabenZeilen[i], ' ', ' ')
			}
		}
	}
	// Aufbau finaler Textzeilen mit Einzug.
	einzugForge := strings.Repeat(" ", 10)
	var result []string
	applyGradient := func(line string) string {
		if len(line) == 0 {
			return line
		}
		runes := []rune(line)
		breite := len(runes)
		segmente := len(styles)
		if segmente == 0 {
			return line
		}
		basis := breite / segmente
		rest := breite % segmente
		teile := make([]string, 0, segmente)
		offset := 0
		for s := 0; s < segmente; s++ {
			laenge := basis
			if s == segmente-1 {
				laenge += rest
			}
			if laenge <= 0 {
				continue
			}
			segment := string(runes[offset : offset+laenge])
			teile = append(teile, styles[s].Bold(true).Render(segment))
			offset += laenge
		}
		return strings.Join(teile, "")
	}
	for i := 0; i < 3; i++ {
		rohRunes := append([]rune(einzugForge), buchstabenZeilen[i]...)
		// Kürzen auf verfügbare Breite (Rune-basiert für UTF-8 Sicherheit).
		if len(rohRunes) > h.breite {
			rohRunes = rohRunes[:h.breite]
		}
		roh := string(rohRunes)
		result = append(result, applyGradient(roh))
	}
	return result
}
