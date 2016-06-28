package jap

import (
	"testing"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/jws"
)

func TestCIDContext(t *testing.T) {
	ctx := NewCIDContext(context.Background(), "TESTSID")
	sid, ok := CIDFromContext(ctx)
	if !ok || sid != "TESTSID" {
		t.Error("Could not find CID in context.")
	}
}

func TestSignJWTPanicsOnNilKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected signJWT to panic with nil private key")
		}
	}()
	signJWT(context.Background(), jws.ClaimSet{}, nil, nil)
}
