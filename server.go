package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/alixleger/cours1/api"
	"github.com/alixleger/cours1/handler"
	"github.com/alixleger/cours1/mongodb"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/rs/cors"
)

const databaseName = "test"
const serverURI = "mongodb://localhost:27017"

func main() {
	dbClient := mongodb.New(serverURI, databaseName)
	fmt.Println("Connected to MongoDB!")

	router := mux.NewRouter()

	// RPC API Server
	go func() {
		persons := api.Persons{DBClient: &dbClient}

		s := rpc.NewServer()
		s.RegisterCodec(json.NewCodec(), "application/json")
		s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

		c := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"POST", "OPTIONS"},
		})

		s.RegisterService(&persons, "")

		router.Handle("/rpc", s)
		log.Fatal(http.ListenAndServe(":1234", c.Handler(router)))
	}()

	// Fullstack App Server and API client

	formTmpl := template.Must(template.ParseFiles(filepath.Join("views", "form.tmpl"), filepath.Join("views", "layout.tmpl")))
	listTmpl := template.Must(template.ParseFiles(filepath.Join("views", "list.tmpl"), filepath.Join("views", "layout.tmpl")))
	apiClientTmpl := template.Must(template.ParseFiles(filepath.Join("views", "api_client.tmpl"), filepath.Join("views", "layout.tmpl")))

	var validInput = regexp.MustCompile(`^[a-zA-ZàâæçéèêëîïôœùûüÿÀÂÆÇnÉÈÊËÎÏÔŒÙÛÜŸ\-\s]+$`)

	routerHandlers := handler.Handlers{
		FormTemplate:      formTmpl,
		ListTemplate:      listTmpl,
		APIClientTemplate: apiClientTmpl,
		ValidInput:        validInput,
		DBClient:          &dbClient,
	}

	router.HandleFunc("/", routerHandlers.ListHandler)

	router.HandleFunc("/person/delete/{id}", routerHandlers.DeleteHandler)

	router.HandleFunc("/person/edit", routerHandlers.UpdateHandler)

	router.HandleFunc("/person/{id}", routerHandlers.GetPersonHandler)

	router.HandleFunc("/client", routerHandlers.APIClientHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
