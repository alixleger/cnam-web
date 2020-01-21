package api

import (
	"net/http"

	"github.com/alixleger/cours1/entity"
	"github.com/alixleger/cours1/mongodb"
)

// Persons is the method type of the RPC API
type Persons struct {
	DBClient *mongodb.MongoClient
}

// GetPersons is the RPC operator to get all mongodb person entities
func (p *Persons) GetPersons(r *http.Request, args *string, result *[]entity.Person) error {
	persons, err := p.DBClient.GetPersons()
	if err == nil {
		*result = persons
	}

	return err
}

// GetPerson is the RPC operator to get a mongodb person entity from an ID
func (p *Persons) GetPerson(r *http.Request, args *string, result *entity.Person) error {
	person, err := p.DBClient.GetPerson(*args)
	if err == nil {
		*result = person
	}

	return err
}

// InsertPerson is the RPC operator to insert a mongodb person entity
func (p *Persons) InsertPerson(r *http.Request, args *entity.Person, result *string) error {
	p.DBClient.InsertPerson(*args)
	*result = "OK"
	return nil
}

// UpdatePerson is the RPC operator to update a mongodb person entity
func (p *Persons) UpdatePerson(r *http.Request, args *entity.Person, result *string) error {
	p.DBClient.UpdatePerson(*args)
	*result = "OK"
	return nil
}

// DeletePerson is the RPC operator to delete a mongodb person entity
func (p *Persons) DeletePerson(r *http.Request, args *string, result *string) error {
	p.DBClient.DeletePerson(*args)
	*result = "OK"
	return nil
}
