package jap

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/net/trace"
)

// GoogleLogin returns a handler which attempts to extract a client ID from its
// context and sends the information to Google to validate the user. If no
// client ID exists in the context it panics.
//
// The handler may return one of the following errors:
//
//   400 BadRequest          – If the id_token form param is missing.
//   408 RequestTimeout      – If the contexts deadline was exceeded.
//   500 InternalServerError – If the upstream returns a response we don't understand.
//   502 BadGateway          – If an upstream service fails to respond for another reason.
func GoogleLogin(ctx context.Context, key *rsa.PrivateKey) func(http.ResponseWriter, *http.Request) {
	cid, ok := CIDFromContext(ctx)
	if !ok {
		panic("No client ID found in the context")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tr := trace.New("jap.tokenlogin", r.URL.Path)
		defer tr.Finish()
		ctx = trace.NewContext(ctx, tr)

		idtoken := r.FormValue("id_token")
		if idtoken == "" {
			writeError(ctx, w, "id_token login param missing", http.StatusBadRequest)
			return
		}

		// TODO: We should skip the Post request and verify the JWT signature with
		//       Google's public key.
		tr.LazyPrintf("Starting Google token validation…")
		inforesp, err := ctxhttp.PostForm(ctx, http.DefaultClient,
			"https://www.googleapis.com/oauth2/v3/tokeninfo",
			url.Values{
				"id_token": []string{idtoken},
			},
		)
		if err != nil {
			switch err {
			case context.DeadlineExceeded:
				writeError(ctx, w, "The deadline was exceeded", http.StatusRequestTimeout)
			default:
				writeError(ctx, w, "Upstream request failed", http.StatusBadGateway)
			}
			return
		}
		if inforesp.StatusCode < 200 || inforesp.StatusCode > 299 {
			writeError(ctx, w, "Unexpected response from upstream", http.StatusInternalServerError)
			return
		}
		tr.LazyPrintf("Received claims from Google…")
		claims := struct {
			Aud          string `json:"aud"`
			Email        string `json:"email"`
			HostedDomain string `json:"hd"`
			Locale       string `json:"locale"`
		}{}
		if err := json.NewDecoder(inforesp.Body).Decode(&claims); err != nil {
			writeError(ctx, w, "Error decoding upstream response", http.StatusInternalServerError)
			return
		}
		if claims.Aud != cid {
			writeError(ctx, w, "Error decoding upstream response", http.StatusInternalServerError)
			return
		}
	}
}
