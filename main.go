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
	value string
	suit  string
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

func (d *deck) newDeck() *deck {
	suits := []string{"C", "H", "S", "D"}
	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	cards := []card{}
	for _, v := range values {
		for _, s := range suits {
			cards = append(cards, card{value: v, suit: s})
		}
	}

	d.cards = cards

	return d
}

func (d *deck) shareCards(numCards int) []card {

	if numCards <= 0 {
		fmt.Printf("Number of cards must be non zero and non negative")
		return []card{}
	}

	newHand, newDeck := d.cards[:numCards], d.cards[numCards:]
	d.cards = newDeck

	return newHand
}

type gamePlayer struct {
	id   string
	name string
	hand []card
}

type room struct {
	id             string
	name           string
	deck           *deck
	currPlayerTurn int
	players        []*gamePlayer
	team1          []*gamePlayer
	team2          []*gamePlayer
	team1Score     int
	team2Score     int
	currTrump      card
	playerBeg      bool
	playerStay     bool
	actionList     [][]string
}

// Adds player to the room
func (r *room) addPlayer(player *gamePlayer) error {

	r.players = append(r.players, player)
	return nil
}

// Checks if a player's id is already in the room
func (r *room) isPlayerInRoom(player *gamePlayer) bool {

	for _, jp := range r.players {
		if jp.id == player.id {
			return true
		}
	}
	return false
}

func (r *room) checkIsRoomFull() bool {
	return len(r.players) == 4
}

func (r *room) startGame() {

	r.team1 = []*gamePlayer{r.players[0], r.players[2]}
	r.team2 = []*gamePlayer{r.players[1], r.players[3]}

	r.team1Score = 0
	r.team1Score = 0

	r.currPlayerTurn = rand.Intn(4)
	r.deck = r.deck.newDeck()
	r.deck.shuffle()

	r.players[r.currPlayerTurn].hand = r.deck.shareCards(6)

	dealerIdx := (r.currPlayerTurn-2)%4 + 1
	r.players[dealerIdx].hand = r.deck.shareCards(6)

	r.currTrump = r.deck.shareCards(1)[0]

	return
}

func (r *room) playerBegAction(pId string) {
	fmt.Printf("player {%v} has begged\n", pId)

	r.playerBeg = true
	r.playerStay = false
	return
}

func (r *room) playerStayAction(pId string) {
	fmt.Printf("player {%v} has stayed\n", pId)

	r.playerBeg = false
	r.playerStay = true
	return
}

func (r *room) dealerGiveOneAction(pId string) {
	fmt.Printf("dealer {%v} gave one point\n", pId)

	return
}

func (r *room) processAction(playerId, playerAction, cardPlayed string) {

	switch playerAction {
	case "BEG":
		r.playerBegAction(playerId)
	case "STAY":
		r.playerStayAction(playerId)
	case "GIVE_ONE":
		r.dealerGiveOneAction(playerId)
	case "GO_AGAIN":
		return
	case "PLAY_CARD":
		return
	default:
	}

	return
}

func (r *room) brodcastState() {

}

type roomManager struct {
	rooms map[string]*room
}

// Return a random string as id for the room based on length. If id is given, the same id is returned
func (rm *roomManager) generateRoomId(length int, userRoomId string) (string, error) {

	if length < 4 {
		return "", errors.New("Length of room id must be more than 3 characters")
	}

	if userRoomId != "" {
		return userRoomId, nil
	}

	charset := "1234567890abcdefghijklmnopqrstuvwxyz"
	roomId := make([]byte, length)

	for i := 0; i < length; i++ {
		roomId[i] = charset[rand.Intn(len(charset))]
	}

	return string(roomId), nil
}

// Create a new room and assign host to the new room
func (rm *roomManager) addNewRoom(w http.ResponseWriter, r *http.Request) {

	//Define structure of request body
	defer r.Body.Close()
	var request struct {
		RoomId   string `json:"room_id"`
		RoomName string `json:"room_name"`
		HostId   string `json:"host_id"`
		HostName string `json:"host_name"`
	}

	// Decode body of request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	fmt.Printf("creating room with id: {%v} \n", request)

	//put variables from body into variables
	userRoomId := request.RoomId
	userRoomName := request.RoomName
	hostId := request.HostId
	hostName := request.HostName

	roomId, _ := rm.generateRoomId(4, userRoomId)

	//Check if room exists. If it does, write that room exists and return
	//NOTE: Need to modify response so that if the room exists, it returns an approprite response
	_, ok := rm.rooms[roomId]
	if ok {
		io.WriteString(w, fmt.Sprintf("Room with id{%v} already exists\n", roomId))
		fmt.Printf("room{%v} already exists\n", roomId)
		return
	}

	//Create and store room pointer!
	//NOTE: Since we create a room pointer, modifying the newRoom variable changes the value in the array as well
	newRoom := &room{id: roomId, name: userRoomName, players: []*gamePlayer{}, deck: &deck{}}
	rm.rooms[roomId] = newRoom

	//Player/Host joins the room they created
	hostGamePlayer := &gamePlayer{id: hostId, name: hostName}
	newRoom.addPlayer(hostGamePlayer)

	// Send response of room id and room name to user
	response := map[string]string{
		"room_id":   roomId,
		"room_name": userRoomName,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	// io.WriteString(w, fmt.Sprintf("Room {%v} has been created with ID: %v\n", body.RoomName, roomId))

	// Debugging
	fmt.Printf("room{%v} created\n", roomId)
	fmt.Printf("%+v\n", newRoom)
}

// Add player to the respective room's list of players
func (rm *roomManager) joinRoom(w http.ResponseWriter, r *http.Request) {

	roomId := r.PathValue("id")
	currRoom, roomFound := rm.rooms[roomId]

	//NOTE: Add response for room not found
	if roomFound != true {
		fmt.Println("The room listed was not found")
		return
	}

	roomName := rm.rooms[roomId].name

	if currRoom.checkIsRoomFull() {
		fmt.Println("The room is full!")
		return
	}

	var joinRoomReqBody struct {
		PlayerId   string `json:"player_id"`
		PlayerName string `json:"player_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&joinRoomReqBody); err != nil {
		fmt.Printf("There was an error with the json body:\n%v\n", err)
		return
	}

	playerId := joinRoomReqBody.PlayerId
	playerName := joinRoomReqBody.PlayerName
	newPlayer := &gamePlayer{id: playerId, name: playerName}

	if err := currRoom.isPlayerInRoom(newPlayer); err == true {
		fmt.Println("This player is already in the room")
		return
	}

	err := rm.rooms[roomId].addPlayer(newPlayer)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	// Send response of room id and room name to user
	response := map[string]string{
		"room_id":   roomId,
		"room_name": roomName,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	io.WriteString(w, fmt.Sprintf("You have joined {%v} with ID: {%v}\n", roomName, roomId))
	fmt.Printf("player {%v} joined room {%v}\n", playerName, roomId)
}

func (rm *roomManager) startGame(w http.ResponseWriter, r *http.Request) {
	fmt.Println("starting room game")

	roomId := r.PathValue("id")
	currRoom, roomFound := rm.rooms[roomId]

	//NOTE: Add response for room not found
	if roomFound != true {
		fmt.Println("The room listed was not found")
		return
	}

	//NOTE: Add response for room is not full
	if currRoom.checkIsRoomFull() == false {
		fmt.Println("The must be full to start")
		return
	}

	var requestBody struct {
		HostId string `json:"host_id"`
	}
	//NOTE: Return appropriate response in response writer for client
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Printf("Error decoding request body: %v", requestBody)
		return
	}

	//Start game
	currRoom.startGame()
	for i, p := range currRoom.players {
		fmt.Printf("Player {%v}: %+v\n", i, p)
	}

	return
}

func (rm *roomManager) processGameAction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("process game action")

	roomId := r.PathValue("id")
	currRoom, roomFound := rm.rooms[roomId]

	//NOTE: Add response for room not found
	if roomFound != true {
		fmt.Println("The room listed was not found")
		return
	}

	var requestBody struct {
		Action     string `json:"action"`
		PlayerId   string `json:"player_id"`
		CardPlayed string `json:"card_played"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Printf("There was an error decoding: %v", err)
		return
	}

	playerId := requestBody.PlayerId
	playerAction := requestBody.Action
	cardPlayed := requestBody.CardPlayed
	currRoom.processAction(playerId, playerAction, cardPlayed)

}

func (rm *roomManager) sseGameStateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("broadcasting game state to player")

	// Set http headers required for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Needed locally for CORS requests
	//NOTE: Come back to this to understand why this setting is here
	w.Header().Set("Access-Control-Allow-Origin", "*")

	roomId := r.PathValue("roomId")
	playerId := r.PathValue("playerId")
	currRoom, roomFound := rm.rooms[roomId]

	//NOTE: Add response for room not found
	if roomFound != true {
		fmt.Println("The room listed was not found")
		return
	}

	var requestBody struct {
		PlayerId string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Printf("There was an error decoding: %v", err)
		return
	}

}

// Simple Greeting when the home page is run
func greetPlayer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("player has joined\n")
	io.WriteString(w, "Welcome to BringTen!\n")
}

func main() {

	roomManager := &roomManager{
		rooms: make(map[string]*room),
	}

	mux := http.NewServeMux() // Routes request to approriate handler
	mux.HandleFunc("/", greetPlayer)
	mux.HandleFunc("POST /room", roomManager.addNewRoom)
	mux.HandleFunc("POST /room/{id}/join", roomManager.joinRoom)
	mux.HandleFunc("POST /room/{id}/start", roomManager.startGame)
	mux.HandleFunc("POST /room/{roomId}/{playerId}/state", roomManager.sseGameStateHandler)

	fmt.Println("Server is up!")

	err := http.ListenAndServe(":8080", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
