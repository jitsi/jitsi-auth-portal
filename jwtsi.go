package jwtsi

import (
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
)

const (
	clientIDKey int = iota
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
