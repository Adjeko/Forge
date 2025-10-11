package exec_test

import (
	"forge/src/exec"
	"testing"
	"time"
)

// T080 [US3] Test: Configurable interval override via --monitor-interval (set to 2s)
// Since CLI flag parsing not implemented yet, we assume a constructor that accepts interval.
func TestMonitorIntervalOverride(t *testing.T) {
	m := exec.NewFakeMonitor("override1", 2*time.Second)
	s := exec.NewMonitorScheduler([]exec.Monitor{m}, 2*time.Second)
	go s.Start()
	defer s.Stop()
	// Sleep 5s -> expect >=2 polls (approx 0,2,4)
	time.Sleep(5 * time.Second)
	if m.PollCount() < 2 {
		t.Fatalf("expected at least 2 polls in 5s with 2s interval; got %d", m.PollCount())
	}
}
