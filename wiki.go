package main

import (
    "html/template" // to keep html in separate file
    "io/ioutil"
    "log"
    "net/http"
)

/* Page Data Structure
  Two fields, Title and Body
  []byte means a "bite slice"
    - type expected by the io libraries we will use
*/
type Page struct {
  Title string
  Body []byte
}

/* Save method for a Page
  - "This is a method named save that takes as its receiver p,
  a pointer to Page. It takes no parameters and returns a value of type error"
  - Will save the Page's Body to a text file using Title as the file name
  - error is the return type of WriteFile (a standard library function that writes
    a byte slice to a file)
  - If successful, Page.save() will return nil
  - 0600 is passed to Writefile to indicate the file should be created with r/w permissions for the current user
*/
func (p *Page) save() error{
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600)
}


/* Load a Page
    - Constructs a filename from title parameter
    - Reads the file's contents into variable body
    - Returns a pointer to Page literal constructed and an error (nil for no error)
    - ioutil.ReadFile() returns []byte and error
*/
func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil{
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

/* viewHandler that allows users to view a wiki Page
  - Will handle URLS prefixed with /view/
  - First extracts page title from r.URL.PATH
  - Loads the page data, formats the page with a string of simple HTML
  - Writes it to w, the http.ResponseWriter
*/
func viewHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/view/"):]
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }
  renderTemplate(w, "view", p)
}

/* editHandler
  - template.ParseFiles will read the contents of edit.html and return
    a *template.Template
  - t.Execute executes a template, writing the generated HTML to the http.ResponseWriter
  - .Title and .Body dotted identifiers refer to p.Title and p.Body
  - Template directives are enclosed in double curly braces in html {{ .Title }}
  - printf "%s" .Body instruction in html is a function call that outputs
*/
func editHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/edit/"):]
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit", p)
}

/* Render Template
  - Handles errors
*/
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
  t, err := template.ParseFiles(tmpl + ".html")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  err = t.Execute(w,p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}


/* Save a page */
func saveHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/save/"):]
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err := p.save()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}


/* Main */
func main() {

  // Page Functions
  // p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
  // p1.save()
  // p2, _ := loadPage("TestPage")
  // fmt.Println(string(p2.Body))

  // Handler
  // localhost:8080/view/[filename]
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/save/", saveHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))

}
