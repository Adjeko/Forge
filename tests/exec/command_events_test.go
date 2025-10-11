package exec_test

import (
	"context"
	"forge/src/exec"
	"forge/src/logging"
	"forge/src/output"
	"testing"
)

// T053 [OBS] Test: command start/end events emitted with correct fields
func TestCommandEventsEmitted(t *testing.T) {
	logging.ResetEvents()
	buf := output.NewBuffer(0)
	primitive := exec.Whitelist[0]
	exit, err := exec.RunSingleCommand(context.Background(), primitive, buf)
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if exit != 0 {
		t.Fatalf("expected exit 0 got %d", exit)
	}
	events := logging.CapturedEvents()
	startIdx := -1
	endIdx := -1
	for i, e := range events {
		if e.Name == logging.EventCommandStart {
			startIdx = i
		}
		if e.Name == logging.EventCommandEnd {
			endIdx = i
		}
	}
	if startIdx == -1 || endIdx == -1 || endIdx < startIdx {
		t.Fatalf("missing or misordered command events: start=%d end=%d", startIdx, endIdx)
	}
}
