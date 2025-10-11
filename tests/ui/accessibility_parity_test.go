package ui_test

import (
	"forge/src/ui/accessibility"
	"testing"
)

// T075 accessibility parity snapshot
func TestAccessibilityParitySnapshot(t *testing.T) {
	accessibility.Register(accessibility.HotkeyAction{ID: "cmd.run", Description: "Run Command", Keys: []string{"enter"}, ZoneID: "zone-command-run"})
	if !accessibility.ParityOK() {
		t.Fatalf("expected parity ok with registered action")
	}
}
