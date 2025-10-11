package perf_test

import (
	"forge/src/output"
	"forge/src/ui/components"
	"testing"
)

// Benchmark buffer append and viewport rendering with 50k lines.
func BenchmarkBufferAppendAndViewport(b *testing.B) {
	buf := output.NewBuffer(60000)
	vp := components.NewOutputViewport(buf, 50)
	// Pre-fill 50k lines
	for i := 0; i < 50000; i++ {
		buf.Append("line")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Append("x")
		_ = vp.View()
	}
}
