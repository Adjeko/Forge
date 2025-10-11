package exec_test

import (
	"forge/src/exec"
	"forge/src/logging"
	"testing"
	"time"
)

// T054 [OBS] Test: monitor poll event sequence (start -> result)
func TestMonitorPollEventsSequence(t *testing.T) {
	logging.ResetEvents()
	m := exec.NewFakeMonitor("m-seq", 1*time.Second)
	s := exec.NewMonitorScheduler([]exec.Monitor{m}, 1*time.Second)
	go s.Start()
	defer s.Stop()
	// wait enough for at least one poll
	time.Sleep(1500 * time.Millisecond)
	events := logging.CapturedEvents()
	startFound := false
	resultFound := false
	for _, e := range events {
		if e.Name == logging.EventMonitorPollStart {
			startFound = true
		}
		if e.Name == logging.EventMonitorPollResult {
			resultFound = true
			if !startFound {
				t.Fatalf("result before start")
			}
		}
	}
	if !startFound || !resultFound {
		t.Fatalf("missing start/result events start=%v result=%v", startFound, resultFound)
	}
}
