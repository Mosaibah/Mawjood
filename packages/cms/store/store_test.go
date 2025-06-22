package store

import (
	"testing"
)

func TestStoreInterface(t *testing.T) {
	var _ Interface = &ContentData{}

	t.Log("ContentData successfully implements the Interface")
}
