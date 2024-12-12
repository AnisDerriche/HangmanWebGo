package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var Variable string

func Index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/Index.html"))
	t.Execute(w, nil)
}

func Jeux(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/jeux.html"))
	t.Execute(w, nil)
}

func Hang(w http.ResponseWriter, r *http.Request) {
	input := r.PostFormValue("input")
	strhtml := fmt.Sprintf(" %s ", input)
	tmpl, _ := template.New("t").Parse(strhtml)
	tmpl.Execute(w, nil)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", Index)
	http.HandleFunc("/jeux", Jeux)
	http.HandleFunc("/lettre/", Hang)
	fmt.Println("http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
