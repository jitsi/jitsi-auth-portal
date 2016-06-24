# Jitsi Authentication Portal

[![GoDoc](https://godoc.org/github.com/jitsi/jap?status.svg)](https://godoc.org/github.com/jitsi/jap)
[![License](https://img.shields.io/badge/license-FreeBSD-blue.svg)](https://opensource.org/licenses/BSD-2-Clause)

The Jitsi Authentication Portal is an OAuth2 frontend that allows authentication
with a number of third party providers and generates short-lived, signed, JWT
(jot) tokens to assert the users identity to Jitsi Meet.

The package contains a number of handlers which can be used to build your own
compatible login service. There is also an example service in the `cmd/jap`
directory which provides a simple frontend.

To get started, install the `jap` command and run it:

```go
go get github.com/jitsi/jap
go install github.com/jitsi/jap/cmd/jap
jap -help
```

## License

The package may be used under the terms of the BSD 2-Clause License a copy of
which may be found in the file [LICENSE.md][LICENSE].

[LICENSE]: ./LICENSE.md
