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
