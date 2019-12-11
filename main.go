package main

import (
	"html/template"
	"log"
	"net/http"
)

// Person is a type of person
type Person struct {
	Firstname string
	Lastname  string
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		person := Person{
			Firstname: r.FormValue("firstname"),
			Lastname:  r.FormValue("lastname"),
		}
		t, _ := template.ParseFiles("person.html")
		t.Execute(w, person)
		return
	}

	t, _ := template.ParseFiles("form.html")
	t.Execute(w, nil)
}
