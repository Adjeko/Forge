package exec_test

import (
	"context"
	"testing"
	"time"

	"forge/src/exec"
	"forge/src/output"
)

// T090 Monitor/workflow separation: failing monitor state must not alter workflow exit or progress percent.
func TestMonitorWorkflowSeparation(t *testing.T) {
	// Create a workflow with one trivial whitelisted command (simulate instant success by using a harmless command)
	wf := &exec.Workflow{Steps: []exec.ExecutionStep{{Cmd: exec.Whitelist[0]}}}
	buf := output.NewBuffer(10)
	// Create a monitor that will quickly go to failed state
	failing := exec.NewFakeMonitorWithType("failing-monitor", 10*time.Millisecond, "ping")
	scheduler := exec.NewMonitorScheduler([]exec.Monitor{failing}, 10*time.Millisecond)
	go scheduler.Start()
	// Execute workflow
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	failedAt, err := exec.RunWorkflow(ctx, wf, buf)
	if err != nil && failedAt >= 0 {
		t.Fatalf("workflow should succeed independently of monitor failure, got err=%v failedAt=%d", err, failedAt)
	}
	// Wait briefly to allow monitor to poll and potentially set failure state
	time.Sleep(50 * time.Millisecond)
	// At this point workflow should have completed; monitor failure should not inject buffer failure lines except its own log semantics
	// Validate buffer does not contain synthetic workflow failure text
	for _, line := range buf.ViewportSlice(buf.Len()) {
		if line == "MONITOR INDUCED WORKFLOW FAILURE" {
			t.Fatalf("found forbidden coupling line in buffer")
		}
	}
}
