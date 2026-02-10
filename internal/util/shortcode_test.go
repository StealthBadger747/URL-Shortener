package util

import (
	"strings"
	"testing"
)

func TestRandomCodeLength(t *testing.T) {
	code, err := RandomCode(8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(code) != 8 {
		t.Fatalf("expected length 8, got %d", len(code))
	}
}

func TestRandomCodeCharset(t *testing.T) {
	code, err := RandomCode(32)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range code {
		if !strings.ContainsRune(alphabet, r) {
			t.Fatalf("unexpected rune %q", r)
		}
	}
}
