package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// headerModel repräsentiert den separaten Kopfbereich (Header) der Anwendung.
// Er ist genau 4 Zeilen hoch und besteht in jeder Zeile aus der Wiederholung des Musters "//////".
// In der ersten Zeile wird zusätzlich " SEW" hinter das Muster gesetzt.
type headerModel struct {
	breite int  // aktuelle verfügbare Breite im Terminal
	fertig bool // Flag ob eine Breite gesetzt wurde
}

// newHeaderModel erzeugt ein leeres Header-Modell.
func newHeaderModel() headerModel { return headerModel{} }

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
	// Erste Zeile mit Zusatz SEW (Abtrennung durch Leerzeichen für Lesbarkeit).
	erste := lipgloss.NewStyle().Bold(true).Render(linieBasis[:h.breite-len(" SEW")] + " SEW")
	// Weitere Zeilen rein aus dem Muster, auf Breite geschnitten.
	zweite := linieBasis[:h.breite]
	dritte := linieBasis[:h.breite]
	vierte := linieBasis[:h.breite]
	return strings.Join([]string{erste, zweite, dritte, vierte}, "\n")
}
