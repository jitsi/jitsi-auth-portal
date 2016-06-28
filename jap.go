package jap

import (
	"crypto/rsa"
	"errors"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"golang.org/x/oauth2/jws"
)

const (
	clientIDKey int = iota
)

var (
	errPermissionDenied = errors.New("Permission denied")
)

// CIDFromContext returns the client ID bound to the context, if any.
func CIDFromContext(ctx context.Context) (cid string, ok bool) {
	cid, ok = ctx.Value(clientIDKey).(string)
	return
}

// NewCIDContext returns a copy of the parent context and associates it with a
// client id.
func NewCIDContext(ctx context.Context, cid string) context.Context {
	return context.WithValue(ctx, clientIDKey, cid)
}

func writeError(ctx context.Context, w http.ResponseWriter, msg string, status int) {
	tr, ok := trace.FromContext(ctx)
	if ok {
		tr.LazyPrintf(msg)
		tr.SetError()
	}
	http.Error(w, msg, status)
}

// PermissionChecker is a function that's used for checking if the email
// associated with a given token has permission to perform some action.
type PermissionChecker func(tok string) (bool, error)

func signJWT(
	ctx context.Context,
	claims jws.ClaimSet,
	key *rsa.PrivateKey,
	permCheck PermissionChecker) (tok string, err error) {
	// Assert that we actually get a key. We don't want bugs that result in nil
	// keys to go unnoticed; we want them to break everything. This would probably
	// happen in the crypto functions anyways, but I want it to be testable.
	if key == nil {
		panic("got nil RSA private key; something is very, very wrong.")
	}

	tr, ok := trace.FromContext(ctx)

	header := jws.Header{
		Algorithm: "RS256",
	}
	if ok {
		tr.LazyPrintf("Signing JWT…")
	}
	tok, err = jws.Encode(&header, &claims, key)
	if err != nil {
		return tok, err
	}
	if ok {
		tr.LazyPrintf("Done signing JWT.")
	}
	if permCheck != nil {
		tr.LazyPrintf("Checking permissions…")
		// TODO(ssw): Retry if there's an error?
		if ok, err := permCheck(tok); !ok || err != nil {
			return tok, errPermissionDenied
		}
	}
	return tok, nil
}
