package perf_test

import (
	"forge/src/ui/model"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// T047 [ACC] Real test: output scroll maintains ≤200 ms input latency
func TestOutputScrollLatency(t *testing.T) {
	m := model.NewRootModel()
	// Move focus to Output (tab twice)
	for i := 0; i < 2; i++ {
		nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
		m = nm.(model.RootModel)
	}
	start := time.Now()
	// Perform multiple PgUp scrolls
	for i := 0; i < 5; i++ {
		nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgUp})
		m = nm.(model.RootModel)
	}
	elapsed := time.Since(start)
	if elapsed > 200*time.Millisecond {
		t.Fatalf("scroll latency exceeded 200ms: %v", elapsed)
	}
	// Scroll down
	start2 := time.Now()
	for i := 0; i < 5; i++ {
		nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
		m = nm.(model.RootModel)
	}
	if time.Since(start2) > 200*time.Millisecond {
		t.Fatalf("scroll down latency exceeded 200ms")
	}
}
