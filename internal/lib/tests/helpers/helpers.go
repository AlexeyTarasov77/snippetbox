package helpers

import "testing"


func Equal[T comparable](t *testing.T, actual T, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}