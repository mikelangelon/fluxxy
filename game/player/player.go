package player

import "game/card"

type Players []*Player
type Player struct {
	ID int64 `json:"id" bson:"_id"`
	Keeper card.Cards
	Cards card.Cards  `json:"cards"`
}

func (p Players) FindPlayer(playerId int64) *Player{
	for _, player := range p {
		if player.ID == playerId {
			return player
		}
	}
	return nil
}
