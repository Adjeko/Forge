package ui_test

import (
	"forge/src/logging"
	"forge/src/ui/model"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// T058 [OBS] Test: UI focus change & help toggle events captured
func TestUIFocusAndHelpEvents(t *testing.T) {
	logging.ResetEvents()
	m := model.NewRootModel()
	// Tab through two focuses
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = nm.(model.RootModel)
	nm, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = nm.(model.RootModel)
	// Toggle help overlay
	nm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m = nm.(model.RootModel)
	events := logging.CapturedEvents()
	focusCount := 0
	helpToggle := 0
	for _, e := range events {
		if e.Name == logging.EventUIFocusChange {
			focusCount++
		}
		if e.Name == logging.EventUIHelpToggle {
			helpToggle++
		}
	}
	if focusCount < 2 {
		t.Fatalf("expected >=2 focus change events got %d", focusCount)
	}
	if helpToggle != 1 {
		t.Fatalf("expected 1 help toggle event got %d", helpToggle)
	}
}
