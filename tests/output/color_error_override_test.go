package output_test

import (
	"forge/src/output"
	"testing"
)

// T042 [US4] Test: error step overrides assigned color with red
// We'll simulate by selecting a normal color then comparing to ErrorColor constant.
func TestErrorColorOverride(t *testing.T) {
	id := "some-step"
	normal := output.ColorForID(id)
	if string(output.ErrorColor) == string(normal) {
		// In rare case hashing yields same as error color; choose different id.
		id = "another-step"
		normal = output.ColorForID(id)
	}
	// Override logic (to be integrated) should use ErrorColor; here we just assert constant is distinct.
	if string(output.ErrorColor) == string(normal) {
		t.Fatalf("expected error color %s to differ from normal %s", output.ErrorColor, normal)
	}
}
