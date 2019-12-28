package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/alixleger/cours1/entity"
	"github.com/alixleger/cours1/mongodb"
)

const databaseName = "test"
const serverURI = "mongodb://localhost:27017"

func main() {
	dbClient := mongodb.New(serverURI, databaseName)
	fmt.Println("Connected to MongoDB!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			tmpl := template.Must(template.ParseFiles(filepath.Join("views", "person.html")))
			person := entity.Person{
				Firstname: r.FormValue("firstname"),
				Lastname:  r.FormValue("lastname"),
			}
			tmpl.Execute(w, person)
			dbClient.InsertPerson(person)
			return
		}

		tmpl := template.Must(template.ParseFiles(filepath.Join("views", "form.html")))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		persons, err := dbClient.GetPersons()

		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to get all persons : %v", err))
		}

		tmpl := template.Must(template.ParseFiles(filepath.Join("views", "list.html")))
		tmpl.Execute(w, struct{ Persons []entity.Person }{persons})
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
