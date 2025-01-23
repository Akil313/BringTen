package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

type httpResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   *errorInfo  `json:"error"`
}

type errorInfo struct {
	Code    string `json:"code"`
	Details string `json:"details"`
}

func sendResponse(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}, error *errorInfo) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := httpResponse{
		Success: success,
		Message: message,
		Data:    data,
		Error:   error,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func stringToCard(cardString string) card {
	if cardString == "" {
		return card{value: "Z", suit: "Z"}
	}

	cardSlice := strings.Split(cardString, "x")
	if len(cardSlice) < 2 {
		return card{value: "Z", suit: "Z"}
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
	return []byte(fmt.Sprintf(`"%sx%s"`, c.value, c.suit)), nil
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

func (d *deck) shareSelectedCards(idx int) []card {

	selectedCards := []card{
		{value: "J", suit: "C"},
		{value: "K", suit: "C"},
		{value: "4", suit: "H"},
		{value: "6", suit: "H"},
	}

	return []card{selectedCards[idx]}
}

type gamePlayer struct {
	pos        int
	id         string
	name       string
	hand       []card
	validHand  []card
	team       *team
	clientChan chan *gameState
}

func (p *gamePlayer) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, p.id)), nil
}

func (p *gamePlayer) removeCardFromHand(playedCard card) {
	if len(p.hand) < 1 {
		return
	}

	cardHandIdx := slices.Index(p.hand, playedCard)
	cardValidHandIdx := slices.Index(p.validHand, playedCard)

	if cardHandIdx == -1 || cardValidHandIdx == -1 {
		return
	}

	copy(p.hand[cardHandIdx:], p.hand[cardHandIdx+1:])
	p.hand[len(p.hand)-1] = card{}
	p.hand = p.hand[:len(p.hand)-1]
	//
	// copy(p.validHand[cardValidHandIdx:], p.validHand[cardValidHandIdx+1:])
	// p.validHand[len(p.validHand)-1] = card{}
	// p.validHand = p.validHand[:len(p.validHand)-1]
}

type gameState struct {
	Name       string        `json:"name"`
	Hand       []card        `json:"hand"`
	ValidHand  []card        `json:"valid_hand"`
	Deck       *deck         `json:"deck"`
	PlayerTurn int           `json:"curr_turn"`
	Dealer     int           `json:"dealer"`
	Players    []*gamePlayer `json:"players"`
	Team1Score int           `json:"team_1_score"`
	Team2Score int           `json:"team_2_score"`
	Trump      card          `json:"trump"`
	Lift       []card        `json:"lift"`
	PlayerBeg  bool          `json:"player_beg"`
	RoundStart bool          `json:"round_start"`
	PlayerStay bool          `json:"player_stay"`
	Winner     string        `json:"winner"`
}

type team struct {
	name    string
	players []*gamePlayer
	score   int
	lift    []card
}

type room struct {
	id                  string
	name                string
	deck                *deck
	playerTurn          int
	dealerIdx           int
	roundFirstPlayerIdx int
	players             []*gamePlayer
	teams               []*team
	callCard            card
	trump               card
	playerBeg           bool
	playerStay          bool
	actionList          [][]string
	roundStart          bool
	lift                []card
	highCard            card
	lowCard             card
	jackPlayed          bool
	jackPoint           *team
	hangJackPoint       *team
	winner              *team
}

// Adds player to the room
func (r *room) addPlayer(player *gamePlayer) error {

	r.players = append(r.players, player)
	return nil
}

// Checks if a player's id is already in the room
func (r *room) isPlayerInRoom(id string) (*gamePlayer, int, bool) {

	for i, jp := range r.players {
		if jp.id == id {
			return r.players[i], i, true
		}
	}
	return &gamePlayer{}, -1, false
}

func (r *room) checkIsRoomFull() bool {
	return len(r.players) == 4
}

func mod(a, b int) int {
	return ((a % b) + b) % b
}

func (r *room) isOnlySuitInHand(hand []card) bool {

	sameSuit := hand[0].suit

	for _, c := range hand {
		if c.suit != sameSuit {
			return false
		}
		continue
	}
	return true
}

// Is the card the same suit as the call card
func (r *room) isCallSuit(c card) bool {
	nilCard := card{value: "Z", suit: "Z"}
	if r.callCard == nilCard {
		return true
	}

	return c.suit == r.callCard.suit
}

func (r *room) isCallSuitInHand(hand []card) bool {

	for _, c := range hand {
		if c.suit == r.callCard.suit {
			return true
		}
	}
	return false
}

func (r *room) isHighestTrumpInLift(c card) bool {

	for _, liftCard := range r.lift {
		if liftCard.suit != r.trump.suit {
			continue
		}

		if c.cardValue() < liftCard.cardValue() {
			return false
		}
	}
	return true
}

func (r *room) highestCardInLift() card {

	if len(r.lift) < 1 {
		return card{}
	}

	cardList := func(lift []card) (ret []card) {
		for _, c := range lift {
			if c.suit == r.trump.suit || c.suit == r.callCard.suit {
				ret = append(ret, c)
			}
		}
		return ret
	}(r.lift)

	trumpCards := func(lift []card) (ret []card) {
		for _, c := range lift {
			if c.suit == r.trump.suit {
				ret = append(ret, c)
			}
		}
		return ret
	}(r.lift)

	if len(trumpCards) > 0 {
		cardList = trumpCards
	}

	highestCard := slices.MaxFunc(cardList, func(a, b card) int {
		return cmp.Compare(a.cardValue(), b.cardValue())
	})

	return highestCard
}

func (r *room) validCards(hand []card) []card {

	validHand := []card{}

	if len(hand) < 1 {
		return validHand
	}

	if len(r.lift) == 0 {
		return hand
	}

	if r.isOnlySuitInHand(hand) == true {
		return hand
	}

	for _, c := range hand {

		if r.isCallSuit(c) == true {
			validHand = append(validHand, c)
			continue
		}

		if c.suit != r.trump.suit {
			if !r.isCallSuitInHand(hand) {
				validHand = append(validHand, c)
			}
			continue
		}

		if r.isHighestTrumpInLift(c) {
			validHand = append(validHand, c)
			continue
		}
	}
	return validHand
}

// This function will return a hand with cards if players are allowed to see their cards
// If they are not, they will get an empty hand
func (r *room) canPlayerSeeHand(playerIdx int, hand []card) []card {

	if r.roundStart == true {
		return hand
	}

	if playerIdx == r.dealerIdx || playerIdx == r.roundFirstPlayerIdx {
		return hand
	}

	return []card{}
}

func (r *room) startGame() {

	for _, p := range r.players {
		_, p.pos, _ = r.isPlayerInRoom(p.id)
	}

	team1 := &team{
		name:    "team1",
		players: []*gamePlayer{r.players[0], r.players[2]},
		score:   0,
	}
	r.teams = append(r.teams, team1)

	team2 := &team{
		name:    "team2",
		players: []*gamePlayer{r.players[1], r.players[3]},
		score:   0,
	}
	r.teams = append(r.teams, team2)

	teamList := []*team{team1, team2}
	for i := range 4 {
		r.players[i].team = teamList[mod(i, 2)]
	}

	r.dealerIdx = rand.Intn(4)
	r.roundFirstPlayerIdx = mod((r.dealerIdx + 1), 4)
	r.playerTurn = r.roundFirstPlayerIdx
	r.deck = r.deck.newDeck()
	r.deck.shuffle()

	for i := range 4 {
		p := r.players[mod((r.playerTurn+i), 4)]
		p.hand = append(p.hand, r.deck.shareCards(1)...)
	}

	fmt.Printf("dealerIdx: %v\n", r.dealerIdx)
	fmt.Printf("firstPlayerIdx: %v\n", r.roundFirstPlayerIdx)

	fmt.Printf("flipping trump\n")
	r.trump = r.deck.shareCards(1)[0]
	fmt.Printf("New Trump: %v\n", r.trump)
	r.checkKickPoints()

	return
}

func (r *room) playerBegAction(player *gamePlayer) {
	fmt.Printf("player {%v} has begged\n", player.id)

	if r.roundStart == true {
		fmt.Printf("these actions cannot be taken if the round has started\n")
		return
	}

	if player != r.players[r.playerTurn] {
		fmt.Printf("player that wasn't current tried to beg\n")
		return
	}

	if r.playerBeg == true || r.playerStay == true {
		fmt.Println("player has already begged or stayed")
		return
	}

	r.playerBeg = true
	r.playerStay = false

	r.broadcastState()
	return
}

func (r *room) playerStayAction(player *gamePlayer) {
	fmt.Printf("player {%v} has stayed\n", player.id)

	if r.roundStart == true {
		fmt.Printf("these actions cannot be taken if the round has started\n")
		return
	}

	if player != r.players[r.playerTurn] {
		fmt.Printf("player that wasn't current tried to beg\n")
		return
	}

	if r.playerBeg == true || r.playerStay == true {
		fmt.Println("player has already begged or stayed")
	}

	r.playerBeg = false
	r.playerStay = true

	r.broadcastState()
	return
}

func (r *room) dealerGiveOneAction(player *gamePlayer) {
	fmt.Printf("dealer {%v} gave one point\n", player.id)

	if r.roundStart == true {
		fmt.Printf("these actions cannot be taken if the round has started\n")
		return
	}

	if player != r.players[r.dealerIdx] {
		fmt.Printf("player that wasn't current dealer tried to go again\n")
		return
	}

	if r.playerBeg != true {
		fmt.Printf("cannot go again if player hasnt begged\n")
		return
	}

	r.players[r.playerTurn].team.score += 1
	r.roundStart = true

	r.broadcastState()
	return
}

func (r *room) dealerGoAgain(player *gamePlayer) {
	fmt.Printf("dealer {%v} go again\n", player.id)

	if r.roundStart == true {
		fmt.Printf("cannot go again if the round has started\n")
		return
	}

	if player != r.players[r.dealerIdx] {
		fmt.Printf("player that wasn't current dealer tried to go again\n")
		return
	}

	if r.playerBeg != true {
		fmt.Printf("cannot go again if player hasnt begged\n")
		return
	}

	// Keeping suit of trump to check if next trump is the same as first
	startTrump := r.trump

	for startTrump.suit == r.trump.suit && len(r.deck.cards) > 4 {
		for i := range 4 {
			p := r.players[mod((r.playerTurn+i), 4)]

			fmt.Printf("deal 1 cards to player: %s\n", p.id)
			p.hand = append(p.hand, r.deck.shareCards(3)...)
		}

		fmt.Printf("flipping trump\n")
		r.trump = r.deck.shareCards(1)[0]
		fmt.Printf("New Trump: %v\n", r.trump)
		r.checkKickPoints()
	}

	for startTrump.suit == r.trump.suit && len(r.deck.cards) > 0 {
		fmt.Printf("flipping trump\n")
		r.trump = r.deck.shareCards(1)[0]
		fmt.Printf("New Trump: %v\n", r.trump)
		r.checkKickPoints()
	}

	r.roundStart = true

	if startTrump.suit == r.trump.suit {
		r.setupNextRound()
	}

	r.broadcastState()
	return
}

func (r *room) checkKickPoints() {

	if r.trump == (card{}) {
		return
	}

	switch r.trump.suit {
	case "J":
		r.players[r.dealerIdx].team.score += 3
	case "6":
		r.players[r.dealerIdx].team.score += 2
	case "A":
		r.players[r.dealerIdx].team.score += 1
	default:
		r.players[r.dealerIdx].team.score += 0
	}

	return
}

func (r *room) checkHighPoint(playedCard card) {
	fmt.Printf("Check High Point. Played: %+v\tTrump: %+v\n", playedCard.suit, r.trump.suit)

	if playedCard.suit != r.trump.suit {
		fmt.Println("Not right suit")
		return
	}

	if r.highCard == (card{}) {
		fmt.Println("No card set for high card")
		r.highCard = playedCard
		return
	}

	if playedCard.cardValue() < r.highCard.cardValue() {
		fmt.Printf("Look at values: %v\tTrump: %v\n", playedCard.cardValue(), r.highCard.cardValue())
		return
	}

	fmt.Println("Give high point")
	r.highCard = playedCard
	return
}

func (r *room) checkLowPoint(playedCard card) {
	fmt.Printf("Check Low Point. Played: %+v\tTrump: %+v\n", playedCard.suit, r.trump.suit)

	if playedCard.suit != r.trump.suit {
		return
	}

	if r.lowCard == (card{}) {
		fmt.Println("No card set for low card")
		r.lowCard = playedCard
		return
	}

	if playedCard.cardValue() > r.lowCard.cardValue() {
		fmt.Printf("Look at values: %v\tTrump: %v\n", playedCard.cardValue(), r.lowCard.cardValue())
		return
	}

	fmt.Println("Give low point")
	r.lowCard = playedCard
	return
}

func (r *room) checkJackPoint(player *gamePlayer, playedCard card) {
	if playedCard.suit != r.trump.suit {
		return
	}

	if r.jackPlayed == true {
		return
	}

	if playedCard.cardValue() != 11 {
		return
	}

	r.jackPlayed = true
	r.jackPoint = player.team
	return
}

func (r *room) checkHangJackPoint() {

	var jackIdx int
	if jackIdx = slices.IndexFunc(r.lift, func(c card) bool {
		return c.value == "J" && c.suit == r.trump.suit
	}); jackIdx == -1 {
		return
	}

	jackCard := r.lift[jackIdx]

	if r.isHighestTrumpInLift(jackCard) {
		fmt.Println("Jack highest in lift")
		return
	}

	var highestTrump card
	for _, c := range r.lift {
		if r.isHighestTrumpInLift(c) {
			highestTrump = c
		}
	}

	if jackCard.playedBy.team == highestTrump.playedBy.team {
		return
	}

	r.hangJackPoint = highestTrump.playedBy.team
	r.jackPoint = &team{}
	return
}

func (r *room) isRoundOver() bool {

	for _, p := range r.players {
		if len(p.hand) > 0 {
			return false
		}
	}
	return true
}

func (r *room) addGamePointScore() {

	team1Score := 0
	team2Score := 0

	for _, c := range r.players[0].team.lift {
		team1Score += c.cardGameValue()
	}

	for _, c := range r.players[1].team.lift {
		team2Score += c.cardGameValue()
	}

	if team1Score > team2Score {
		fmt.Printf("Giving game point to: %v\n", r.players[0].team.name)
		r.players[0].team.score += 1
	} else if team2Score > team1Score {
		fmt.Printf("Giving game point to: %v\n", r.players[1].team.name)
		r.players[1].team.score += 1
	} else {
		fmt.Printf("Game point is a draw.\nGiving game point to: %v\n", r.players[mod(r.dealerIdx+1, 4)].team.name)
		r.players[mod(r.dealerIdx+1, 4)].team.score += 1
	}
}

func (r *room) isGameOver() bool {

	if r.players[0].team.score >= 2 {
		fmt.Printf("%v is the winner!", r.players[0].team.name)
		r.winner = r.players[0].team
		return true
	}

	if r.players[1].team.score >= 2 {
		fmt.Printf("%v is the winner!", r.players[1].team.name)
		r.winner = r.players[1].team
		return true
	}

	return false
}

func (r *room) cleanUpRound() {
	fmt.Println("Cleaning up round")

	//Add High point
	if r.highCard.playedBy != nil {
		fmt.Printf("Giving hight point to: %v\n", r.highCard.playedBy.team.name)
		r.highCard.playedBy.team.score += 1
	}
	if r.isGameOver() {
		return
	}

	//Add Low point
	if r.lowCard.playedBy != nil {
		fmt.Printf("Giving low point to: %v\n", r.lowCard.playedBy.team.name)
		r.lowCard.playedBy.team.score += 1
	}
	if r.isGameOver() {
		return
	}

	//Add HangJack point
	if r.hangJackPoint != nil {
		fmt.Printf("Giving hang jack point to: %v\n", r.hangJackPoint.name)
		r.hangJackPoint.score += 3
	}
	if r.isGameOver() {
		return
	}

	//Add Jack point
	if r.jackPoint != nil {
		fmt.Printf("Giving jack point to: %v\n\n", r.jackPoint.name)
		r.jackPoint.score += 1
	}
	if r.isGameOver() {
		return
	}

	//Add Game point
	r.addGamePointScore()
	if r.isGameOver() {
		return
	}

	r.setupNextRound()
}

func (r *room) setupNextRound() {

	fmt.Println("Setting up next room")
	r.roundStart = false
	r.callCard = card{}
	r.trump = card{}
	r.playerBeg = false
	r.playerStay = false
	r.highCard = card{}
	r.lowCard = card{}
	r.jackPlayed = false
	r.jackPoint = nil
	r.hangJackPoint = nil
	r.winner = nil

	r.lift = []card{}
	r.dealerIdx = mod(r.dealerIdx+1, 4)
	r.playerTurn = mod(r.dealerIdx+1, 4)
	r.deck = r.deck.newDeck()
	r.deck.shuffle()

	for i := range 4 {
		p := r.players[mod((r.playerTurn+i), 4)]
		p.hand = append(p.hand, r.deck.shareCards(1)...)
	}

	fmt.Printf("dealerIdx: %v\n", r.dealerIdx)
	fmt.Printf("firstPlayerIdx: %v\n", r.roundFirstPlayerIdx)

	fmt.Printf("flipping trump\n")
	r.trump = r.deck.shareCards(1)[0]
	fmt.Printf("New Trump: %v\n", r.trump)
	r.checkKickPoints()

	return
}

func (r *room) playCard(player *gamePlayer, cardString string) {

	if playerIdx := slices.Index(r.players, player); playerIdx != r.playerTurn {
		return
	}

	playedCard := stringToCard(cardString)

	if playedCard == (card{}) {
		return
	}

	isValid := slices.Contains(player.validHand, playedCard)
	if isValid == false {
		return
	}

	player.removeCardFromHand(playedCard)
	playedCard.playedBy = player
	if len(r.lift) == 0 {
		r.callCard = playedCard
	}

	r.lift = append(r.lift, playedCard)

	//NOTE: Check for hang Jack
	r.checkHighPoint(playedCard)
	r.checkLowPoint(playedCard)
	r.checkJackPoint(player, playedCard)

	r.playerTurn = mod(r.playerTurn+1, 4)

	if len(r.lift) == 4 {
		r.checkHangJackPoint()

		highestCard := r.highestCardInLift()
		r.playerTurn = slices.Index(r.players, highestCard.playedBy)
		highestCard.playedBy.team.lift = append(highestCard.playedBy.team.lift, r.lift...)
		r.lift = []card{}
	}

	if r.isRoundOver() {
		r.cleanUpRound()
	}

	r.broadcastState()
}

func (r *room) broadcastState() {
	fmt.Println("\nbroadcasting state to players from room")

	for _, player := range r.players {
		player.validHand = r.validCards(player.hand)

		newGameState := &gameState{
			Name:       player.name,
			Hand:       r.canPlayerSeeHand(player.pos, player.hand),
			ValidHand:  r.canPlayerSeeHand(player.pos, player.validHand),
			Deck:       r.deck,
			Dealer:     r.dealerIdx,
			PlayerTurn: r.playerTurn,
			Players:    r.players,
			Team1Score: r.teams[0].score,
			Team2Score: r.teams[1].score,
			Trump:      r.trump,
			Lift:       r.lift,
			PlayerBeg:  r.playerBeg,
			PlayerStay: r.playerStay,
			RoundStart: r.roundStart,
			Winner: func() string {
				if r.winner != nil {
					return r.winner.name
				}
				return "None"
			}(),
		}

		player.clientChan <- newGameState
	}
}

func (r *room) processAction(player *gamePlayer, playerAction, cardPlayed string) {
	fmt.Printf("procees: %v, %v, %v\n", player.id, playerAction, cardPlayed)

	switch playerAction {
	case "BEG":
		r.playerBegAction(player)
	case "STAY":
		r.playerStayAction(player)
	case "GIVE_ONE":
		r.dealerGiveOneAction(player)
	case "GO_AGAIN":
		r.dealerGoAgain(player)
		return
	case "PLAY_CARD":
		r.playCard(player, cardPlayed)
		return
	default:
	}

	return
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

		message := "There was an error somewhere"
		error := &errorInfo{Code: "400", Details: "Invalid JSON format"}

		sendResponse(w, http.StatusBadRequest, false, message, nil, error)

		return
	}

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

		message := fmt.Sprintf("Room with id{%v} already exists\n", roomId)
		error := &errorInfo{Code: "400", Details: "Room Conflict"}

		sendResponse(w, http.StatusBadRequest, false, message, nil, error)
		fmt.Printf("room{%v} already exists\n", roomId)
		return
	}

	//Create and store room pointer!
	//NOTE: Since we create a room pointer, modifying the newRoom variable changes the value in the array as well
	newRoom := &room{
		id:      roomId,
		name:    userRoomName,
		players: []*gamePlayer{},
		deck:    &deck{},
	}
	rm.rooms[roomId] = newRoom

	//Player/Host joins the room they created
	hostGamePlayer := &gamePlayer{id: hostId, name: hostName, hand: []card{}, clientChan: make(chan *gameState)}
	newRoom.addPlayer(hostGamePlayer)

	// Send response of room id and room name to user
	response := map[string]string{
		"room_id":   newRoom.id,
		"room_name": newRoom.name,
	}

	message := "Room has been successfully created!"
	sendResponse(w, http.StatusOK, true, message, response, nil)

	fmt.Printf("creating room %v: %v \n", roomId, userRoomName)

}

// Add player to the respective room's list of players
func (rm *roomManager) joinRoom(w http.ResponseWriter, r *http.Request) {

	roomId := r.PathValue("id")
	currRoom, roomFound := rm.rooms[roomId]

	//NOTE: Add response for room not found
	if roomFound != true {
		fmt.Println("The room listed was not found")

		message := fmt.Sprintf("The room requested was not found")
		error := &errorInfo{Code: "400", Details: "The ID of the room submitted was not found in the list of active rooms"}

		sendResponse(w, http.StatusBadRequest, false, message, nil, error)
		return
	}

	if currRoom.checkIsRoomFull() {
		fmt.Println("The room is full!")
		message := fmt.Sprintf("The room you are trying to join is full")
		error := &errorInfo{Code: "400", Details: "The room the player is trying to join is full"}

		sendResponse(w, http.StatusBadRequest, false, message, nil, error)
		return
	}

	var joinRoomReqBody struct {
		PlayerId   string `json:"player_id"`
		PlayerName string `json:"player_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&joinRoomReqBody); err != nil {
		fmt.Printf("There was an error with the json body:\n%v\n", err)

		message := fmt.Sprint("There was an internal error")
		errorDetails := fmt.Sprint("There was an error with the json body formatting")
		error := &errorInfo{Code: "400", Details: errorDetails}

		sendResponse(w, http.StatusBadRequest, false, message, nil, error)
		return
	}

	playerId := joinRoomReqBody.PlayerId
	playerName := joinRoomReqBody.PlayerName
	newPlayer := &gamePlayer{id: playerId, name: playerName, hand: []card{}, clientChan: make(chan *gameState)}

	if _, _, err := currRoom.isPlayerInRoom(playerId); err == true {
		fmt.Println("This player is already in the room")

		message := fmt.Sprint("You have already joined this room")
		errorDetails := fmt.Sprint("The player trying to join is already listed to be in the room")
		error := &errorInfo{Code: "400", Details: errorDetails}

		sendResponse(w, http.StatusBadRequest, false, message, nil, error)
		return
	}

	err := rm.rooms[roomId].addPlayer(newPlayer)
	if err != nil {
		fmt.Printf("%v", err)

		message := fmt.Sprint("There was an error somewhere")
		errorDetails := fmt.Sprint("An error occured when trying to add the room to the list of rooms")
		error := &errorInfo{Code: "400", Details: errorDetails}

		sendResponse(w, http.StatusBadRequest, false, message, nil, error)
		return
	}

	// Send response of room id and room name to user
	response := map[string]string{
		"room_id":   currRoom.id,
		"room_name": currRoom.name,
	}
	message := "You have successfully joined the room! :)"
	sendResponse(w, http.StatusOK, true, message, response, nil)

	fmt.Printf("player {%v} joined room {%v}\n", playerName, roomId)
}

func (rm *roomManager) startGame(w http.ResponseWriter, r *http.Request) {
	fmt.Println("starting room game")

	roomId := r.PathValue("id")
	currRoom, roomFound := rm.rooms[roomId]

	//NOTE: Add response for room not found
	if roomFound != true {

		message := fmt.Sprint("There was an error somewhere")
		errorDetails := fmt.Sprint("An error because the id of the room that tried to start was not found")
		error := &errorInfo{Code: "400", Details: errorDetails}

		sendResponse(w, http.StatusBadRequest, false, message, nil, error)
		fmt.Println("The room listed was not found")
		return
	}

	//NOTE: Add response for room is not full
	if currRoom.checkIsRoomFull() == false {

		message := fmt.Sprint("Room cant be started if it isnt full")
		errorDetails := fmt.Sprint("An error because the room tried to start without enough players")
		error := &errorInfo{Code: "400", Details: errorDetails}

		sendResponse(w, http.StatusBadRequest, false, message, nil, error)
		fmt.Println("The room must be full to start")
		return
	}

	var requestBody struct {
		HostId string `json:"host_id"`
	}
	//NOTE: Return appropriate response in response writer for client
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {

		message := fmt.Sprint("There was an internal error")
		errorDetails := fmt.Sprint("There was an error with the json body formatting")
		error := &errorInfo{Code: "400", Details: errorDetails}

		sendResponse(w, http.StatusBadRequest, false, message, nil, error)

		fmt.Printf("Error decoding request body: %v", requestBody)
		return
	}

	//Start game
	currRoom.startGame()
	for i, p := range currRoom.players {
		fmt.Printf("\nPlayer {%v}: %+v\n", i, p)
	}

	// Send response of room id and room name to user
	response := map[string]string{}
	message := "The room has started! >:)"
	currRoom.broadcastState()
	sendResponse(w, http.StatusOK, true, message, response, nil)

	return
}

func (rm *roomManager) processGameAction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("process game action")

	roomId := r.PathValue("roomId")
	playerId := r.PathValue("playerId")
	currRoom, roomFound := rm.rooms[roomId]

	//NOTE: Add response for room not found
	if roomFound != true {

		message := fmt.Sprint("The room for that action could not be found")
		errorDetails := fmt.Sprint("The id for the room that the action occured in could not be found")
		error := &errorInfo{Code: "400", Details: errorDetails}

		sendResponse(w, http.StatusBadRequest, false, message, nil, error)

		fmt.Println("The room listed was not found")
		return
	}

	player, _, _ := currRoom.isPlayerInRoom(playerId)

	var requestBody struct {
		Action     string `json:"action"`
		CardPlayed string `json:"card_played"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Printf("There was an error decoding: %v", err)
		return
	}

	playerAction := requestBody.Action
	cardPlayed := requestBody.CardPlayed
	currRoom.processAction(player, playerAction, cardPlayed)

	message := "Action Successful"
	response := map[string]string{}
	sendResponse(w, http.StatusOK, true, message, response, nil)
}

func (rm *roomManager) sseGameStateHandler(w http.ResponseWriter, r *http.Request) {
	roomId := r.PathValue("roomId")
	playerId := r.PathValue("playerId")
	currRoom, roomFound := rm.rooms[roomId]

	//NOTE: Add response for room not found
	if roomFound != true {
		fmt.Println("The room listed was not found")
		return
	}

	player, _, _ := currRoom.isPlayerInRoom(playerId)

	// Set http headers required for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Needed locally for CORS requests
	//NOTE: Come back to this to understand why this setting is here
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	//Create channel for client disconnection
	clientGone := r.Context().Done()

	rc := http.NewResponseController(w)
	for {
		select {
		case <-clientGone:
			fmt.Println("Client disconnected")
			return
		case gameState := <-player.clientChan:
			//send event to client
			jsonBytes, err := json.Marshal(&gameState)

			if err != nil {
				fmt.Println("There was an error with the JSON conversion")
				return
			}

			_, err = fmt.Fprintf(w, "data: %+v\n\n", string(jsonBytes))

			if err != nil {
				return
			}
			err = rc.Flush()
			if err != nil {
				return
			}
		}
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
	mux.HandleFunc("POST /room/{roomId}/{playerId}/action", roomManager.processGameAction)
	mux.HandleFunc("GET /room/{roomId}/{playerId}/state", roomManager.sseGameStateHandler)

	fmt.Println("Server is up!")

	err := http.ListenAndServe(":8080", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
