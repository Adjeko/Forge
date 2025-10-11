package ui_test

import (
	"forge/src/ui/model"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// T045 [ACC] Real test: help overlay toggles & renders within ≤150 ms
func TestHelpOverlayToggleLatency(t *testing.T) {
	m := model.NewRootModel()
	start := time.Now()
	// simulate '?' key press
	km := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	nm, _ := m.Update(km)
	rendered := nm.View()
	if !strings.Contains(rendered, "Help (? to hide)") {
		t.Fatalf("expected help overlay title in view after toggle")
	}
	elapsed := time.Since(start)
	if elapsed > 150*time.Millisecond {
		t.Fatalf("overlay toggle exceeded 150ms elapsed=%v", elapsed)
	}
}

// removed naive helpers; using strings.Contains instead.
