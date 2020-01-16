package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	var persons []entity.Person

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

// GetPerson function return an entity in the collection "person" from an ID
func (mongoClient *MongoClient) GetPerson(ID string) (entity.Person, error) {
	person := entity.Person{}
	objID, _ := primitive.ObjectIDFromHex(ID)
	err := mongoClient.personCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&person)
	return person, err
}

// InsertPerson function insert an entity in the collection "person"
func (mongoClient *MongoClient) InsertPerson(person entity.Person) {
	mongoClient.personCollection.InsertOne(context.Background(), person)
}

// DeletePerson function insert an entity in the collection "person" from an ID
func (mongoClient *MongoClient) DeletePerson(ID string) {
	objID, _ := primitive.ObjectIDFromHex(ID)
	mongoClient.personCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
}

// UpdatePerson function update an entity in the collection "person"
func (mongoClient *MongoClient) UpdatePerson(person entity.Person) error {
	pByte, err := bson.Marshal(person)
	if err != nil {
		return err
	}

	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		return err
	}

	mongoClient.personCollection.UpdateOne(context.Background(), bson.M{"_id": person.ID}, bson.D{{Key: "$set", Value: update}})
	return nil
}
