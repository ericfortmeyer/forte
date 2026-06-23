package archive

import (
	"errors"
	"testing"
)

func TestNewSkippableError(t *testing.T) {
	msg := "test error message"
	err := NewSkippableError(msg)

	if err == nil {
		t.Fatal("NewSkippableError returned nil")
	}
	if err.message != msg {
		t.Errorf("expected message %q, got %q", msg, err.message)
	}
}

func TestSkippableError_Error(t *testing.T) {
	msg := "validation failed"
	err := NewSkippableError(msg)

	if err.Error() != msg {
		t.Errorf("expected %q, got %q", msg, err.Error())
	}
}

func TestIsSkippable_WithNil(t *testing.T) {
	if IsSkippable(nil) {
		t.Error("IsSkippable(nil) should return false")
	}
}

func TestIsSkippable_WithSkippableError(t *testing.T) {
	err := NewSkippableError("archive not found")
	if !IsSkippable(err) {
		t.Error("IsSkippable should return true for SkippableError")
	}
}

func TestIsSkippable_WithNonSkippableError(t *testing.T) {
	err := errors.New("some other error")
	if IsSkippable(err) {
		t.Error("IsSkippable should return false for non-SkippableError")
	}
}

func TestIsSkippable_WithWrappedSkippableError(t *testing.T) {
	original := NewSkippableError("missing archive")
	wrapped := errors.Join(errors.New("context"), original)

	if !IsSkippable(wrapped) {
		t.Error("IsSkippable should work with wrapped errors via errors.As")
	}
}
