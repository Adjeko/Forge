package exec_test

import (
	"forge/src/exec"
	"testing"
	"time"
)

// T071 hanging detection simulation
func TestHangingDetection(t *testing.T) {
	step := exec.ExecutionStep{Start: time.Now().Add(-2 * time.Minute), LastOutput: time.Now().Add(-61 * time.Second)}
	if !step.IsHanging(time.Now()) {
		t.Fatalf("expected step hanging after 61s inactivity")
	}
}
