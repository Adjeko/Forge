package exec_test

import (
	"forge/src/ui/model"
	"os"
	"testing"
)

// Verify FORGE_WORKFLOW_STEPS env flag controls number of steps in workflow.
func TestWorkflowStepsEnvFlag(t *testing.T) {
	os.Setenv("FORGE_WORKFLOW_STEPS", "3")
	m := model.NewRootModel()
	if m.WorkflowLen() != 3 {
		t.Fatalf("expected workflow length 3, got %d", m.WorkflowLen())
	}
	os.Unsetenv("FORGE_WORKFLOW_STEPS")
	m2 := model.NewRootModel()
	if m2.WorkflowLen() != 1 {
		t.Fatalf("expected default workflow length 1, got %d", m2.WorkflowLen())
	}
}
