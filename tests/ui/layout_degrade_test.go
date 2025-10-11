package ui_test

import "testing"

// T070 layout degrade width<100
func Degrades(width int) bool { return width < 100 }

func TestLayoutDegrade(t *testing.T) {
	if !Degrades(99) {
		t.Fatalf("expected degrade at width 99")
	}
	if Degrades(120) {
		t.Fatalf("did not expect degrade at width 120")
	}
}
