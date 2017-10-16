package mainmodel

import (
	"math/rand"
	"time"
)

func initCards() []*Card {
	cards := []*Card{}
	cards = append(cards,
		keeper("The Party"),
		keeper("Cookies"),
		keeper("Milk"),
		keeper("The Sun"),
		keeper("The Moon"),
		keeper("Chocolate"),
		keeper("Dreams"),
		keeper("Time"),
		keeper("Sleep"),
		keeper("Music"),
		keeper("The Toaster"),
		keeper("Money"),
		keeper("The Rocket"),
		keeper("Television"),
		keeper("Bread"),
		keeper("Love"),
		keeper("Peace"),

		rule("Draw 2", "Draw", "2"),
		rule("Draw 3", "Draw", "3"),
		rule("Draw 4", "Draw", "4"),
		rule("Draw 5", "Draw", "5"),
		rule("Hand Limit 0", "HandLimit", "0"),
		rule("Hand Limit 1", "HandLimit", "1"),
		rule("Hand Limit 2", "HandLimit", "2"),
		rule("Play 2", "Play", "2"),
		rule("Play 3", "Play", "3"),
		rule("Play 4", "Play", "4"),
		rule("Play All", "Play", "All"),
		rule("Play All But 1", "Play", "But1"),

		goal("Squishy Chocolate", []string{"Chocolate", "The Sun"}),
		goal("Star Gazing", []string{"Cosmos", "Eye"}),
		goal("Time is Money", []string{"Time", "Money"}),
		goal("Toast", []string{"Bread", "The Toaster"}),
	)
	return cards
}

func DoShuffle(cards []*Card) {
	rand.Seed(time.Now().UnixNano())
	Shuffle(cards)
}

func keeper(Descripion string) *Card {
	return &Card{ID: Descripion, Type: "Keeper", Description: Descripion}
}

func rule(Descripion string, RuleType string, AdditionalInfo string) *Card {
	return &Card{ID: Descripion, Type: "Rule", Description: Descripion, RuleType: RuleType, AdditionalInfo: AdditionalInfo}
}

func goal(Descripion string, AdditionalCards []string) *Card {
	return &Card{ID: Descripion, Type: "Goal", Description: Descripion, AdditionalCards: AdditionalCards}
}

func Shuffle(a []*Card) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}
