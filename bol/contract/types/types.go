package types

type Version struct {
	ID          int64      `json:"id" bson:"id"`
	Version     int64      `json:"version" bson:"version"`
}
