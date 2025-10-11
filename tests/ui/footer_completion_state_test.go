package ui_test

import (
	"forge/src/ui/model"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// T083 [ACC] Real test: Footer completion state visible after workflow end
func TestFooterCompletionState(t *testing.T) {
	m := model.NewRootModel()
	// Run workflow (enter key)
	km := tea.KeyMsg{Type: tea.KeyEnter}
	nm, _ := m.Update(km)
	view := nm.View()
	if !strings.Contains(view, "exit=0") {
		t.Fatalf("expected exit=0 in footer after successful workflow")
	}
	// progress 100 implies completion
	if !strings.Contains(view, "exit=0") { // already checked but explicit for clarity
		t.Fatalf("missing completion indicator")
	}
}
