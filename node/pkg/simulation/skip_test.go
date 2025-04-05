package simulation

import (
	"testing"
)

// TestSkipped is a test that will be skipped when running with -short flag
func TestSkipped(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	// This test will always pass if not skipped
	t.Log("This test should be skipped when running with -short flag")
}
