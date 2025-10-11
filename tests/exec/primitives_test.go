package exec_test

import (
	"forge/src/exec"
	"testing"
)

func TestIsWhitelisted(t *testing.T) {
	if !exec.IsWhitelisted("cmd /C echo ok") {
		t.Fatalf("expected echo ok command to be whitelisted")
	}
	if exec.IsWhitelisted("rm -rf /") {
		t.Fatalf("expected dangerous command not whitelisted")
	}
}
