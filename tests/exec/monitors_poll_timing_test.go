package exec_test

import (
	"forge/src/exec"
	"testing"
	"time"
)

// T034 [US3] Test: monitor poll updates state within <=6s (default interval 5s)
// Assumption: We will implement a monitor scheduler with default interval = 5s.
// This test uses a fake monitor whose Poll method increments a counter and sets lastPoll time.
// After starting the scheduler and advancing time (simulated via sleeping in test), we assert at least one poll occurred within 6s.
func TestMonitorPollUpdatesWithinSixSeconds(t *testing.T) {
	m := exec.NewFakeMonitor("fake1", 5*time.Second)
	s := exec.NewMonitorScheduler([]exec.Monitor{m}, 5*time.Second)
	go s.Start() // non-blocking
	defer s.Stop()
	// Wait up to 6s.
	time.Sleep(6 * time.Second)
	if m.PollCount() == 0 {
		// Provide detailed failure to help debugging.
		last := m.LastPollTime()
		elapsed := time.Since(last)
		if last.IsZero() {
			elapsed = 0
		}
		// We expect at least one poll under default interval + grace.
		if elapsed > 6*time.Second || last.IsZero() {
			// emphasise requirement FR-018 cadence.
			// Using t.Fatalf ensures immediate clarity.
			t.Fatalf("expected >=1 poll within 6s (interval=5s); pollCount=%d lastPoll=%v elapsed=%v", m.PollCount(), last, elapsed)
		}
	}
}
