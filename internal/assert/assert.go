package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectdSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectdSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectdSubstring)
	}
}
