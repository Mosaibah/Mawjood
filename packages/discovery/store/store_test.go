package store

import (
	"testing"
)

// Simple test to verify the store interface implementation
func TestStoreInterface(t *testing.T) {
	// This test verifies that our ContentData struct implements the Interface
	var _ Interface = &ContentData{}

	// Test passes if compilation succeeds
	t.Log("ContentData successfully implements the Interface")
}
