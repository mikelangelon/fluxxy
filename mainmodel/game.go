package mainmodel

type Game struct {
	ID int64 `json:"id"`
	Players []*Player `json:"players"`
	Rules []*Card `json:"rules"`
	Cards []*Card `json:"cards"`
}

type Player struct {
	ID int64 `json:"id" bson:"_id"`
	Keeper []*Card
}

type Card struct {
	ID string
	Type string
	Description string
	RuleType string
	AdditionalInfo string
	AdditionalCards []string
}