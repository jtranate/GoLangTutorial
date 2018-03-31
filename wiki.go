package main

import (
    "html/template" // to keep html in separate file
    "io/ioutil"
    "log"
    "net/http"
    "regexp"
    "errors" // To create new errors
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
  filename := "data/" + p.Title + ".txt"
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
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
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
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit", p)
}

/* Save a page */
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


func rootHandler(w http.ResponseWriter, r *http.Request){
  http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}


/* Gloabl templates variable
  - Call ParseFiles once at program initialization, parsing all templates into a single *Template
  - Then can use ExecuteTemplate method to render a specific template
  - Must is a convenience wrapper that panics when passed a non-nil error value, otherwise returns the *Template unaltered
    - Panic is appropriate here if template can't be loaded, so it will exit the program
  - ParseFiles can take any number of strings
*/
var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))



/* Render Template
  - Handles errors
*/
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
  err := templates.ExecuteTemplate(w, tmpl + ".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

/* Validate title with regular expression
  - regexp.MustCompile will parse and compile the regex and return a
    regexp.Regexp.Mustcompile is distinct from Compile in that it will panic if expression
    compilation fails, while Compile returns an error as a second parameter.
*/
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// Don't need because we added makeHandler
/* Function to validate path and extract the page title */
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w,r)
    return "", errors.New("Invalid Page Title")
  }
  return m[2], nil // The title is the second subexpression
 }

/* Wrapper function that takes a function and returns a function of type http.HandlerFunc
  - returned function called a closure bc it encloses values defined outside of it
  - the variable fn is enclosed by the closure
  - fn will be one of our save, edit, or view handlers
  - Closure extracts the title from the request path and validates it with the TitleValidator regexp
    - if invalid, error will be written to the ResponseWriter using the http.NotFound function
    - if valid, enclosed handler function fn will be called with the ResponseWriter, Request, and title as arguments

*/
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    // Here we will extract the page title from the Request and call the provided handler 'fn'
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
      http.NotFound(w, r)
      return
    }
    fn(w, r, m[2])
  }
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
  http.HandleFunc("/", rootHandler)
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
  log.Fatal(http.ListenAndServe(":8080", nil))

}
