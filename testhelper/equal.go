package testhelper

import "testing"

func MustEqual[T comparable](t *testing.T, a, b T) {
	t.Helper()
	if a != b {
		t.Fatalf("%v and %v must be equal", a, b)
	}
}

func MustNotEqual[T comparable](t *testing.T, a, b T) {
	t.Helper()
	if a == b {
		t.Fatalf("%v and %v must not be equal", a, b)
	}
}

func MustNilError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("error is not nil: %v", err)
	}
}
