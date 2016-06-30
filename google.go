package jap

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/net/trace"
	"golang.org/x/oauth2/jws"
)

// GoogleLogin returns a handler which attempts to extract a client ID from its
// context and sends the information to Google to validate the user. If no
// client ID exists in the context it panics.
//
// The handler may return one of the following errors:
//
//   400 BadRequest          – If the id_token form param is missing.
//   401 StatusUnauthorzed   — If the permCheck function returns false.
//   408 RequestTimeout      – If the contexts deadline was exceeded.
//   500 InternalServerError – If the upstream returns a response we don't understand.
//   502 BadGateway          – If an upstream service fails to respond for another reason.
func GoogleLogin(ctx context.Context, key *rsa.PrivateKey, permCheck PermissionChecker) func(http.ResponseWriter, *http.Request) {
	cid, ok := CIDFromContext(ctx)
	if !ok {
		panic("No client ID found in the context")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tr := trace.New("jap.GoogleLogin", r.URL.Path)
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
		defer inforesp.Body.Close()
		if inforesp.StatusCode < 200 || inforesp.StatusCode > 299 {
			writeError(ctx, w, "Unexpected response from upstream", http.StatusInternalServerError)
			return
		}
		tr.LazyPrintf("Received claims from Google.")
		meta := struct {
			Aud          string `json:"aud"`
			Email        string `json:"email"`
			Verified     string `json:"email_verified"` // For some reason Google makes this a string.
			HostedDomain string `json:"hd"`
			Locale       string `json:"locale"`
		}{}
		tr.LazyPrintf("Decoding claims from Google…")
		if err := json.NewDecoder(inforesp.Body).Decode(&meta); err != nil {
			writeError(ctx, w, "Error decoding upstream response", http.StatusInternalServerError)
			return
		}
		if meta.Aud != cid || ((meta.Email == "" || meta.Verified != "true") && meta.HostedDomain == "") {
			writeError(ctx, w, "Error decoding upstream response", http.StatusInternalServerError)
			return
		}
		tr.LazyPrintf("Decoded claims from Google.")

		// Generate the JWT
		iat := time.Now().Unix()
		claims := jws.ClaimSet{
			Scope:         "login",
			Iss:           meta.Email,
			Exp:           iat + (5 * 60), // TODO(ssw): Make the expiration time configurable.
			Iat:           iat,
			PrivateClaims: map[string]interface{}{},
		}
		if meta.Locale != "" {
			claims.PrivateClaims["locale"] = meta.Locale
		}
		if meta.HostedDomain != "" {
			claims.PrivateClaims["domain"] = meta.HostedDomain
		}
		if room := r.FormValue("room"); room != "" {
			claims.PrivateClaims["room"] = room
		}

		tok, err := signJWT(ctx, claims, key, permCheck)
		switch err {
		case errPermissionDenied:
			tr.LazyPrintf("Error signing token: %s", err.Error())
			writeError(ctx, w, "Permission check failed", http.StatusUnauthorized)
			return
		default:
			if err != nil {
				writeError(ctx, w, "Error encoding JWS", http.StatusInternalServerError)
				return
			}
		}
		fmt.Fprintf(w, tok)
	}
}
