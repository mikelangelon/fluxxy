package section

import (
	"bol/contract/db"
	"time"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"bol/contract/types"
)

const SectionCollectionName = "section"

var (
	// SectionDB represents the MongoDB Booking collection.
	SectionDB func() *db.Collection
)

type Section struct {
	ID          int64      `json:"id" bson:"_id"`
	Version     int64      `json:"version" bson:"version"`
	Description string     `json:"description" bson:"description"`
	Title       string     `json:"title" bson:"title"`
	Text        string     `json:"text" bson:"text"`
	Actor       string     `json:"actor" bson:"actor"`
	History     []*Section `bson:"history" json:"-"`
}
type SectionVersion types.Version

func Insert(section Section) error {
	c := SectionDB()
	defer c.Close()

	section.ID = time.Now().UnixNano()
	section.Version = section.ID

	if err := c.Insert(&section); err != nil {
		fmt.Errorf("Error inserting contract %s", err)
		return err
	}
	return nil
}


func Update(section Section) error {
	c := SectionDB()
	defer c.Close()

	current, err := Lookup(section.ID)
	if err != nil {
		fmt.Errorf("Error doing lookups %s", err)
		return fmt.Errorf("Error doing lookups! : %v", err)
	}

	section.History = append(current.History, current)
	current.History = nil

	query := bson.M{
		"_id":     current.ID,
		"version": current.Version,
	}

	section.Version = time.Now().UnixNano()

	switch err := c.Update(query, &section); {
	case err == mgo.ErrNotFound:
		fmt.Errorf("Not found %s", err)
		return fmt.Errorf("not found : %v", err)
	case err != nil:
		fmt.Errorf("Error updating %s", err)
		return fmt.Errorf("Error updating:  %s", err)
	}

	return nil
}

func Lookup(id int64) (*Section, error) {
	current := new(Section)

	c := SectionDB()
	defer c.Close()

	switch err := c.FindId(id).One(current); {
	case err == mgo.ErrNotFound:
		return nil, fmt.Errorf("not found : %v", err)
	case err != nil:
		return nil, fmt.Errorf("Error retrieving : %v", err)
	}
	return current, nil
}

func LookupByVersion(id int64, version int64) (*Section, error) {
	current := new(Section)

	c := SectionDB()
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