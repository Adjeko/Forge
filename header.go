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
	var renderGradient = func(roherText string, bold bool, suffix string) string {
		if len(roherText) == 0 {
			return ""
		}
		breite := len(roherText)
		segmente := len(styles)
		basisLaenge := breite / segmente
		rest := breite % segmente
		var teile []string
		offset := 0
		for s := 0; s < segmente; s++ {
			laenge := basisLaenge
			if s == segmente-1 { // letztes Segment bekommt Rest
				laenge += rest
			}
			if laenge == 0 {
				continue
			}
			segmentText := roherText[offset : offset+laenge]
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

	// Erste Zeile mit Zusatz SEW (Abtrennung durch Leerzeichen für Lesbarkeit).
	ersteRoh := linieBasis[:h.breite-len(" SEW")] // Platz für Suffix freilassen
	erste := renderGradient(ersteRoh, true, " SEW")
	zweite := renderGradient(linieBasis[:h.breite], false, "")
	dritte := renderGradient(linieBasis[:h.breite], false, "")
	vierte := renderGradient(linieBasis[:h.breite], false, "")
	return strings.Join([]string{erste, zweite, dritte, vierte}, "\n")
}
