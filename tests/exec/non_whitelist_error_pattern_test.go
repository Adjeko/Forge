package exec_test

import (
	"forge/src/exec"
	"strings"
	"testing"
)

// T072 non-whitelist error pattern
func TestNonWhitelistErrorPattern(t *testing.T) {
	err := exec.ValidateCommand("echo 'not allowed'")
	if err == nil {
		t.Fatalf("expected error for non-whitelist command")
	}
	if got := err.Error(); !strings.HasPrefix(got, "ERR_NON_WHITELIST:") {
		t.Fatalf("unexpected error pattern: %s", got)
	}
}
