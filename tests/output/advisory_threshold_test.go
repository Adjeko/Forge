package output_test

import (
	"forge/src/output"
	"testing"
)

// T073 advisory threshold warning trigger
func TestAdvisoryThreshold(t *testing.T) {
	b := output.NewBuffer(100000)
	for i := 0; i < 100001; i++ {
		b.Append("line")
	}
	if !b.WarningIssued() {
		t.Fatalf("expected warning issued after threshold exceed")
	}
	// Verify advisory line appended with prefix
	lines := b.ViewportSlice(b.Len())
	found := false
	for _, l := range lines {
		if len(l) >= 9 && l[:9] == "ADVISORY:" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("advisory line with 'ADVISORY:' prefix not found in buffer lines=%d", len(lines))
	}
}
