// Package jwtsi contains HTTP handlers and utilities for authenticating against
// a range of OAuth2 providers and returning signed JWT assertions about the
// authenticated user.
package jwtsi // import "github.com/jitsi/jwtsi"

// BUG: Jwtsi does not support TLS. To access the service with TLS (which you
//      really should be doing), use a reverse proxy such as Nginx.
