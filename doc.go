// The jwtsi command launches an OAuth2 server that generates a JSON Web
// Signature (JWS) to prove the users identity to other Jitsi services.
//
// To get started run jwtsi -help
//
// Jwtsi does not have an option to listen for HTTPS connections. To use TLS,
// put Jwtsi behind a reverse proxy such as nginx.
package main // import "github.com/jitsi/jwtsi"

const help = `The jwtsi command launches an OAuth2 server that generates a JSON
Web Signature (JWS) to prove the users identity to other Jitsi services.

To use the supported providers, a few environment variables must be set:

ENV:

  GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET: Needed to support login with Google.`

// BUG: Jwtsi does not support TLS. To access the service with TLS (which you
//      really should be doing), use a reverse proxy such as Nginx.
