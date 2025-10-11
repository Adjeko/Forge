package output_test

import (
	"forge/src/output"
	"testing"
)

// T017 color assignment deterministic
func TestColorDeterminism(t *testing.T) {
	c1 := output.ColorForID("git-status")
	c2 := output.ColorForID("git-status")
	if c1 != c2 {
		t.Fatalf("expected same color for identical primitive ID")
	}
}
