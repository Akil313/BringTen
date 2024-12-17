//go:build ignore

package main

import (
	"errors"
	"fmt"
	"io"
	// "log"
	"net/http"
	"os"
)

type card struct {
	suit  string
	value string
}

type deck struct {
	cards []card
}

func newDeck() *deck {
	suits := []string{"C", "H", "S", "D"}
	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	cards := []card{}
	for _, v := range values {
		for _, s := range suits {
			cards = append(cards, card{value: v, suit: s})
		}
	}

	return &deck{cards: cards}
}

func greetPlayer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("player has joined")
	io.WriteString(w, "Welcome to BringTen!\n")
}

func createDeck(w http.ResponseWriter, r *http.Request) {

	tempDeck := newDeck()

	fmt.Printf("new deck made")
	io.WriteString(w, fmt.Sprintf("%+v\n", tempDeck))
}

func main() {

	mux := http.NewServeMux() // Routes request to approriate handler
	mux.HandleFunc("/", greetPlayer)
	mux.HandleFunc("/deck", createDeck)

	fmt.Println("Server is up!")

	err := http.ListenAndServe(":8080", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
