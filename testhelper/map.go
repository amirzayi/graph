package testhelper

import "testing"

func MapMustContainKeys[T comparable](t *testing.T, m map[T]any, keys ...T) {
	t.Helper()

	var notExistedKeys []T
	for _, key := range keys {
		if _, exists := m[key]; !exists {
			notExistedKeys = append(notExistedKeys, key)
		}
	}

	if len(notExistedKeys) > 0 {
		t.Fatalf("map doesnt contains %v keys", notExistedKeys)
	}
}
