package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type card struct {
	suit  string
	value string
}

type deck struct {
	cards []card
}

// Fisher-Yates shuffle
func (d *deck) shuffle() error {

	if len(d.cards) != 52 {
		return errors.New("deck is not properly filled\n")
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	shufIndex := r.Perm(len(d.cards))

	newDeck := make([]card, 52)

	for i, v := range shufIndex {
		newDeck[i] = d.cards[v]
	}

	d.cards = newDeck

	return nil
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

type room struct {
	id      string
	name    string
	canJoin bool
	deck    deck
}

type roomManager struct {
	rooms map[string]*room
}

func (rm *roomManager) generateRoomId(length int) (string, error) {

	if length < 4 {
		return "", errors.New("Length of room id must be more than 3 characters")
	}

	charset := "1234567890abcdefghijklmnopqrstuvwxyz"
	roomId := make([]byte, length)

	for i := 0; i < length; i++ {
		roomId[i] = charset[rand.Intn(len(charset))]
	}

	return string(roomId), nil
}

func (rm *roomManager) addNewRoom(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Method is: %v\n", r.Method)
	roomId, _ := rm.generateRoomId(4)

	defer r.Body.Close()
	var body struct {
		RoomName string `json:"room_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	roomName := body.RoomName

	rm.rooms[roomId] = &room{id: roomId, name: roomName, canJoin: true}

	io.WriteString(w, fmt.Sprintf("Room {%v} has been created with ID: %v\n", body.RoomName, roomId))
	fmt.Printf("room{%v} created\n", roomId)
}

func (rm *roomManager) joinRoom(w http.ResponseWriter, r *http.Request) {

	roomId := r.PathValue("id")
	fmt.Printf(roomId)
	roomName := rm.rooms[roomId].name
	io.WriteString(w, fmt.Sprintf("You have joined {%v} with ID: {%v}\n", roomName, roomId))
	fmt.Printf("player joined room {%v}\n", roomId)
}

func greetPlayer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("player has joined\n")
	io.WriteString(w, "Welcome to BringTen!\n")
}

func createDeck(w http.ResponseWriter, r *http.Request) {

	tempDeck := newDeck()
	err := tempDeck.shuffle()

	if err != nil {
		fmt.Printf("%v", err)
	}

	io.WriteString(w, fmt.Sprintf("%+v\n", tempDeck))
}

func main() {

	roomManager := &roomManager{
		rooms: make(map[string]*room),
	}

	mux := http.NewServeMux() // Routes request to approriate handler
	mux.HandleFunc("/", greetPlayer)
	mux.HandleFunc("/deck", createDeck)
	mux.HandleFunc("POST /room/{id}/join", roomManager.joinRoom)
	mux.HandleFunc("POST /room", roomManager.addNewRoom)

	fmt.Println("Server is up!")

	err := http.ListenAndServe(":8080", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
