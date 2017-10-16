package mainmodel

import (
	"testing"
	"fmt"
)

func TestShuffle(t *testing.T) {
	cards := initCards()
	DoShuffle(cards)
	for _, c :=range cards{
		fmt.Println(c.Description)
	}
}
