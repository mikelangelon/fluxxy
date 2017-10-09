package contract

import (
	"bol/contract/db"
	"time"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"bol/contract/contract_template"
)

const ContractCollectionName = "contract"

var (
	// ContractDB represents the MongoDB Booking collection.
	ContractDB func() *db.Collection
)

type Contract struct {
	ID              int64                              `json:"id" bson:"_id"`
	Version         int64                              `json:"version" bson:"version"`
	Title           string                             `json:"title" bson:"title"`
	Description     string                             `json:"description" bson:"description"`
	Actor           string                             `json:"actor" bson:"actor"`
	TemplateVersion *contract_template.TemplateVersion `json:"templateVersion" bson:"templateVersion"`
	Params          []*Param                           `json:"params" bson:"params"`
	ParamsMap       map[string]string                  `json:"paramMap" bson:"paramMap"`
	History         []*Contract                        `bson:"history" json:"-"`
}

type Param struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

func (c Contract) GetBSON() (interface{}, error) {
	return struct {
		ID              int64                              `json:"id" bson:"_id"`
		Version         int64                              `json:"version" bson:"version"`
		Title           string                             `json:"title" bson:"title"`
		Description     string                             `json:"description" bson:"description"`
		Actor           string                             `json:"actor" bson:"actor"`
		TemplateVersion *contract_template.TemplateVersion `json:"templateVersion" bson:"templateVersion"`
		Params          []*Param                           `json:"params" bson:"params"`
		ParamsMap       []*Param                           `json:"paramMap" bson:"paramMap"`
		History         []*Contract                        `bson:"history" json:"-"`
	}{
		ID:              c.ID,
		Version:         c.Version,
		Title:           c.Title,
		Description:     c.Description,
		Actor:           c.Actor,
		TemplateVersion: c.TemplateVersion,
		Params:          c.Params,
		ParamsMap:       MapToParams(c.ParamsMap),
		History:         c.History,
	}, nil
}

func MapToParams(paramsMap map[string]string) []*Param {
	params := []*Param{}
	for value, key := range paramsMap {
		params = append(params, &Param{key, value})
	}
	return params
}

func (c *Contract) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		ID              int64                              `json:"id" bson:"_id"`
		Title           string                             `json:"title" bson:"title"`
		Description     string                             `json:"description" bson:"description"`
		Actor           string                             `json:"actor" bson:"actor"`
		TemplateVersion *contract_template.TemplateVersion `json:"templateVersion" bson:"templateVersion"`
		Params          []*Param                           `json:"params" bson:"params"`
		ParamsMap       []*Param                           `json:"paramMap" bson:"paramMap"`
	})
	bsonErr := raw.Unmarshal(decoded)

	if bsonErr == nil {
		c.ID = decoded.ID
		c.Title = decoded.Title
		c.Description = decoded.Description
		c.Actor = decoded.Actor
		c.TemplateVersion = decoded.TemplateVersion
		c.Params = decoded.Params
		paramsMap := map[string]string{}
		for _, param := range decoded.ParamsMap {
			paramsMap[param.Key] = param.Value
		}
		c.ParamsMap = paramsMap
		return nil
	} else {
		return bsonErr
	}
	return nil
}

func Insert(contract Contract) error {
	c := ContractDB()
	defer c.Close()

	contract.ID = time.Now().UnixNano()

	if err := c.Insert(&contract); err != nil {
		fmt.Errorf("Error inserting contract %s", err)
		return err
	}
	return nil
}

func Update(contract Contract) error {
	c := ContractDB()
	defer c.Close()

	current, err := Lookup(contract.ID)
	if err != nil {
		return fmt.Errorf("Error doing lookups! : %v", err)
	}

	contract.History = append(current.History, current)
	current.History = nil

	query := bson.M{
		"_id": current.ID,
	}

	contract.Version = time.Now().UnixNano()

	switch err := c.Update(query, &contract); {
	case err == mgo.ErrNotFound:
		return fmt.Errorf("not found : %v", err)
	case err != nil:
		return fmt.Errorf("Error updating:  %s", err)
	}

	return nil
}

func Lookup(id int64) (*Contract, error) {
	current := new(Contract)

	c := ContractDB()
	defer c.Close()

	switch err := c.FindId(id).One(current); {
	case err == mgo.ErrNotFound:
		return nil, fmt.Errorf("not found : %v", err)
	case err != nil:
		return nil, fmt.Errorf("Error retrieving : %v", err)
	}
	return current, nil
}