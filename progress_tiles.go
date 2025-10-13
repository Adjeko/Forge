package main

import (
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProgressTile repräsentiert eine zweizeilige Kachel mit ID+Name und Fortschrittsbalken.
type ProgressTile struct {
	ID      string
	Name    string
	Bar     progress.Model
	Percent float64
}

func NewProgressTile(id, name string) ProgressTile {
	// Gradient entspricht dem FORGE Wort (Start #B50013 -> Ende #CCCCCC)
	bar := progress.New(
		progress.WithGradient("#B50013", "#CCCCCC"),
		progress.WithoutPercentage(),
	)
	bar.Width = 10 // Breite wird bei der Darstellung dynamisch angepasst
	return ProgressTile{ID: id, Name: name, Bar: bar, Percent: 0.0}
}

// Tick animiert die Kachel; gibt ein Cmd für weiteres Ticken zurück.
func (t *ProgressTile) Tick() tea.Cmd {
	// einfacher Auto-Inkrement Loop
	if t.Percent < 1.0 {
		step := 0.02 + 0.06*math.Sin(float64(time.Now().UnixNano())/5e8)
		p := t.Percent + step
		if p > 1.0 {
			p = 1.0
		}
		t.Percent = p
	}
	return tea.Tick(120*time.Millisecond, func(time.Time) tea.Msg { return progressTickMsg{} })
}

type progressTickMsg struct{}

// RenderProgressTiles rendert zwei Zeilen (ID+Name Reihe, Balken Reihe) über volle Breite.
// Kacheln horizontal angeordnet mit 1 Zeichen Padding um einen vertikalen Trenner wie " │ ".
func RenderProgressTiles(m Model, totalWidth int) []string {
	if totalWidth <= 0 {
		return []string{"", ""}
	}
	if len(m.tiles) == 0 {
		return []string{strings.Repeat(" ", totalWidth), strings.Repeat(" ", totalWidth)}
	}
	// pro-Kachel Breite bestimmen; Rest verteilen um vermeidbares End-Padding zu verhindern
	divider := "│"
	pad := 1
	divChunk := len(divider) + 2*pad
	segments := len(m.tiles)
	avail := totalWidth - (segments-1)*divChunk
	// Extrem schmal
	if avail < segments {
		avail = segments
	}
	base := avail / segments
	rem := avail % segments
	widths := make([]int, segments)
	for i := 0; i < segments; i++ {
		w := base
		if i < rem { // Rest auffüllen
			w++
		}
		if w < 4 {
			w = 4
		}
		widths[i] = w
	}
	// Letzte Breite anpassen sodass Summe(widths)+(segments-1)*divChunk == totalWidth exakt
	currentSum := 0
	for _, w := range widths {
		currentSum += w
	}
	needed := totalWidth - (segments-1)*divChunk
	if needed > 0 && currentSum != needed { // korrigieren
		delta := needed - currentSum
		widths[len(widths)-1] += delta
		if widths[len(widths)-1] < 4 { // Minimum sichern
			widths[len(widths)-1] = 4
		}
	}
	// Balkenbreiten je Kachel anpassen (mindestens 2 Zeichen für Label Bereich; Balken füllt visuell gesamte Kachelbreite)
	for i := range m.tiles {
		bw := widths[i]
		if bw < 4 { // Minimum sichern
			bw = 4
		}
		m.tiles[i].Bar.Width = bw
	}

	ensureVisibleWidth := func(s string, w int) string {
		if w <= 0 {
			return ""
		}
		vis := lipgloss.Width(s)
		if vis == w {
			return s
		}
		if vis < w { // auffüllen
			return s + strings.Repeat(" ", w-vis)
		}
		// vis > w -> Runen abschneiden bis passend
		runes := []rune(s)
		// naive: remove from end; ANSI sequences may be cut -> safer: rebuild accumulating visible width while skipping ANSI
		var b strings.Builder
		current := 0
		inEsc := false
		for _, r := range runes {
			if r == 0x1b { // ESC Start Escape-Sequenz
				inEsc = true
				b.WriteRune(r)
				continue
			}
			if inEsc {
				b.WriteRune(r)
				// Ende CSI Sequenz? Grobe Prüfung
				if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
					inEsc = false
				}
				continue
			}
			rw := lipgloss.Width(string(r))
			if current+rw > w { // passt nicht mehr
				break
			}
			b.WriteRune(r)
			current += rw
		}
		if current < w { // auffüllen
			b.WriteString(strings.Repeat(" ", w-current))
		}
		return b.String()
	}
	// erste Zeile (ID + Name gepadded), zweite Zeile (Fortschrittsbalken) aus den Breiten bauen
	var line1Parts []string
	var line2Parts []string
	for i := range m.tiles {
		t := m.tiles[i]
		w := widths[i]
		label := ensureVisibleWidth(t.ID+" "+t.Name, w)
		line1Parts = append(line1Parts, label)
		barView := ensureVisibleWidth(t.Bar.ViewAs(t.Percent), w) // Prozent in Ansicht
		line2Parts = append(line2Parts, barView)
	}
	var joinWithDiv func(parts []string) string
	joinWithDiv = func(parts []string) string {
		if len(parts) == 0 { // leer -> leere Zeile in Breite
			return strings.Repeat(" ", totalWidth)
		}
		var b strings.Builder
		for i, p := range parts {
			if i > 0 {
				b.WriteString(strings.Repeat(" ", pad))
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Render(divider))
				b.WriteString(strings.Repeat(" ", pad))
			}
			b.WriteString(p)
		}
		out := b.String()
		visLen := lipgloss.Width(out)
		if visLen < totalWidth { // Extra Leerzeichen verteilen
			// Extra Leerzeichen über Kacheln verteilen (round-robin) statt Block am Ende
			extra := totalWidth - visLen
			idx := len(parts) - 1
			for extra > 0 && idx >= 0 {
				parts[idx] += " "
				extra--
				idx--
				if idx < 0 {
					idx = len(parts) - 1
				}
			}
			// einmal neu bauen
			return joinWithDiv(parts)
		}
		if visLen > totalWidth { // zu lang -> kürzen
			// von letzter Kachel kürzen
			diff := visLen - totalWidth
			last := parts[len(parts)-1]
			runes := []rune(last)
			for diff > 0 && len(runes) > 0 { // solange diff vorhanden
				runes = runes[:len(runes)-1]
				diff = lipgloss.Width(out) - totalWidth
			}
			parts[len(parts)-1] = string(runes)
			return joinWithDiv(parts)
		}
		return out
	}
	l1 := joinWithDiv(line1Parts)
	l2 := joinWithDiv(line2Parts)
	return []string{l1, l2} // zwei Zeilen zurückgeben
}
