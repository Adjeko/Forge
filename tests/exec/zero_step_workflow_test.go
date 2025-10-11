package exec_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"forge/src/exec"
	"forge/src/logging"
	"forge/src/output"
)

// T084 Zero-step workflow validation: running empty workflow should return an error immediately
// and emit workflow start/end events without panic.
func TestZeroStepWorkflowValidation(t *testing.T) {
	logging.ResetEvents()
	wf := &exec.Workflow{Steps: []exec.ExecutionStep{}} // zero steps
	buf := output.NewBuffer(100)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	failedAt, err := exec.RunWorkflow(ctx, wf, buf)
	if err == nil {
		t.Fatalf("expected error for zero-step workflow, got nil")
	}
	if !errors.Is(err, exec.ErrEmptyWorkflow) {
		t.Fatalf("expected ErrEmptyWorkflow, got %v", err)
	}
	if failedAt != -1 {
		t.Fatalf("expected failedAt -1 for structural validation failure, got %d", failedAt)
	}
	events := logging.CapturedEvents()
	if len(events) < 2 {
		t.Fatalf("expected at least 2 events (start/end), got %d", len(events))
	}
	if events[0].Name != logging.EventWorkflowStart {
		t.Fatalf("first event should be workflow start, got %s", events[0].Name)
	}
	last := events[len(events)-1]
	if last.Name != logging.EventWorkflowEnd {
		t.Fatalf("last event should be workflow end, got %s", last.Name)
	}
	// ensure attributes include failedIndex -1
	foundFailedIndex := false
	for i := 0; i < len(last.Attrs)-1; i += 2 { // slog key/value pairs
		if key, ok := last.Attrs[i].(string); ok && key == "failedIndex" {
			foundFailedIndex = true
			if vi, ok2 := last.Attrs[i+1].(int); ok2 && vi != -1 {
				t.Fatalf("expected failedIndex -1, got %d", vi)
			}
		}
	}
	if !foundFailedIndex {
		t.Fatalf("failedIndex attribute not found in workflow end event")
	}
}
