package exec_test

import (
	"forge/src/exec"
	"testing"
)

// T016 whitelist validation rejects non-listed command
func TestWhitelistValidationRejects(t *testing.T) {
	if err := exec.ValidateCommand("rm -rf /"); err == nil {
		// ValidateCommand should produce error for non-whitelist
		t.Fatalf("expected validation error for non-whitelisted command")
	}
}
