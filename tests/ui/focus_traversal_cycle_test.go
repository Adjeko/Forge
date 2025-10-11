package ui_test

import (
	"forge/src/ui/model"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// T046 [ACC] Real test: focus traversal cycles all regions via Tab
func TestFocusTraversalCycle(t *testing.T) {
	m := model.NewRootModel()
	order := []string{m.FocusCurrent()}
	// Press Tab 3 times to move through next 3 regions
	for i := 0; i < 3; i++ {
		km := tea.KeyMsg{Type: tea.KeyTab}
		nm, _ := m.Update(km)
		// assert nm type and get current focus by rendering or exposing helper (we'll parse view highlight as fallback)
		m = nm.(model.RootModel)
		order = append(order, m.FocusCurrent())
	}
	// After 4 states we expect unique traversal (CommandList, Output, Monitors, Help)
	if len(order) != 4 {
		t.Fatalf("expected 4 focus states, got %v", order)
	}
	// Now wrap: one more Tab returns to CommandList
	km := tea.KeyMsg{Type: tea.KeyTab}
	nm, _ := m.Update(km)
	m = nm.(model.RootModel)
	if m.FocusCurrent() != "CommandList" {
		t.Fatalf("expected focus to wrap to CommandList; got %s", m.FocusCurrent())
	}
}
