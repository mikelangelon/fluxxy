package contract_template

import (
	"bol/contract/db"
	"bol/contract/section"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"bol/contract/types"
)

const TemplateCollectionName = "template"

var (
	// TemplateDB represents the MongoDB Booking collection.
	TemplateDB func() *db.Collection
)

type Template struct {
	ID          int64                     `json:"id" bson:"_id"`
	Version     int64                     `json:"version" bson:"version"`
	Description string                    `json:"description" bson:"description"`
	Title       string                    `json:"title" bson:"title"`
	Actor       string                    `json:"actor" bson:"actor"`
	Sections    []*section.SectionVersion `json:"sections" bson:"sections"`
	History     []*Template               `bson:"history" json:"-"`
}
type TemplateVersion types.Version

func EnsureTemplateIndice() error {
	{
		idx := mgo.Index{
			Name:   "example",
			Unique: false,
			Key:    []string{"id"},
		}

		c := TemplateDB()
		defer c.Close()

		if err := c.EnsureIndex(idx); err != nil {
			return fmt.Errorf("Error creating index: %s", err)
		}
	}
	return nil
}

func Insert(template Template) error {
	c := TemplateDB()
	defer c.Close()

	template.ID = time.Now().UnixNano()
	template.Version = template.ID

	if err := c.Insert(&template); err != nil {
		fmt.Errorf("Error inserting contract_template %s", err)
		return err
	}
	return nil
}

func Update(template Template) error {
	c := TemplateDB()
	defer c.Close()

	current, err := Lookup(template.ID)
	if err != nil {
		return fmt.Errorf("Error doing lookups! : %v", err)
	}

	template.History = append(current.History, current)
	current.History = nil

	query := bson.M{
		"_id":     current.ID,
		"version": current.Version,
	}

	template.Version = time.Now().UnixNano()

	switch err := c.Update(query, &template); {
	case err == mgo.ErrNotFound:
		return fmt.Errorf("not found : %v", err)
	case err != nil:
		return fmt.Errorf("Error updating:  %s", err)
	}

	return nil
}

func Lookup(id int64) (*Template, error) {
	current := new(Template)

	c := TemplateDB()
	defer c.Close()

	switch err := c.FindId(id).One(current); {
	case err == mgo.ErrNotFound:
		return nil, fmt.Errorf("not found : %v", err)
	case err != nil:
		return nil, fmt.Errorf("Error retrieving : %v", err)
	}
	return current, nil
}

func LookupByVersion(id int64, version int64) (*Template, error) {
	current := new(Template)

	c := TemplateDB()
	defer c.Close()

	switch err := c.FindId(id).One(current); {
	case err == mgo.ErrNotFound:
		return nil, fmt.Errorf("not found : %v", err)
	case err != nil:
		return nil, fmt.Errorf("Error retrieving : %v", err)
	}
	sections := current.History
	sections = append(sections, current)
	for _, section := range sections {
		if section.Version == version {
			return section, nil
		}
	}
	return nil, fmt.Errorf("Non existing version ")
}