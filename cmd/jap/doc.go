// This file was generate by go generate; DO NOT EDIT

// The jap command launches an OAuth2 server that generates a JSON
// Web Signature (JWS) to prove the users identity to other Jitsi services.
//
// Environment
//
// To use the supported providers, a few environment variables must be set:
//
//   GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET: Needed to support login with Google.
//
// The private key used to sign JWTs can also be loaded directly from an
// environment variable:
//
//   JAP_PRIVATE_KEY
//
// Signals
//
// On POSIX based systems templates can be reloaded in a running process by sending
// the process a SIGHUP. For more information on POSIX signals, see the signal(7)
// man page.
//
// For more information try:
//
//    jap -help
package main

//go:generate go run gen.go

const help = `The jap command launches an OAuth2 server that generates a JSON
Web Signature (JWS) to prove the users identity to other Jitsi services.

Environment

To use the supported providers, a few environment variables must be set:

  GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET: Needed to support login with Google.

The private key used to sign JWTs can also be loaded directly from an
environment variable:

  JAP_PRIVATE_KEY

Signals

On POSIX based systems templates can be reloaded in a running process by sending
the process a SIGHUP. For more information on POSIX signals, see the signal(7)
man page.`
