// Package jap (Jitsi Authentication Provider) contains HTTP handlers and
// utilities for authenticating against a range of OAuth2 providers and
// returning signed JWT assertions about the authenticated user.
package jap // import "github.com/jitsi/jap"

// BUG: Jap does not support TLS. To access the service with TLS (which you
//      really should be doing), use a reverse proxy such as Nginx.
