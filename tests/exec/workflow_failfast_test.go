package exec_test

import (
	"context"
	"forge/src/exec"
	"forge/src/output"
	"testing"
)

// T026 workflow stops on first failing step (simulate by injecting non-whitelisted second command)

func TestWorkflowFailFast(t *testing.T) {
	wf := &exec.Workflow{Steps: []exec.ExecutionStep{{Cmd: exec.Whitelist[0]}, {Cmd: exec.PrimitiveCommand{ID: "bad", Label: "bad", Cmd: "badcmd"}}}}
	buf := output.NewBuffer(0)
	failedAt, err := exec.RunWorkflow(context.Background(), wf, buf)
	if failedAt != 1 || err == nil {
		t.Fatalf("expected failure at second step with error; got failedAt=%d err=%v", failedAt, err)
	}
}
