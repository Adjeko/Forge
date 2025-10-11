package ui_test

import (
	"forge/src/ui/model"
	"testing"
)

// T074 composition-first smoke
func TestCompositionSmoke(t *testing.T) {
	m := model.NewRootModel()
	if m.View() == "" {
		t.Fatalf("expected non-empty view")
	}
}
