package webserver

import (
	"fmt"
	"golinks/database"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func Run(port int) {

	http.HandleFunc("/", requestHandler)
	http.HandleFunc("/create/", registerNewView)
	http.HandleFunc("/favicon.ico", http.NotFound)

	log.Printf("Server running on port :%d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	if url, ok := database.Lookup(path); ok {
		log.Printf("Redirecting to %s", url)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		log.Printf("Path %s not found; Asking to register.", path)
		registerLinkTemplate := template.Must(template.ParseFiles("views/not_found.html"))
		if err := registerLinkTemplate.Execute(w, path); err != nil {
			log.Print(err)
		}
	}
}

type context struct {
	Path      string
	Url       string
	Submitted bool
}

func registerNewView(w http.ResponseWriter, r *http.Request) {
	registerLinkTemplate := template.Must(template.ParseFiles("views/register_link.html"))
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
		}
		path := r.Form.Get("path")
		url := r.Form.Get("url")
		if err := database.RegisterUrl(path, url); err != nil {
			log.Fatal(err)
		}
		cont := context{Path: path, Url: url, Submitted: true}
		if err := registerLinkTemplate.Execute(w, cont); err != nil {
			log.Print(err)
		}

	} else {
		path := strings.TrimLeft(r.URL.RawQuery, "path=")
		cont := context{Path: path, Url: "", Submitted: false}
		if err := registerLinkTemplate.Execute(w, cont); err != nil {
			log.Print(err)
		}
	}
}