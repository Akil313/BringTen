package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var SCORE_LIMIT int = 6
var HAND_SIZE int = 6
var EXPIRED_TIME float64 = 15.0
var allowedOrigins = []string{
	"http://165.227.221.32:3000",
}

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
	if c.value == "" && c.suit == "" {
		return []byte(`""`), nil
	}
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
	Pos        int    `json:"pos"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	hand       []card
	validHand  []card
	team       *team
	clientChan chan *gameState
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
	Position   int           `json:"position"`
	RoomName   string        `json:"room_name"`
	Hand       []card        `json:"hand"`
	ValidHand  []card        `json:"valid_hand"`
	Deck       int           `json:"deck"`
	PlayerTurn int           `json:"curr_turn"`
	Dealer     int           `json:"dealer"`
	Players    []*gamePlayer `json:"players"`
	Team1Score int           `json:"team_1_score"`
	Team2Score int           `json:"team_2_score"`
	Trump      card          `json:"trump"`
	Lift       []card        `json:"lift"`
	PlayerBeg  bool          `json:"player_beg"`
	RoundStart bool          `json:"round_start"`
	GameStart  bool          `json:"game_start"`
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
	host                *gamePlayer
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
	gameStart           bool
	lift                []card
	highCard            card
	lowCard             card
	jackPlayed          bool
	jackPoint           *team
	hangJackPoint       *team
	winner              *team
	lastActionTime      time.Time
}

func (r *room) updateLastActionTime() error {

	r.lastActionTime = time.Now()

	return nil
}

func (r *room) isRoomExpired() bool {

	elapsed := time.Since(r.lastActionTime)
	minsElapsed := elapsed.Minutes()

	if minsElapsed >= 30.0 {
		return true
	}

	return false
}

// Adds player to the room
func (r *room) addPlayer(player *gamePlayer) error {

	r.players = append(r.players, player)
	r.updateLastActionTime()
	return nil
}

// Checks if a player's id is already in the room
func (r *room) isPlayerInRoom(id string) (*gamePlayer, int, bool) {

	for i, jp := range r.players {
		if jp.Id == id {
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
		_, p.Pos, _ = r.isPlayerInRoom(p.Id)
	}

	for i := range 4 {
		r.players[i].team = r.teams[mod(i, 2)]
	}

	r.gameStart = true
	r.roundStart = false
	r.dealerIdx = rand.Intn(4)
	r.roundFirstPlayerIdx = mod((r.dealerIdx + 1), 4)
	r.playerTurn = r.roundFirstPlayerIdx
	r.deck = r.deck.newDeck()
	r.deck.shuffle()

	for i := range 4 {
		p := r.players[mod((r.playerTurn+i), 4)]
		p.hand = append(p.hand, r.deck.shareCards(HAND_SIZE)...)
	}

	fmt.Printf("dealerIdx: %v\n", r.dealerIdx)
	fmt.Printf("firstPlayerIdx: %v\n", r.roundFirstPlayerIdx)

	fmt.Printf("flipping trump\n")
	r.trump = r.deck.shareCards(1)[0]
	fmt.Printf("New Trump: %v\n", r.trump)
	r.checkKickPoints()

	r.updateLastActionTime()

	return
}

func (r *room) playerBegAction(player *gamePlayer) {
	fmt.Printf("player {%v} has begged\n", player.Id)

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
	fmt.Printf("player {%v} has stayed\n", player.Id)

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
	r.roundStart = true

	r.broadcastState()
	return
}

func (r *room) dealerGiveOneAction(player *gamePlayer) {
	fmt.Printf("dealer {%v} gave one point\n", player.Id)

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
	fmt.Printf("dealer {%v} go again\n", player.Id)

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

			fmt.Printf("deal 1 cards to player: %s\n", p.Id)
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
		fmt.Println("First trump played. Give high point")
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
		fmt.Println("Not right suit")
		return
	}

	if r.lowCard == (card{}) {
		fmt.Println("First trump played. Give low point")
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

	if r.players[0].team.score >= SCORE_LIMIT {
		fmt.Printf("%v is the winner!", r.players[0].team.name)
		r.winner = r.players[0].team
		return true
	}

	if r.players[1].team.score >= SCORE_LIMIT {
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
	r.roundFirstPlayerIdx = mod((r.dealerIdx + 1), 4)
	r.playerTurn = r.roundFirstPlayerIdx
	r.deck = r.deck.newDeck()
	r.deck.shuffle()

	for i := range 4 {
		p := r.players[mod((r.playerTurn+i), 4)]
		p.hand = append(p.hand, r.deck.shareCards(HAND_SIZE)...)
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
	fmt.Printf("\nbroadcasting state from room: %v, roomsize: %v\n", r.name, len(r.players))

	for _, player := range r.players {
		player.validHand = r.validCards(player.hand)

		newGameState := &gameState{
			RoomName:   r.name,
			Name:       player.Name,
			Position:   player.Pos,
			Hand:       r.canPlayerSeeHand(player.Pos, player.hand),
			ValidHand:  r.canPlayerSeeHand(player.Pos, player.validHand),
			Deck:       len(r.deck.cards),
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
			GameStart:  r.gameStart,
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
	fmt.Printf("procees: %v, %v, %v\n", player.Id, playerAction, cardPlayed)

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

	r.updateLastActionTime()

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

func (rm *roomManager) generatePlayerId(length int) (string, error) {

	if length < 4 {
		return "", errors.New("Length of player id must be more than 3")
	}

	charset := "1234567890abcdefghijklmnopqrstuvwxyz"
	playerId := make([]byte, length)

	for i := 0; i < length; i++ {
		playerId[i] = charset[rand.Intn(len(charset))]
	}

	return string(playerId), nil
}

// Create a new room and assign host to the new room
func (rm *roomManager) addNewRoom(w http.ResponseWriter, r *http.Request) {

	enableCors(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

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
	hostId, _ := rm.generatePlayerId(6)
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
		teams: []*team{
			{
				name:    "team1",
				players: []*gamePlayer{},
				score:   0,
			},
			{
				name:    "team2",
				players: []*gamePlayer{},
				score:   0,
			},
		},
	}

	newRoom.updateLastActionTime()

	rm.rooms[roomId] = newRoom

	//Player/Host joins the room they created
	hostGamePlayer := &gamePlayer{Id: hostId, Name: hostName, hand: []card{}, clientChan: make(chan *gameState)}
	newRoom.host = hostGamePlayer
	newRoom.addPlayer(hostGamePlayer)

	// Send response of room id and room name to user
	response := map[string]string{
		"room_id":   newRoom.id,
		"room_name": newRoom.name,
		"host_id":   hostId,
		"host_name": hostName,
	}

	message := "Room has been successfully created!"
	sendResponse(w, http.StatusOK, true, message, response, nil)

	fmt.Printf("creating room %v: %v \n", roomId, userRoomName)

	return
}

// Add player to the respective room's list of players
func (rm *roomManager) joinRoom(w http.ResponseWriter, r *http.Request) {

	enableCors(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	roomId := vars["id"]
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
		// PlayerId   string `json:"player_id"`
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

	// playerId := joinRoomReqBody.PlayerId
	playerId, _ := rm.generatePlayerId(6)
	playerName := joinRoomReqBody.PlayerName
	newPlayer := &gamePlayer{Id: playerId, Name: playerName, hand: []card{}, clientChan: make(chan *gameState)}

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
		"room_id":     currRoom.id,
		"room_name":   currRoom.name,
		"player_id":   playerId,
		"player_name": playerName,
	}
	message := "You have successfully joined the room! :)"
	sendResponse(w, http.StatusOK, true, message, response, nil)

	fmt.Printf("player {%v} joined room {%v}\n", playerName, roomId)
}

func (rm *roomManager) deleteRoomById(roomKey string) error {

	delete(rm.rooms, roomKey)

	return nil
}

func (rm *roomManager) deleteRoom(w http.ResponseWriter, r *http.Request) {
	enableCors(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	roomId := vars["id"]

	_, roomFound := rm.rooms[roomId]

	if roomFound != true {
		fmt.Println("The room listed was not found")
		message := "Room %v does not exist :("
		sendResponse(w, http.StatusNotFound, false, message, nil, nil)
		return
	}

	rm.deleteRoomById(roomId)

	fmt.Printf("Deleted room with id %v\n", roomId)

	message := "Room %v has been removed :)"
	sendResponse(w, http.StatusOK, true, message, nil, nil)
}

func (rm *roomManager) checkAllRoomsExpired() error {

	countRemoved := 0
	for _, room := range rm.rooms {
		if room.isRoomExpired() {
			rm.deleteRoomById(room.id)
			countRemoved++
		}
	}

	fmt.Println("# of expired rooms removed: ", countRemoved)
	return nil
}

func (rm *roomManager) startGame(w http.ResponseWriter, r *http.Request) {

	enableCors(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	fmt.Println("starting room game")

	vars := mux.Vars(r)
	roomId := vars["id"]
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

	enableCors(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Println("process game action")

	vars := mux.Vars(r)
	roomId := vars["roomId"]
	playerId := vars["playerId"]
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

func (rm *roomManager) getRooms(w http.ResponseWriter, r *http.Request) {

	enableCors(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	type simpleRoomDetails struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Host       string `json:"host"`
		NumPlayers int    `json:"numPlayers"`
	}

	rooms := map[string]simpleRoomDetails{}

	for _, r := range rm.rooms {
		rooms[r.id] = simpleRoomDetails{ID: r.id, Name: r.name, Host: r.host.Name, NumPlayers: len(r.players)}
	}

	message := "Rooms returned! :)"
	sendResponse(w, http.StatusOK, true, message, rooms, nil)

	return
}

func (rm *roomManager) sseGameStateHandler(w http.ResponseWriter, r *http.Request) {

	enableCors(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	roomId := vars["roomId"]
	playerId := vars["playerId"]
	currRoom, roomFound := rm.rooms[roomId]

	fmt.Printf("%v, %v", roomId, playerId)

	//NOTE: Add response for room not found
	if roomFound != true {
		fmt.Println("The room listed was not found")
		message := "Room could not be found"
		sendResponse(w, http.StatusNotFound, false, message, nil, nil)
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

	// Trigger initial broadcast when a player connects
	go currRoom.broadcastState()

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
	enableCors(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Printf("Welcome to BringTen! Server Side\n")

	message := "Welcome to BringTen!\n"
	sendResponse(w, http.StatusOK, true, message, nil, nil)

}

func enableCors(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	if isAllowedOrigin(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func isAllowedOrigin(origin string) bool {
	for _, o := range allowedOrigins {
		if o == origin {
			return true
		}
	}
	return false
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
			w.Header().Set("Vary", "Origin")
		}

		// Short-circuit preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the actual handler
		next.ServeHTTP(w, r)
	})
}

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or couldn't load the .env file")
	}

	origins := os.Getenv("ORIGINS")
	if origins != "" {
		allowedOrigins = strings.Split(origins, ",")
	}

}

func main() {

	roomManager := &roomManager{
		rooms: make(map[string]*room),
	}

	r := mux.NewRouter()
	r.HandleFunc("/", greetPlayer).Methods("GET")
	r.HandleFunc("/rooms", roomManager.addNewRoom).Methods("POST")
	r.HandleFunc("/rooms/{id}/join", roomManager.joinRoom).Methods("POST")
	r.HandleFunc("/rooms/{id}/start", roomManager.startGame).Methods("POST")
	r.HandleFunc("/rooms/{id}/delete", roomManager.deleteRoom).Methods("DELETE")
	r.HandleFunc("/rooms/{roomId}/{playerId}/action", roomManager.processGameAction).Methods("POST", "OPTIONS")
	r.HandleFunc("/rooms", roomManager.getRooms).Methods("GET")
	r.HandleFunc("/rooms/{roomId}/{playerId}/state", roomManager.sseGameStateHandler).Methods("GET")

	fmt.Println("Server is up!")

	//Periodically check for expired rooms to be deleted
	expiredRoomTicker := time.NewTicker(time.Duration(EXPIRED_TIME) * time.Minute) // Ticks every minute
	defer expiredRoomTicker.Stop()

	go func() {
		for {
			select {
			case t := <-expiredRoomTicker.C:
				roomManager.checkAllRoomsExpired()
				fmt.Println("Chekcing expired rooms at: ", t)
			}
		}
	}()

	err := http.ListenAndServe(":8080", withCORS(r))

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
