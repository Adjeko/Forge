package output_test

import (
	"forge/src/output"
	"testing"
)

// T041 [US4] Test: same primitive sequence yields identical color set across two runs
func TestColorSequenceStabilityAcrossRuns(t *testing.T) {
	ids := []string{"alpha", "beta", "gamma", "delta"}
	first := make([]string, len(ids))
	second := make([]string, len(ids))
	for i, id := range ids {
		first[i] = string(output.ColorForID(id))
	}
	for i, id := range ids {
		second[i] = string(output.ColorForID(id))
	}
	for i := range ids {
		if first[i] != second[i] {
			t.Fatalf("color instability for id %s: %s vs %s", ids[i], first[i], second[i])
		}
	}
}
