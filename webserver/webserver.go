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
		log.Printf("Registring path %s", path)
		registerNewView(w, r, path)
	}
}

func registerNewView(w http.ResponseWriter, r *http.Request, path string) {
	registerLinkTemplate := template.Must(template.ParseFiles("views/register_link.html"))
	registerLinkTemplate.Execute(w, path)
}
