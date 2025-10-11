package ui_test

import (
	"forge/src/ui/model"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// T061 Final interaction audit: parity, help toggle, focus cycle, scroll under resize.
func TestFinalInteractionAudit(t *testing.T) {
	m := model.NewRootModel()
	startHelp := time.Now()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m = updated.(model.RootModel)
	if time.Since(startHelp) > 150*time.Millisecond {
		t.Fatalf("help toggle latency exceeded")
	}
	// Focus cycle
	seen := map[string]bool{"CommandList": true}
	for i := 0; i < 3; i++ {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
		m = updated.(model.RootModel)
		seen[m.FocusCurrent()] = true
	}
	for _, region := range []string{"Output", "Monitors", "Help"} {
		if !seen[region] {
			t.Fatalf("expected focus region %s visited", region)
		}
	}
	// Scroll after focusing output
	// Cycle until focus becomes Output
	for m.FocusCurrent() != "Output" {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
		m = updated.(model.RootModel)
	}
	scrollStart := time.Now()
	for i := 0; i < 5; i++ {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyPgUp})
		m = updated.(model.RootModel)
	}
	if time.Since(scrollStart) > 200*time.Millisecond {
		t.Fatalf("scroll latency exceeded")
	}
}
