package ui_test

import (
	"forge/src/ui/accessibility"
	"forge/src/ui/model"
	"forge/src/ui/zones"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// T077 [ACC] Real test: Hotkey & BubbleZone action parity audit
func TestHotkeyZoneParity(t *testing.T) {
	accessibility.Reset()
	zones.Reset()
	m := model.NewRootModel()
	// ensure Init runs to append any parity warnings
	if cmd := m.Init(); cmd != nil {
		cmd()
	}
	actions := accessibility.List()
	zoneSet := map[string]bool{}
	for _, z := range zones.Zones() {
		zoneSet[z] = true
	}
	missing := 0
	for _, a := range actions {
		if !zoneSet[a.ZoneID] {
			missing++
		}
	}
	if missing > 0 {
		t.Fatalf("found %d actions missing zones", missing)
	}
	// Add scroll actions parity demonstration
	accessibility.Register(accessibility.HotkeyAction{ID: "scroll-up", Description: "Scroll up", Keys: []string{"pgup"}, ZoneID: "action:scroll-up"})
	accessibility.Register(accessibility.HotkeyAction{ID: "scroll-down", Description: "Scroll down", Keys: []string{"pgdown"}, ZoneID: "action:scroll-down"})
	zones.RegisterZone("action:scroll-up")
	zones.RegisterZone("action:scroll-down")
	km := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	m.Update(km) // toggle overlay - parity should still hold
	actions2 := accessibility.List()
	for _, a := range actions2 {
		if a.ZoneID == "" {
			t.Fatalf("action %s missing zone", a.ID)
		}
	}
}
