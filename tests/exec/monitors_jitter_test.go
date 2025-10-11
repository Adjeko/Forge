package exec_test

import (
	"forge/src/exec"
	"testing"
	"time"
)

// T035 [US3] Test: jitter prevents synchronized poll (timestamps variance)
// We create multiple fake monitors and run scheduler; after some polls, verify their poll times are not identical.
// Allow tolerance of 500ms; we expect divergence due to jitter.
func TestMonitorJitterPreventsSynchronization(t *testing.T) {
	monitors := []exec.Monitor{
		exec.NewFakeMonitor("m1", 5*time.Second),
		exec.NewFakeMonitor("m2", 5*time.Second),
		exec.NewFakeMonitor("m3", 5*time.Second),
	}
	s := exec.NewMonitorScheduler(monitors, 5*time.Second)
	go s.Start()
	defer s.Stop()
	// Wait enough for at least one poll each.
	time.Sleep(6 * time.Second)
	// Collect poll times.
	var times []time.Time
	for _, m := range monitors {
		if fm, ok := m.(*exec.FakeMonitor); ok {
			lt := fm.LastPollTime()
			if lt.IsZero() {
				t.Fatalf("monitor %s did not poll within 6s", fm.ID())
			}
			times = append(times, lt)
		}
	}
	// Check variance: not all within 500ms window.
	var min, max time.Time
	for i, ts := range times {
		if i == 0 || ts.Before(min) {
			min = ts
		}
		if i == 0 || ts.After(max) {
			max = ts
		}
	}
	if max.Sub(min) < 500*time.Millisecond {
		t.Fatalf("expected jitter to spread poll times by >=500ms; got spread=%v", max.Sub(min))
	}
}
