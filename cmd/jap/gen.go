// +build ignore

// This file is used to generate the help text and documentation for the jap
// command.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const help = `The jap command launches an OAuth2 server that generates a JSON
Web Signature (JWS) to prove the users identity to other Jitsi services.

Environment

To use the supported providers, a few environment variables must be set:

  GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET: Needed to support login with Google.

Signals

On POSIX based systems templates can be reloaded in a running process by sending
the process a SIGHUP. For more information on POSIX signals, see the signal(7)
man page.`

func main() {
	fh, err := os.Create("doc.go")
	if err != nil {
		panic(err.Error())
	}
	defer fh.Close()

	fmt.Fprintln(fh, "// This file was generate by go generate; DO NOT EDIT\n")

	r := strings.NewReader(help)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			fmt.Fprintln(fh, "//")
			continue
		}
		fmt.Fprintf(fh, "// %s\n", text)
	}
	fmt.Fprintln(fh, "package main\n")

	fmt.Fprintln(fh, "//go:generate go run gen.go\n")

	fmt.Fprint(fh, "const help = `")
	fmt.Fprint(fh, help)
	fmt.Fprintln(fh, "`")
}
