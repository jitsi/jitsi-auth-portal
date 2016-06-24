// The jap command launches an OAuth2 server that generates a JSON
// Web Signature (JWS) to prove the users identity to other Jitsi services.
//
// To use the supported providers, a few environment variables must be set:
//
// ENV:
//
//   GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET: Needed to support login with Google.
package main

//go:generate go run gen.go

const help = `The jap command launches an OAuth2 server that generates a JSON
Web Signature (JWS) to prove the users identity to other Jitsi services.

To use the supported providers, a few environment variables must be set:

ENV:

  GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET: Needed to support login with Google.`
