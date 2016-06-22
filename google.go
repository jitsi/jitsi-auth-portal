package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/net/trace"
)

// Errors:
//  400 BadRequest          – If the id_token form param is missing.
//  408 RequestTimeout      – If the contexts deadline was exceeded.
//  500 InternalServerError – If the upstream returns a response we don't understand.
//  502 BadGateway          – If an upstream service fails to respond for another reason.
func googleLoginHandler(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tr := trace.New("jwtsi.tokenlogin", r.URL.Path)
		defer tr.Finish()

		idtoken := r.FormValue("id_token")
		if idtoken == "" {
			tr.LazyPrintf("Requires id_token login param")
			tr.SetError()
			http.Error(w, "id_token login param missing", http.StatusBadRequest)
			return
		}

		inforesp, err := ctxhttp.PostForm(ctx, http.DefaultClient,
			"https://www.googleapis.com/oauth2/v3/tokeninfo",
			url.Values{
				"id_token": []string{idtoken},
			},
		)
		if err != nil {
			tr.LazyPrintf("Error validating token:", err.Error())
			tr.SetError()
			switch err {
			case context.DeadlineExceeded:
				http.Error(w, "The deadline was exceeded", http.StatusRequestTimeout)
			default:
				http.Error(w, "Upstream request failed", http.StatusBadGateway)
			}
			return
		}
		if inforesp.StatusCode < 200 || inforesp.StatusCode > 299 {
			tr.LazyPrintf("Error validating token: Got unexpected response", inforesp.StatusCode)
			tr.SetError()
			http.Error(w, "Unexpected response from upstream", http.StatusInternalServerError)
			return
		}
		values := struct {
			Aud          string `json:"aud"`
			Email        string `json:"email"`
			HostedDomain string `json:"hd"`
			Locale       string `json:"locale"`
		}{}
		if err := json.NewDecoder(inforesp.Body).Decode(&values); err != nil {
			tr.LazyPrintf("Error decoding upstream response")
			tr.SetError()
			http.Error(w, "Error decoding upstream response", http.StatusInternalServerError)
			return
		}
		// TODO: More work.
		log.Println(values)
	}
}
