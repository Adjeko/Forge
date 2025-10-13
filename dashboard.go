package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// stripANSI entfernt ANSI Escape Sequenzen für Breitenberechnung.
func stripANSI(s string) string {
	// Einfache Zustandsmaschine: ESC [ ... Buchstaben entfernen
	b := strings.Builder{}
	inEsc := false
	seq := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !inEsc {
			if c == 0x1b { // ESC
				inEsc = true
				seq = false
				continue
			}
			b.WriteByte(c)
		} else {
			// Nach ESC erwarten wir '[' oder andere; überspringen bis Buchstabe
			if !seq {
				if c == '[' || c == ']' || c == '(' || c == ')' || c == '>' || c == '<' || c == 'O' || c == 'P' {
					seq = true
					continue
				}
			}
			// Ende der CSI Sequenz wenn Buchstabe zwischen @ und ~
			if seq && (c >= '@' && c <= '~') {
				inEsc = false
				seq = false
			}
		}
	}
	return b.String()
}

// RenderDashboard baut die Zeilen zwischen Header und Footer.
// Enthalten sind:
// - eine leere Trennzeile nach dem Header
// - Terminalausgabe-Bereich (Viewport), unten ausgerichtet
// - eine Trennzeile mit rechtsbündigem Prozent-Indikator
// - eine leere Zeile (vor Footer)
func RenderDashboard(m Model, innerW, innerH, headerHeight, footerHeight int) []string {
	lines := make([]string, 0, innerH)

	// Verfügbare Zeilen zwischen Header und Footer (Footer kann mehrzeilig sein)
	avail := innerH - headerHeight - footerHeight
	if avail < 0 {
		avail = 0
	}

	// Benötigt: 1 (Header-Separ.) + outputRows + 1 (Separ.) + 1 (Leer)
	// => outputRows = avail - 3
	outputRows := avail - 3
	if outputRows < 0 {
		outputRows = 0
	}

	// 1) Leere Zeile nach Header
	lines = append(lines, strings.Repeat(" ", innerW))

	// 2) Ausgabebereich
	if m.showHelp {
		// Hilfe im Ausgabebereich anzeigen
		helpLines := []string{
			"",
			"Taste     Aktion",
			"^+Y        Hilfe Umschalten",
			"Esc       Hilfe Schließen",
			"^+C        Befehl Ausführen",
			"^+Q        Programm Beenden",
			"",
			"Diese Leiste unten zeigt stets die aktuell gültigen Hotkeys.",
		}
		helpBlock := contentStyle.Width(innerW).Render(strings.Join(helpLines, "\n"))
		padded := padToLines(helpBlock, outputRows)
		if outputRows > 0 {
			lines = append(lines, strings.Split(padded, "\n")...)
		}
	} else {
		if outputRows > 0 {
			// Viewport rendern; kann kürzer sein als outputRows, dann oben mit Leerzeilen auffüllen
			vpView := m.vp.View()
			vpLines := strings.Split(vpView, "\n")
			if len(vpLines) > outputRows {
				vpLines = vpLines[len(vpLines)-outputRows:]
			}
			// oben auffüllen
			for i := 0; i < outputRows-len(vpLines); i++ {
				lines = append(lines, strings.Repeat(" ", innerW))
			}
			// jede Zeile auf innerW beschneiden
			for _, ln := range vpLines {
				r := []rune(ln)
				if len(r) > innerW {
					ln = string(r[:innerW])
				}
				if len([]rune(ln)) < innerW {
					ln = ln + strings.Repeat(" ", innerW-len([]rune(ln)))
				}
				lines = append(lines, ln)
			}
		}
	}

	// Modal-Overlay via Lipgloss-Komposition
	if (m.showModal || m.showTree || m.showCreate || m.showFlows || m.showFlowsCreate || m.showStatus || m.showStatusCreate || m.showGitStatus) && outputRows > 0 && innerW > 20 {
		var formView string
		if m.showCreate && m.createForm != nil {
			contentW := innerW - 10
			if contentW < 40 {
				contentW = innerW - 6
			}
			if contentW < 40 {
				contentW = innerW - 4
			}
			if contentW < 40 {
				contentW = innerW
			}
			m.createForm.SetWidth(contentW)
			formView = m.createForm.View()
		} else if (m.showTree || m.showFlows || m.showStatus) && m.tree != nil && !m.showFlowsCreate && !m.showStatusCreate {
			// Breite teilen: 2/3 Baum, 1/3 Detail (minus Padding)
			contentW := innerW - 10
			if contentW < 30 {
				contentW = innerW - 6
			}
			if contentW < 30 {
				contentW = innerW - 4
			}
			if contentW < 30 {
				contentW = innerW - 2
			}
			if contentW < 30 {
				contentW = innerW - 0
			}
			if contentW < 30 {
				contentW = innerW
			}
			if contentW < 30 {
				contentW = 30
			}
			m.tree.SetWidth(contentW * 2 / 3)
			m.tree.SetDetailWidth(contentW - m.tree.Width())
			formView = m.tree.View()
		} else if m.showFlowsCreate && m.flowsCreateForm != nil {
			contentW := innerW - 10
			if contentW < 40 {
				contentW = innerW - 6
			}
			if contentW < 40 {
				contentW = innerW - 4
			}
			if contentW < 40 {
				contentW = innerW
			}
			m.flowsCreateForm.SetWidth(contentW)
			formView = m.flowsCreateForm.View()
		} else if m.showStatusCreate && m.statusCreateForm != nil {
			contentW := innerW - 10
			if contentW < 40 {
				contentW = innerW - 6
			}
			if contentW < 40 {
				contentW = innerW - 4
			}
			if contentW < 40 {
				contentW = innerW
			}
			m.statusCreateForm.SetWidth(contentW)
			formView = m.statusCreateForm.View()
		} else if m.gitStatusForm != nil && m.showGitStatus {
			contentW := innerW - 20
			if contentW < 40 {
				contentW = innerW - 10
			}
			if contentW < 30 {
				contentW = innerW
			}
			m.gitStatusForm.SetWidth(contentW)
			formView = m.gitStatusForm.View()
		} else if m.form != nil {
			// Ziel-Innenbreite (Rand + Padding ~4 Zeichen ausnehmen)
			target := innerW - 10
			if target < 20 {
				target = innerW - 6
			}
			if target < 20 {
				target = 20
			}
			m.form.SetWidth(target)
			formView = m.form.View()
		} else {
			formView = "Befehl ausführen\n\nDrücke Esc oder ^+C zum Schließen."
		}
		// Stil für Modal
		modalStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF4A5A")). // helleres SEW-Rot
			Padding(0, 1)
		rendered := modalStyle.Render(formView)
		modalLines := strings.Split(rendered, "\n")
		// Platzierung bestimmen
		modalHeight := len(modalLines)
		modalWidth := 0
		for _, l := range modalLines {
			if lw := len([]rune(stripANSI(l))); lw > modalWidth {
				modalWidth = lw
			}
		}
		if modalWidth > innerW-2 {
			modalWidth = innerW - 2
		}
		// Vertikal im Ausgabebereich zentrieren
		startY := 1 // erste Ausgabeline ist Index 1
		if outputRows > modalHeight {
			startY += (outputRows - modalHeight) / 2
		}
		// Horizontal zentrieren
		startX := 0
		if innerW > modalWidth {
			startX = (innerW - modalWidth) / 2
		}
		// Zeilen überlagern
		for i, ml := range modalLines {
			row := startY + i
			if row >= len(lines) || row < 1 || row > outputRows { // innerhalb Ausgabebereich bleiben
				continue
			}
			base := []rune(lines[row])
			// Basislänge == innerW sicherstellen
			if len(base) < innerW {
				base = append(base, []rune(strings.Repeat(" ", innerW-len(base)))...)
			}
			// ml an startX einsetzen (ggf. abschneiden – ANSI nicht mittendrin brechen)
			if startX >= 0 && startX < innerW {
				left := string(base[:startX])
				lineOut := left + ml
				if len([]rune(stripANSI(lineOut))) < innerW {
					padN := innerW - len([]rune(stripANSI(lineOut)))
					lineOut += strings.Repeat(" ", padN)
				}
				lines[row] = lineOut
			}
		}
	}

	// 3) Trennzeile zwischen Ausgabebereich und Footer-Abstand; Prozentanzeige rechts
	sep := strings.Repeat(" ", innerW)
	if outputRows > 0 && innerW > 0 {
		// Prozent basierend auf umgebrochenem Inhalt und Viewport-Offset
		wrapped := flattenOutput(m.termOutput, innerW)
		total := len(wrapped)
		bottom := 0
		if total > 0 {
			y := m.vp.YOffset
			if y < 0 {
				y = 0
			}
			shown := m.vp.Height
			if shown < 0 {
				shown = 0
			}
			if shown > outputRows {
				shown = outputRows
			}
			b := y + shown
			if b > total {
				b = total
			}
			bottom = b
		}
		var percent int
		if total <= 0 {
			percent = 0
		} else {
			percent = int(float64(bottom)*100.0/float64(total) + 0.5)
		}
		indicator := fmt.Sprintf("%3d%%", percent)
		startX := innerW - len([]rune(indicator))
		if startX < 0 {
			startX = 0
		}
		sep = insertAt(sep, indicator, startX, innerW)
	}
	lines = append(lines, sep)

	// 4) Leere Zeile zwischen Ausgabe und Footer
	lines = append(lines, strings.Repeat(" ", innerW))

	return lines
}

// RenderPullRequests erzeugt Platzhalter-Zeilen für die Pull Requests Seite.
// Spiegelt Struktur von RenderDashboard, zeigt aber zentrierte Meldung.
func RenderPullRequests(m Model, innerW, innerH, headerHeight, footerHeight int) []string {
	lines := make([]string, 0, innerH)
	avail := innerH - headerHeight - footerHeight
	if avail < 0 {
		avail = 0
	}
	outputRows := avail - 3 // header sep + sep + blank
	if outputRows < 0 {
		outputRows = 0
	}
	// 1) Leer nach Header
	lines = append(lines, strings.Repeat(" ", innerW))
	// 2) Platzhalter-Bereich mit vertikal zentrierter Nachricht
	// Content aus m.pullRequests oder Fallback-Nachricht
	var prLines []string
	if len(m.pullRequests) == 0 {
		prLines = []string{"Noch keine Pull Requests geladen (drücke ^+P)"}
	} else {
		for _, pr := range m.pullRequests {
			// Format: #num [STATUS] title (author) branch
			status := pr.Status
			line := fmt.Sprintf("%s [%s] %s (%s) %s", pr.Number, status, pr.Title, pr.Author, pr.Branch)
			prLines = append(prLines, line)
		}
	}
	// Block zusammenfügen
	msg := strings.Join(prLines, "\n")
	if outputRows > 0 {
		// Breite pro Zeile sicherstellen; Block vertikal zentrieren
		blockLines := strings.Split(msg, "\n")
		// Jede Zeile auf innerW beschneiden
		for i, bl := range blockLines {
			r := []rune(bl)
			if len(r) > innerW {
				blockLines[i] = string(r[:innerW])
			} else if len(r) < innerW {
				blockLines[i] = bl + strings.Repeat(" ", innerW-len(r))
			}
		}
		blockHeight := len(blockLines)
		if blockHeight > outputRows {
			blockLines = blockLines[:outputRows]
			blockHeight = outputRows
		}
		padTop := 0
		if outputRows > blockHeight {
			padTop = (outputRows - blockHeight) / 2
		}
		for i := 0; i < padTop; i++ {
			lines = append(lines, strings.Repeat(" ", innerW))
		}
		lines = append(lines, blockLines...)
		for len(lines) < 1+outputRows { // restliche Zeilen auffüllen
			lines = append(lines, strings.Repeat(" ", innerW))
		}
	}
	// 3) Separatore mit Prozent (immer 0%)
	sep := strings.Repeat(" ", innerW)
	indicator := "  0%"
	startX := innerW - len([]rune(indicator))
	if startX < 0 {
		startX = 0
	}
	sep = insertAt(sep, indicator, startX, innerW)
	lines = append(lines, sep)
	// 4) Leere Zeile
	lines = append(lines, strings.Repeat(" ", innerW))
	return lines
}

// Helfer sind in model.go definiert und werden hier wiederverwendet: insertAt, minInt, flattenOutput
