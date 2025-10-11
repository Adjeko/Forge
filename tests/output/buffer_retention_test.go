package output_test

import (
	"forge/src/output"
	"testing"
)

// T068 viewport retention
func TestViewportRetention(t *testing.T) {
	b := output.NewBuffer(100000)
	for i := 0; i < 6000; i++ {
		b.Append("line")
	}
	slice := b.ViewportSlice(5000)
	if len(slice) != 5000 {
		t.Fatalf("expected 5000 lines retained slice, got %d", len(slice))
	}
}
