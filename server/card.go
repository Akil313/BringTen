package main

import (
	"fmt"
	"strings"
)

func stringToCard(cardString string) card {
	if cardString == "" {
		return card{value: "", suit: ""}
	}

	cardSlice := strings.Split(cardString, "x")
	if len(cardSlice) < 2 {
		return card{value: "", suit: ""}
	}

	return card{value: cardSlice[0], suit: cardSlice[1]}
}

type card struct {
	value    string
	suit     string
	playedBy *gamePlayer
}

func (c card) cardValue() int {
	valueMap := map[string]int{
		"Z": 0, "": 0, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6,
		"7": 7, "8": 8, "9": 9, "10": 10,
		"J": 11, "Q": 12, "K": 13, "A": 14,
	}
	return valueMap[c.value]
}

func (c card) cardGameValue() int {
	valueMap := map[string]int{
		"Z": 0, "2": 0, "3": 0, "4": 0, "5": 0, "6": 0,
		"7": 0, "8": 0, "9": 0, "10": 10,
		"J": 1, "Q": 2, "K": 3, "A": 4,
	}
	return valueMap[c.value]
}

func (c card) MarshalJSON() ([]byte, error) {
	if c.value == "" && c.suit == "" {
		return []byte(`"back"`), nil
	}
	return []byte(fmt.Sprintf(`"%sx%s"`, c.value, c.suit)), nil
}
