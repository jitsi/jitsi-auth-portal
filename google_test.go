package jap

import (
	"testing"

	"golang.org/x/net/context"
)

func TestGoogleHandlerPanicsWithoutCID(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected GoogleLogin to panic if CID missing from context")
		}
	}()
	_ = GoogleLogin(context.Background(), nil, nil)
}

func TestGoogleHandlerDoesNotPanicWithCID(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Did not expect GoogleLogin to panic if provided with a CID")
		}
	}()
	if h := GoogleLogin(NewCIDContext(context.Background(), "TESTSID"), nil, nil); h == nil {
		t.Error("Expected GoogleLogin to return a handler, but got nil")
	}
}
