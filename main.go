package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/alixleger/cours1/entity"
	"github.com/alixleger/cours1/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const databaseName = "test"
const serverURI = "mongodb://localhost:27017"

func main() {
	formTmlp := template.Must(template.ParseFiles(filepath.Join("views", "form.html")))
	listTmpl := template.Must(template.ParseFiles(filepath.Join("views", "list.html")))

	dbClient := mongodb.New(serverURI, databaseName)
	fmt.Println("Connected to MongoDB!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			person := entity.Person{
				Firstname: r.FormValue("firstname"),
				Lastname:  r.FormValue("lastname"),
			}
			dbClient.InsertPerson(person)
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}

		persons, err := dbClient.GetPersons()
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to get all persons : %v", err))
		}
		listTmpl.Execute(w, struct{ Persons []entity.Person }{persons})
	})

	http.HandleFunc("/person/add", func(w http.ResponseWriter, r *http.Request) {
		params := struct {
			Person bool
			Action string
		}{
			false,
			"/",
		}
		formTmlp.Execute(w, params)
	})

	http.HandleFunc("/person/delete/", func(w http.ResponseWriter, r *http.Request) {
		dbClient.DeletePerson(strings.Replace(r.URL.Path, "/person/delete/", "", 1))
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})

	http.HandleFunc("/person/edit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			objID, _ := primitive.ObjectIDFromHex(r.FormValue("id"))
			person := entity.Person{
				ID:        objID,
				Firstname: r.FormValue("firstname"),
				Lastname:  r.FormValue("lastname"),
			}
			dbClient.UpdatePerson(person)
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})

	http.HandleFunc("/person/", func(w http.ResponseWriter, r *http.Request) {
		person, err := dbClient.GetPerson(strings.Replace(r.URL.Path, "/person/", "", 1))
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to get person : %v", err))
		}

		params := struct {
			Person entity.Person
			Action string
		}{
			person,
			"/person/edit",
		}
		formTmlp.Execute(w, params)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
