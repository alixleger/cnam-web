package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Person is a type of person
type Person struct {
	Firstname string
	Lastname  string
}

// Persons list
type Persons struct {
	List []Person
}

var client *mongo.Client

const databseNAME = "test"
const collectionName = "person"
const serverURI = "mongodb://localhost:27017"

func main() {
	client, err := createMongoDBClient(serverURI)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to mongodb server : %v", err))
	}
	fmt.Println("Connected to MongoDB!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			tmpl := template.Must(template.ParseFiles(filepath.Join("views", "person.html")))
			person := Person{
				Firstname: r.FormValue("firstname"),
				Lastname:  r.FormValue("lastname"),
			}
			tmpl.Execute(w, person)
			collection := client.Database(databseNAME).Collection(collectionName)
			collection.InsertOne(context.Background(), person)
			return
		}

		tmpl := template.Must(template.ParseFiles(filepath.Join("views", "form.html")))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		personsList, err := getAllPersons(client)

		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to get all persons : %v", err))
		}

		tmpl := template.Must(template.ParseFiles(filepath.Join("views", "list.html")))
		tmpl.Execute(w, Persons{List: personsList})
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createMongoDBClient(serverURI string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(serverURI))

	if err != nil || client.Ping(context.Background(), readpref.Primary()) != nil {
		return nil, err
	}

	return client, err
}

func getAllPersons(client *mongo.Client) ([]Person, error) {
	databseNAME := "test"
	collectionName := "person"

	collection := client.Database(databseNAME).Collection(collectionName)
	cursor, err := collection.Find(context.Background(), bson.D{{}})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	persons := make([]Person, 0)
	for cursor.Next(context.Background()) {
		person := Person{}
		err := cursor.Decode(&person)
		if err != nil {
			log.Fatal(err)
		}
		if person.Firstname == "" || person.Lastname == "" {
			continue
		}
		persons = append(persons, person)
	}

	return persons, nil
}
