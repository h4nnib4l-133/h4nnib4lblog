// this is a simple webserver reading from text files so basically i could edit them and jsut reload the server to refresh with new content this is probably the worst. but i mean who gives a fuck. Will improve with time. 

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
	"regexp"
	"errors"
)

type Page struct {
	Title string
	Body  template.HTML
}

var templates = template.Must(template.ParseFiles("blog.html", "source.html"))

func loadPage(title string) (*Page, error) {
	filename := title+"_2.html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var s template.HTML
	s = template.HTML(body) 
	return &Page{Title: title, Body: s}, nil
}

var validPath = regexp.MustCompile("^/(blog|source)") 

func getTitle(w http.ResponseWriter, r *http.Request)(string, error){
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil{
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[1], nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
}

func blogHandler(w http.ResponseWriter, r *http.Request){
	title, err := getTitle(w, r)
	if err != nil{
		return
	}
	p, err := loadPage(title)
	renderTemplate(w, "blog", p)
}

func sourceHandler(w http.ResponseWriter, r *http.Request){
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	renderTemplate(w, "source", p)
}


func main() {
	fmt.Println("Listening on server 8080")
	http.HandleFunc("/blog/", blogHandler)
	http.HandleFunc("/source/", sourceHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

