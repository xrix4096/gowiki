package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

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
	filename := title + ".txt"  // Construct filename
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
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
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
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
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
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

	
//
// renderTemplate
// Load the specified HTML template file and subsitute the variables in the
// specified page
//
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}


//
// main
//
func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}

