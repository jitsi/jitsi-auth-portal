// The jwtsi command launches an OAuth2 server that generates a JSON Web
// Signature (JWS) to prove the users identity to other Jitsi services.
//
// To get started run jwtsi -help
//
// Jwtsi does not have an option to listen for HTTPS connections. To use TLS,
// put Jwtsi behind a reverse proxy such as nginx.
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jitsi/jwtsi"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"golang.org/x/text/language"
)

const help = `The jwtsi command launches an OAuth2 server that generates a JSON
Web Signature (JWS) to prove the users identity to other Jitsi services.

To use the supported providers, a few environment variables must be set:

ENV:

  GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET: Needed to support login with Google.`

var (
	addr, pubDir, tmplDir              string
	googleClientSecret, googleClientID string
	redirectURL                        string

	tmpl    *template.Template
	devMode bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nUsage of %s:\n", help, os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&addr, "http", ":http-alt", "The address to listen on.")
	flag.StringVar(&pubDir, "public", "public/", "A directory containing static files to serve.")
	flag.StringVar(&tmplDir, "templates", "templates/", "A directory containing templates to render.")
	flag.StringVar(&redirectURL, "redirect", "https://meet.jit.si", "The URL to redirect back too after performing OAuth.")
	flag.BoolVar(&devMode, "dev", false, "Run in dev mode (reload templates on page refresh).")
	flag.Parse()

	googleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")

	loadTemplates()
}

// Load all templates found in the tmplDir directory; if any of them contain
// errors, panic.
func loadTemplates() {
	files, err := filepath.Glob(filepath.Join(tmplDir, "*.tmpl"))
	switch {
	case err != nil:
		panic(err)
	case len(files) < 1:
		panic("No templates found in " + tmplDir)
	}
	tmpl = template.Must(template.New("jwtsi").ParseFiles(files...))
}

func main() {
	log.Printf("Starting server on %s…\n", addr)

	http.HandleFunc("/googlelogin", jwtsi.GoogleLogin(
		jwtsi.NewCIDContext(context.Background(), googleClientID)))
	http.HandleFunc("/login", loginHandler(context.Background()))
	if pubDir != "" {
		http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(pubDir))))
	}
	log.Fatal(http.ListenAndServe(addr, nil))
}

func loginHandler(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if devMode {
			loadTemplates()
		}

		tr := trace.New("jwtsi.login", r.URL.Path)
		defer tr.Finish()

		tr.LazyPrintf("Executing login.tmpl…")
		err := tmpl.ExecuteTemplate(w, "login.tmpl", Login{
			Lang:           language.English,
			GoogleClientID: googleClientID,
		})
		if err != nil {
			tr.LazyPrintf("Error exeuting login.tmpl:", err.Error())
			tr.SetError()
			return
		}
		tr.LazyPrintf("Done executing login.tmpl…")
	}
}

// Login represents all the information we need to show the login window.
type Login struct {
	Lang           language.Tag
	GoogleClientID string
}
