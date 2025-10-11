package output_test

import (
	"forge/src/output"
	"testing"
)

// T082 [US4] Test: Palette exhaustion reuse strategy (more steps than palette size)
// Current behavior: modulo reuse. We assert that requesting > len(palette) colors returns repeats, not panics.
func TestPaletteExhaustionReuseStrategy(t *testing.T) {
	ids := []string{"s1", "s2", "s3", "s4", "s5", "s6", "s7", "s8", "s9", "s10"}
	seen := map[string]int{}
	for _, id := range ids {
		c := string(output.ColorForID(id))
		seen[c]++
	}
	// Expect at least one color used more than once since len(ids) > palette size (~6)
	multiple := false
	for _, count := range seen {
		if count > 1 {
			multiple = true
			break
		}
	}
	if !multiple {
		t.Fatalf("expected color reuse when palette exhausted; got unique count=%d paletteReuse=%v", len(seen), seen)
	}
}
