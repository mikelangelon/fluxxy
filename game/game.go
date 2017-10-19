package game

import (
	"game/card"
	"time"
	"game/player"
	"math/rand"
)

var Games map[int64]*Game
type Game struct {
	ID int64 `json:"id"`
	Players player.Players `json:"players"`
	Rules card.Cards `json:"rules"`
	Cards card.Cards `json:"cards"`
}

func StartGame(playerId int64) *Game {

	gameId := int64(1)
	players := []*player.Player{{ID: playerId}, {ID: playerId + 1}, {ID: playerId + 2}, {ID: playerId + 3}}
	game := &Game{ID: gameId, Players: players}

	cards := card.InitCards()
	MixToPlayers(cards, players)

	Games[gameId] = game

	return game
}

func PlayCard(gameId int64, playerId int64, cardId string)  {
	//validate card in player hand
	game := Games[gameId]

	player := game.Players.FindPlayer(playerId)
	card, position := player.Cards.FindCard(cardId)
	//do something depending of the type
	switch card.Type{

	}
	//Remove card from hand
	player.Cards = append(player.Cards[:position], player.Cards[position+1:]...)
}
func MixToPlayers(cards []*card.Card, players []*player.Player) {
	rand.Seed(time.Now().UnixNano())
	card.Shuffle(cards)

	for _, player := range players {
		player.Cards = cards[0:4]
		cards = cards[4:]
	}
}
