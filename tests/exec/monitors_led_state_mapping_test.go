package exec_test

import (
	"forge/src/exec"
	"testing"
	"time"
)

// T076 [US3] Test: LED state accuracy classification (OK/Failed/Hanging)
// We simulate state transitions on a fake monitor and assert LED mapping output.
func TestMonitorLEDStateMapping(t *testing.T) {
	m := exec.NewFakeMonitor("led1", 5*time.Second)
	// Simulate transitions using helper methods (will be added on FakeMonitor).
	m.SetResult(true, nil) // OK
	if exec.MonitorLEDColor(m.State()) != "green" {
		t.Fatalf("expected green for OK; got %s", exec.MonitorLEDColor(m.State()))
	}
	m.SetResult(false, nil) // Failed
	if exec.MonitorLEDColor(m.State()) != "red" {
		t.Fatalf("expected red for Failed; got %s", exec.MonitorLEDColor(m.State()))
	}
	// Hanging: simulate lastPoll older than 2*interval and checking state.
	m.SetChecking()
	m.SetLastPoll(time.Now().Add(-12 * time.Second))
	if exec.MonitorLEDColor(m.State()) != "yellow" {
		t.Fatalf("expected yellow for Hanging; got %s", exec.MonitorLEDColor(m.State()))
	}
}
