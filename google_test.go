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
	_ = GoogleLogin(context.Background(), nil)
}
