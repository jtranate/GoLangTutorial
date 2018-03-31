package main

import (
    "fmt"
    "io/ioutil"
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


/* Main */
func main() {
  p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
  p1.save()
  p2, _ := loadPage("TestPage")
  fmt.Println(string(p2.Body))
}
