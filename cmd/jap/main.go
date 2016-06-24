package main

// BUG(ssw): JAP does not support TLS. To access the service with TLS (which you
//           really should be doing), use a reverse proxy such as Nginx.

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jitsi/jap"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"golang.org/x/text/language"
)

var (
	addr, pubDir, tmplDir, keyPath     string
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
	flag.StringVar(&keyPath, "key", os.Getenv("JAP_PRIVATE_KEY"), "An RSA private key in PEM format to use for signing tokens. Defaults to $JAP_PRIVATE_KEY.")
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
		log.Fatal(err)
	case len(files) < 1:
		log.Fatalf("No templates found in %s", tmplDir)
	}
	tmpl = template.Must(template.New("jap").ParseFiles(files...))
}

func loadRSAKeyFromPEM(keyPath string) (*rsa.PrivateKey, error) {
	pembytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	var blk *pem.Block
	for {
		blk, pembytes = pem.Decode(pembytes)
		if blk.Type == "RSA PRIVATE KEY" {
			return x509.ParsePKCS1PrivateKey(blk.Bytes)
		}
	}
	return nil, fmt.Errorf("No RSA private key found in pem file %s", keyPath)
}

func main() {
	if keyPath == "" {
		log.Fatalf("No private key specified. Try: %s -help", os.Args[0])
	}
	key, err := loadRSAKeyFromPEM(keyPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server on %s…\n", addr)

	http.HandleFunc("/googlelogin", jap.GoogleLogin(
		jap.NewCIDContext(context.Background(), googleClientID), key))
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

		tr := trace.New("jap.login", r.URL.Path)
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
