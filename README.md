# Jitsi JWT Service

[![GoDoc](https://godoc.org/github.com/jitsi/jwtsi?status.svg)](https://godoc.org/github.com/jitsi/jwtsi)
[![License](https://img.shields.io/badge/license-FreeBSD-blue.svg)](https://opensource.org/licenses/BSD-2-Clause)

The Jitsi JWT Service is an OAuth2 frontend that allows authentication with a
number of providers and generates a short-lived, signed, JWT (jot) token to
assert the users identity to Jitsi Meet.

The package contains a number of handlers which can be used to build your own
compatible login service. There is also an example service in the `cmd/jwtsi`
directory which provides a nice frontend (which can be loaded in an iframe) that
supports several providers, and the various required login endpoints.

To get started, install the `jwtsi` command and run it:

```go
go get github.com/jitsi/jwtsi
go install github.com/jitsi/jwtsi/cmd/jwtsi
jwtsi -help
```

## License

The package may be used under the terms of the BSD 2-Clause License a copy of
which may be found in the file [LICENSE.md][LICENSE].

[LICENSE]: ./LICENSE.md
