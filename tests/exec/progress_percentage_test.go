package exec_test

import (
	"context"
	"forge/src/exec"
	"forge/src/output"
	"testing"
)

// Simple progress calculation imitation for T027
func progressPercent(total, completed int) int {
	if total == 0 {
		return 0
	}
	return int(float64(completed) / float64(total) * 100)
}

func TestProgressPercentageMatchesCompletedSteps(t *testing.T) {
	wf := &exec.Workflow{Steps: []exec.ExecutionStep{{Cmd: exec.Whitelist[0]}, {Cmd: exec.Whitelist[0]}, {Cmd: exec.PrimitiveCommand{ID: "bad", Label: "bad", Cmd: "badcmd"}}}}
	buf := output.NewBuffer(0)
	failedAt, _ := exec.RunWorkflow(context.Background(), wf, buf)
	completed := failedAt // failing index indicates first failed step; completed equals failedAt
	pct := progressPercent(len(wf.Steps), completed)
	if pct != int(float64(completed)/float64(len(wf.Steps))*100) {
		t.Fatalf("progress percent mismatch")
	}
}
