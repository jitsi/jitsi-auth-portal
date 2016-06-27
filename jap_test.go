package jap

import (
	"testing"

	"golang.org/x/net/context"
)

func TestCIDContext(t *testing.T) {
	ctx := NewCIDContext(context.Background(), "TESTSID")
	sid, ok := CIDFromContext(ctx)
	if !ok || sid != "TESTSID" {
		t.Error("Could not find CID in context.")
	}
}
