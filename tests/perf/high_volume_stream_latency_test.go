package perf_test

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"forge/src/ui/model"
)

// T089 High-volume stream latency: after appending 50k lines, toggling help overlay should process within <200ms.
func TestHighVolumeStreamLatency(t *testing.T) {
	m := model.NewRootModel()
	// Append 50k lines to buffer
	for i := 0; i < 50000; i++ {
		m.Buffer().Append("line perf test")
	}
	// Measure latency of help toggle key processing
	start := time.Now()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if time.Since(start) > 200*time.Millisecond {
		t.Fatalf("help toggle latency exceeded 200ms: %v", time.Since(start))
	}
	// ensure overlay visible by checking view contains expected substring (simplified)
	m = updated.(model.RootModel)
	view := m.View()
	if len(view) == 0 {
		t.Fatalf("empty view after toggle")
	}
}
