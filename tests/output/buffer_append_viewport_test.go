package output_test

import (
	"forge/src/output"
	"testing"
)

// T018 buffer append & viewport integrity
func TestBufferAppendViewportIntegrity(t *testing.T) {
	b := output.NewBuffer(0)
	for i := 0; i < 120; i++ {
		b.Append("line")
	}
	if b.Len() != 120 {
		t.Fatalf("expected 120 lines, got %d", b.Len())
	}
	slice := b.ViewportSlice(50)
	if len(slice) != 50 {
		t.Fatalf("expected last 50 lines slice, got %d", len(slice))
	}
}
