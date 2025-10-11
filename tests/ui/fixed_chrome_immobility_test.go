package ui_test

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"forge/src/ui/components"
	"forge/src/ui/model"
)

// T085 Fixed chrome immobility: header and footer must remain constant when scrolling output.
func TestFixedChromeImmobility(t *testing.T) {
	m := model.NewRootModel()
	headerInitial := components.Header()
	// render initial view
	view1 := m.View()
	if view1[:len(headerInitial)] != headerInitial {
		t.Fatalf("expected header at top of view")
	}
	// append lines to buffer to enable scrolling
	for i := 0; i < 200; i++ {
		m.Buffer().Append("line scrolling test")
	}
	// focus output to allow scroll (simulate tab cycles until output)
	for c := 0; c < 2; c++ { // CommandList -> Output
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}, Alt: false}) // placeholder non-tab
		// real tab navigation requires tea.KeyMsg with Type=tea.KeyTab; use that
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
		m = updated.(model.RootModel)
	}
	// perform scrolls
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	m = updated.(model.RootModel)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	m = updated.(model.RootModel)
	view2 := m.View()
	if view2[:len(headerInitial)] != headerInitial {
		t.Fatalf("header changed after scroll")
	}
	// footer presence check at end of view (last line starts with footer label)
	footer := components.Footer()
	// naive approach: locate last newline and compare suffix contains footer base string
	lastNL := -1
	for i := len(view2) - 1; i >= 0; i-- {
		if view2[i] == '\n' {
			lastNL = i
			break
		}
	}
	if lastNL == -1 {
		t.Fatalf("no newline found in view for footer detection")
	}
	footerLine := view2[lastNL+1:]
	if !contains(footerLine, footer) {
		t.Fatalf("footer line does not contain expected footer content")
	}
	// ensure footer unchanged after additional scroll up
	updatedUp, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	m = updatedUp.(model.RootModel)
	view3 := m.View()
	lastNL2 := -1
	for i := len(view3) - 1; i >= 0; i-- {
		if view3[i] == '\n' {
			lastNL2 = i
			break
		}
	}
	if lastNL2 == -1 {
		t.Fatalf("no newline in view3")
	}
	footerLine2 := view3[lastNL2+1:]
	if footerLine2 != footerLine {
		t.Fatalf("footer changed across scroll operations")
	}
	_ = time.Now() // placeholder to avoid unused imports if extended later
}

// simple key message helper reusing pattern from other tests (minimal no dependency)
// We'll define a tiny shim since tests can't import tea without adding heavy deps; root model returns itself.
// For real tests referencing tea.KeyMsg, adapt as needed.

// contains helper
func contains(s, sub string) bool { return len(sub) == 0 || indexOf(s, sub) >= 0 }
func indexOf(s, sub string) int {
	if len(sub) == 0 {
		return 0
	}
	// naive substring search
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

// expose buffer via method (need to add accessor if missing)
// We assume a Buffer() accessor exists; if not we'll add it.
