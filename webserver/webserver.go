package webserver

import (
	"golinks/database"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func Run(port string) {

	http.HandleFunc("/", requestHandler)
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/create/", registerHandler)

	// Tutorial
	http.HandleFunc("/intro", renderTemplateHandler("views/intro.html"))
	http.HandleFunc("/quite-hard-to-memorize-pattern-here-not-gonna-lie", renderTemplateHandler("views/quite-hard-to-memoroze-pattern-here-not-gonna-lie.html"))
	http.HandleFunc("/register", renderTemplateHandler("views/new_links.html"))
	http.HandleFunc("/finish", renderTemplateHandler("views/finish.html"))

	log.Printf("Server running on port :%s", port)
	http.ListenAndServe(":"+port, nil)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	if path == "" {
		http.Redirect(w, r, "/intro", http.StatusFound)
	} else if url, ok := database.Lookup(path); ok {
		log.Printf("Redirecting to %s", url)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		log.Printf("Path %s not found; Asking to register.", path)
		registerLinkTemplate := loadTemplate("views/not_found.html", defaultWrappers)
		if err := registerLinkTemplate.Execute(w, path); err != nil {
			log.Print(err)
		}
	}
}

func renderTemplateHandler(templatePath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		introTemplate := loadTemplate(templatePath, defaultWrappers)
		introTemplate.Execute(w, "intro")
	}
}

func introHandler(w http.ResponseWriter, r *http.Request) {
	introTemplate := loadTemplate("views/intro.html", defaultWrappers)
	introTemplate.Execute(w, "intro")
}

type context struct {
	Path      string
	Url       string
	Submitted bool
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	registerLinkTemplate := loadTemplate("views/register_link.html", defaultWrappers)
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

// Structure for storing default wrappers.
type wrappers struct {
	header string
	footer string
	navbar string
}

// Default set of header and footer wrappers.
var defaultWrappers = wrappers{
	header: "views/header.html",
	footer: "views/footer.html",
	navbar: "views/navbar.html",
}

// loadTemaplate parses provided template from file, appends footer and header
// and makes sure that template loaded correctly or panics if not.
func loadTemplate(templatePath string, w wrappers) *template.Template {
	templateWithHeaders := withWrappers(templatePath, w)
	return template.Must(template.ParseFiles(templateWithHeaders...))
}

//withWrappers appends header and footer to template.
func withWrappers(templatePath string, w wrappers) (paths []string) {
	paths = append(paths, templatePath)
	paths = append(paths, w.header)
	paths = append(paths, w.footer)
	paths = append(paths, w.navbar)
	return paths
}
