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

To use the supported providers, a few environment variables must be set:

ENV:

  GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET: Needed to support login with Google.`

func main() {
	fh, err := os.Create("doc.go")
	if err != nil {
		panic(err.Error())
	}
	defer fh.Close()

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
