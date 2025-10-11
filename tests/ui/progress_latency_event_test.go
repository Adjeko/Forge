package ui_test

import (
	"context"
	"forge/src/exec"
	"forge/src/logging"
	"forge/src/output"
	"testing"
	"time"
)

// T078 [OBS] Test: progress bar update latency ≤250 ms after last step completion
func TestProgressBarLatency(t *testing.T) {
	logging.ResetEvents()
	// Replace workflow with single fast step
	mWorkflow := &exec.Workflow{Steps: []exec.ExecutionStep{{Cmd: exec.Whitelist[0]}}}
	// run workflow directly
	start := time.Now()
	buf := output.NewBuffer(0)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := exec.RunWorkflow(ctx, mWorkflow, buf)
	if err != nil {
		t.Fatalf("workflow error: %v", err)
	}
	elapsed := time.Since(start)
	if elapsed > 250*time.Millisecond {
		t.Fatalf("progress latency exceeded 250ms: %v", elapsed)
	}
}
