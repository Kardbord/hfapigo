package testutils

import (
	"reflect"
	"testing"
)

// Ptr returns a pointer to the provided value.
func Ptr[T any](value T) *T {
	return &value
}

// RequireNil fails the test if value is not nil.
func RequireNil(t *testing.T, value any) {
	t.Helper()

	if value == nil {
		return
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		if v.IsNil() {
			return
		}
	}
	t.Fatalf("expected nil, got %v", value)
}

// RequireNotNil fails the test if value is nil.
func RequireNotNil(t *testing.T, value any) {
	t.Helper()

	if value == nil {
		t.Fatalf("expected non-nil value")
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		if v.IsNil() {
			t.Fatalf("expected non-nil value")
		}
	}
}
