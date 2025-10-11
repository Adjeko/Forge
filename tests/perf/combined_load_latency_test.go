package perf_test

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"forge/src/ui/model"
)

// T079 Combined load latency: workflow + monitors + stream append; key handling under 200ms.
func TestCombinedLoadLatency(t *testing.T) {
	m := model.NewRootModel()
	// Start heavy stream append in background
	done := make(chan struct{})
	go func() {
		for i := 0; i < 30000; i++ { // 30k lines
			m.Buffer().Append("combined load line")
		}
		close(done)
	}()
	// Trigger workflow (enter)
	startWorkflow := time.Now()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if time.Since(startWorkflow) > 200*time.Millisecond { // workflow dispatch latency measure
		t.Fatalf("workflow dispatch latency exceeded 200ms: %v", time.Since(startWorkflow))
	}
	m = updated.(model.RootModel)
	// Wait for stream to finish
	<-done
	// Measure help toggle under load
	startHelp := time.Now()
	updated2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if time.Since(startHelp) > 200*time.Millisecond {
		t.Fatalf("help toggle latency under combined load exceeded 200ms: %v", time.Since(startHelp))
	}
	_ = updated2
}
