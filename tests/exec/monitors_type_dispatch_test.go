package exec_test

import (
	"forge/src/exec"
	"testing"
	"time"
)

// T081 [US3] Test: Monitor type dispatch & script constraints (ping vs script)
// We will create two fake monitors with declared types and ensure scheduler polls both and that script monitor enforces constraint (e.g., script path non-empty).
func TestMonitorTypeDispatch(t *testing.T) {
	ping := exec.NewFakeMonitorWithType("ping1", 5*time.Second, "ping")
	script := exec.NewFakeMonitorWithType("script1", 5*time.Second, "script")
	script.SetScriptPath("some/path.sh")
	s := exec.NewMonitorScheduler([]exec.Monitor{ping, script}, 5*time.Second)
	go s.Start()
	defer s.Stop()
	time.Sleep(6 * time.Second)
	if ping.PollCount() == 0 || script.PollCount() == 0 {
		t.Fatalf("expected both monitors to have polled; ping=%d script=%d", ping.PollCount(), script.PollCount())
	}
	// Constraint: script path must be set.
	if script.ScriptPath() == "" {
		t.Fatalf("expected script monitor to retain script path")
	}
}
