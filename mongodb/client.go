package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/alixleger/cours1/entity"
)

const personCollectionName = "person"

// MongoClient is a MongoDB client interface
type MongoClient struct {
	personCollection *mongo.Collection
}

// New function is the MongoClient constructor
func New(serverURI string, databaseName string) MongoClient {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(serverURI))

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to mongodb server : %v", err))
	}

	if client.Ping(context.Background(), readpref.Primary()) != nil {
		log.Fatal("Failed to ping mongodb server")
	}

	return MongoClient{personCollection: client.Database(databaseName).Collection(personCollectionName)}
}

// GetPersons function return all entities of the collection "person"
func (mongoClient *MongoClient) GetPersons() ([]entity.Person, error) {
	cursor, err := mongoClient.personCollection.Find(context.Background(), bson.D{{}})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	persons := make([]entity.Person, 0)
	for cursor.Next(context.Background()) {
		person := entity.Person{}
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

// InsertPerson function insert an entity in the collection "person"
func (mongoClient *MongoClient) InsertPerson(person entity.Person) {
	mongoClient.personCollection.InsertOne(context.Background(), person)
}

/*

// DeletePerson function insert an entity in the collection "person"
func (mongoClient *MongoClient) DeletePerson(person entity.Person) {

}

// UpdatePerson function update an entity in the collection "person"
func (mongoClient *MongoClient) UpdatePerson(person entity.Person) {

}

*/
