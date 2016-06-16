package main

import (
	"fmt"
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
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title} // Create page if not existing
	}
	fmt.Fprintf(w, "<h1>Editing %s</h1>" +
		"<form action=\"/save/%s\" method=\"POST\">" +
		"<textarea name=\"body\">%s</textarea><br>" +
		"<input type=\"submit\" value=\"Save\">" +
		"</form>",
		p.Title, p.Title, p.Body)
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

