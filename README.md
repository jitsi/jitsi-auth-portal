# Jitsi Authentication Portal

[![GoDoc](https://godoc.org/github.com/jitsi/jap?status.svg)](https://godoc.org/github.com/jitsi/jap)
[![License](https://img.shields.io/badge/license-FreeBSD-blue.svg)](https://opensource.org/licenses/BSD-2-Clause)

The Jitsi JWT Service is an OAuth2 frontend that allows authentication with a
number of providers and generates a short-lived, signed, JWT (jot) token to
assert the users identity to Jitsi Meet.

To get started, install `jap` and run it:

```go
go get github.com/jitsi/jap
go install github.com/jitsi/jap
jap -help
```

## License

The package may be used under the terms of the BSD 2-Clause License a copy of
which may be found in the file [LICENSE.md][LICENSE].

[LICENSE]: ./LICENSE.md
