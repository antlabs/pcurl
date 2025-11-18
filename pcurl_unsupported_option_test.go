package pcurl

import (
	"testing"
)

func Test_UnsupportedOption(t *testing.T) {
	_, err := ParseAndRequest(`curl --hahaha`)
	if err == nil {
		t.Fatalf("expected error for unsupported option, got nil")
	}
}
