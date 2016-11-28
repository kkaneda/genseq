package testutils

import (
	"errors"
	"testing"
)

// TestIsError
func TestIsError(t *testing.T) {
	err := errors.New("message")
	if !IsError(err, "message") {
		t.Errorf("unexpected failure: %v", err)
	}
	if IsError(err, "bad_message") {
		t.Errorf("unexpected success: %v", err)
	}
	if IsError(nil, "message") {
		t.Errorf("unexpected success")
	}
}
