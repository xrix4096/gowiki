package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

//
// Globals
//
var gTemplates = template.Must(template.ParseFiles("edit.html", "view.html"))
var gValidPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

//
// Page type: Each page has a title and a body
//
type Page struct {
	Title string
	Body  []byte // byte 'slice' TODO learn about those
}

//
// Save method for writing a page to a file:
// Takes a pointer to a page 'p' as input, return type is 'error'
//
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

//
// Load method
// Takes a page name as input, returns a pointer to the loaded page and an error
//
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"             // Construct filename
	body, err := ioutil.ReadFile(filename) // Read data into 'body' var
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil // Construct and return Page
}

//
// viewHandler
// Load the page specified in the request then format and write the response
// to the supplied response writer
//
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

//
// editHandler
// Load or create the page specified in the request into an HTML edit form
//
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title} // Create page if not existing
	}
	renderTemplate(w, "edit", p)
}

//
// saveHandler
// Write the body of the page specified in the request to a file
//
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

//
// renderTemplate
// Load the specified HTML template file and subsitute the variables in the
// specified page
//
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := gTemplates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//
// makeHandler
// Create a closure for the handler functions that validates the URL using our
// regexp before calling the inner handler only if it's valid
//
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := gValidPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2]) // Page title is 2nd regexp match
	}
}



//
// main
//
func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}
