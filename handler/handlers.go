package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/alixleger/cours1/entity"
	"github.com/alixleger/cours1/mongodb"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Handlers is an the handlers manager of the app
type Handlers struct {
	FormTemplate      *template.Template
	ListTemplate      *template.Template
	APIClientTemplate *template.Template
	ValidInput        *regexp.Regexp
	DBClient          *mongodb.MongoClient
}

// ListHandler insert person entity if POST or launch list template if GET
func (h *Handlers) ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		firstnameInput := strings.Trim(r.FormValue("firstname"), " ")
		lastnameInput := strings.Trim(r.FormValue("lastname"), " ")

		if h.ValidInput.MatchString(firstnameInput) && h.ValidInput.MatchString(lastnameInput) {
			person := entity.Person{
				Firstname: firstnameInput,
				Lastname:  lastnameInput,
			}
			h.DBClient.InsertPerson(person)
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	persons, err := h.DBClient.GetPersons()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}

	err = h.ListTemplate.ExecuteTemplate(w, "layout", struct{ Persons []entity.Person }{persons})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}
}

// DeleteHandler delete a person entity
func (h *Handlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	h.DBClient.DeletePerson(vars["id"])
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// UpdateHandler update a person entity
func (h *Handlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		firstnameInput := strings.Trim(r.FormValue("firstname"), " ")
		lastnameInput := strings.Trim(r.FormValue("lastname"), " ")
		objID, err := primitive.ObjectIDFromHex(strings.Trim(r.FormValue("id"), " "))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
		}

		if h.ValidInput.MatchString(firstnameInput) && h.ValidInput.MatchString(lastnameInput) && err == nil {
			person := entity.Person{
				ID:        objID,
				Firstname: firstnameInput,
				Lastname:  lastnameInput,
			}
			h.DBClient.UpdatePerson(person)
		}
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// GetPersonHandler get a person entity
func (h *Handlers) GetPersonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	person, err := h.DBClient.GetPerson(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}

	params := struct {
		Person entity.Person
		Action string
	}{
		person,
		"/person/edit",
	}
	err = h.FormTemplate.ExecuteTemplate(w, "layout", params)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}
}

// APIClientHandler launch api client app
func (h *Handlers) APIClientHandler(w http.ResponseWriter, r *http.Request) {
	err := h.APIClientTemplate.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}
}
