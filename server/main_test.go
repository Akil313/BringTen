package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

const HOST = "http://localhost"
const PORT = "8080"

type response struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
	Error   string            `json:"error"`
}

// decodeEnvelope reads the server's standard {success, message, data, error}
// envelope and returns the Data payload typed as T. Any transport problem,
// non-200 status, decode error, or server-reported failure is a t.Fatalf.
func decodeEnvelope[T any](t *testing.T, res *http.Response) T {
	t.Helper()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("expected 200, got %d; body: %s", res.StatusCode, string(body))
	}

	var env struct {
		Success bool       `json:"success"`
		Message string     `json:"message"`
		Data    T          `json:"data"`
		Error   *errorInfo `json:"error"`
	}
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if !env.Success {
		t.Fatalf("server reported failure: %s (%+v)", env.Message, env.Error)
	}
	return env.Data
}

// createRoomData is the payload shape returned inside `data` for POST /rooms.
type createRoomData struct {
	RoomId   string `json:"room_id"`
	RoomName string `json:"room_name"`
	HostId   string `json:"host_id"`
	HostName string `json:"host_name"`
}

func TestOnePlayerCreateRoom(t *testing.T) {
	fmt.Println("Testing One Player Creating a Room")
	posturl := HOST + ":" + PORT + "/rooms"

	postBodyPlayer1, _ := json.Marshal(map[string]string{
		"host_id":   "rt01pl01",
		"host_name": "Arin",
		"room_id":   "rt01",
		"room_name": "Game Grumps",
	})
	responseBody := bytes.NewBuffer(postBodyPlayer1)

	res, err := http.Post(posturl, "application/json", responseBody)
	if err != nil {
		t.Fatalf("POST %s failed: %v", posturl, err)
	}

	data := decodeEnvelope[createRoomData](t, res)

	if data.RoomId == "" || data.RoomName == "" {
		t.Fatalf("empty room id or name in response: %+v", data)
	}

	fmt.Printf("RoomID: %v\tRoom Name: %v\n", data.RoomId, data.RoomName)
}

func TestRoomCreateFourPlayerJoin(t *testing.T) {

	fmt.Println("Testing one player creating room and 3 joining")
	createRoomUrl := HOST + ":" + PORT + "/room"

	// JSON object for Host creating room
	postBodyPlayer1, _ := json.Marshal(map[string]string{
		"room_id":   "rt02",
		"room_name": "Akatsuki",
		"host_id":   "rt02pl01",
		"host_name": "Itachi",
	})
	requestBody := bytes.NewBuffer(postBodyPlayer1)

	type createRoomRes struct {
		RoomId   string `json:"room_id"`
		RoomName string `json:"room_name"`
	}
	parsedResponse := createRoomRes{}

	//Send create room request
	res, _ := http.Post(createRoomUrl, "application/json", requestBody)
	if err := json.NewDecoder(res.Body).Decode(&parsedResponse); err != nil {
		fmt.Printf("Error in decoding: \n%v\n", err)
	}
	fmt.Printf("%+v\n", parsedResponse)

	roomId := parsedResponse.RoomId
	// roomName := parsedResponse.RoomName

	// Create second player body
	postBodyPlayer2, _ := json.Marshal(map[string]string{
		"player_id":   "rt02pl02",
		"player_name": "Pain",
	})
	requestBody = bytes.NewBuffer(postBodyPlayer2)

	type joinRoomRes struct {
		RoomId   string `json:"room_id"`
		RoomName string `json:"room_name"`
	}
	parsedJoinResponse := joinRoomRes{}

	joinRoomUrl := HOST + ":" + PORT + "/room/" + roomId + "/join"
	res, _ = http.Post(joinRoomUrl, "application/json", requestBody)

	if err := json.NewDecoder(res.Body).Decode(&parsedJoinResponse); err != nil {
		fmt.Printf("Error in decoding 1: \n%v\n", err)
	}
	fmt.Printf("%+v\n", parsedJoinResponse)

	// Create third player body
	postBodyPlayer3, _ := json.Marshal(map[string]string{
		"player_id":   "rt02pl03",
		"player_name": "Kisame",
	})
	requestBody = bytes.NewBuffer(postBodyPlayer3)

	parsedJoinResponse = joinRoomRes{}

	res, _ = http.Post(joinRoomUrl, "application/json", requestBody)

	if err := json.NewDecoder(res.Body).Decode(&parsedJoinResponse); err != nil {
		fmt.Printf("Error in decoding: \n%v\n", err)
	}
	fmt.Printf("%+v\n", parsedJoinResponse)

	// Create third player body
	postBodyPlayer4, _ := json.Marshal(map[string]string{
		"player_id":   "rt02pl04",
		"player_name": "Hidan",
	})
	requestBody = bytes.NewBuffer(postBodyPlayer4)

	parsedJoinResponse = joinRoomRes{}

	res, _ = http.Post(joinRoomUrl, "application/json", requestBody)

	if err := json.NewDecoder(res.Body).Decode(&parsedJoinResponse); err != nil {
		fmt.Printf("Error in decoding: \n%v\n", err)
	}
	fmt.Printf("%+v\n", parsedJoinResponse)

	return
}

func TestStartFullGameRoom(t *testing.T) {

	fmt.Println("Testing one player creating room and 3 joining")

	createRoomUrl := HOST + ":" + PORT + "/room"

	// JSON object for Host creating room
	postBodyPlayer1, _ := json.Marshal(map[string]string{
		"room_id":   "rt03",
		"room_name": "Spice Girls",
		"host_id":   "rt03pl01",
		"host_name": "Baby",
	})
	requestBody := bytes.NewBuffer(postBodyPlayer1)

	parsedResponse := response{}

	//Send create room request
	res, _ := http.Post(createRoomUrl, "application/json", requestBody)
	if err := json.NewDecoder(res.Body).Decode(&parsedResponse); err != nil {
		fmt.Printf("Error in decoding: \n%v\n", err)
	}
	fmt.Printf("%s\n", parsedResponse.Message)
	roomId := parsedResponse.Data["room_id"]
	// roomName := parsedResponse.RoomName

	// Create second player body
	postBodyPlayer2, _ := json.Marshal(map[string]string{
		"player_id":   "rt03pl02",
		"player_name": "Scary",
	})
	requestBody = bytes.NewBuffer(postBodyPlayer2)

	parsedJoinResponse := response{}

	joinRoomUrl := HOST + ":" + PORT + "/room/" + roomId + "/join"
	res, _ = http.Post(joinRoomUrl, "application/json", requestBody)

	if err := json.NewDecoder(res.Body).Decode(&parsedJoinResponse); err != nil {
		fmt.Printf("Error in decoding 1: \n%v\n", err)
	}
	fmt.Printf("%s\n", parsedJoinResponse.Message)

	// Create third player body
	postBodyPlayer3, _ := json.Marshal(map[string]string{
		"player_id":   "rt03pl03",
		"player_name": "Posh",
	})
	requestBody = bytes.NewBuffer(postBodyPlayer3)

	parsedJoinResponse = response{}

	res, _ = http.Post(joinRoomUrl, "application/json", requestBody)

	if err := json.NewDecoder(res.Body).Decode(&parsedJoinResponse); err != nil {
		fmt.Printf("Error in decoding: \n%v\n", err)
	}
	fmt.Printf("%s\n", parsedJoinResponse.Message)

	// Create fourth player body
	postBodyPlayer4, _ := json.Marshal(map[string]string{
		"player_id":   "rt03pl04",
		"player_name": "Sporty",
	})
	requestBody = bytes.NewBuffer(postBodyPlayer4)

	parsedJoinResponse = response{}

	res, _ = http.Post(joinRoomUrl, "application/json", requestBody)

	if err := json.NewDecoder(res.Body).Decode(&parsedJoinResponse); err != nil {
		fmt.Printf("Error in decoding: \n%v\n", err)
	}
	fmt.Printf("%s\n", parsedJoinResponse.Message)

	// Start room game

	startGameUrl := HOST + ":" + PORT + "/room/" + roomId + "/start"
	postBodyStartGame, _ := json.Marshal(map[string]string{"": ""})
	requestBody = bytes.NewBuffer(postBodyStartGame)

	res, _ = http.Post(startGameUrl, "application/json", requestBody)

	return
}

func TestBegAndGiveOne(t *testing.T) {

	type gameState struct {
		Name           string   `json:"name"`
		Hand           []string `json:"hand"`
		ValidHand      []string `json:"valid_hand"`
		CurrPlayerTurn int      `json:"curr_turn"`
		Dealer         int      `json:"dealer"`
		Players        []string `json:"players"`
		Team1Score     int      `json:"team_1_score"`
		Team2Score     int      `json:"team_2_score"`
		CurrTrump      string   `json:"trump"`
		PlayerBeg      bool     `json:"player_beg"`
		RoundStart     bool     `json:"round_start"`
		PlayerStay     bool     `json:"player_stay"`
		Winner         string   `json:"winner"`
	}

	type playerInfo struct {
		id         string
		name       string
		playerChan chan *gameState
		state      *gameState
	}

	fmt.Println("Testing 4 players joining, player begging and dealer giving one")

	createRoomUrl := HOST + ":" + PORT + "/room"

	// JSON object for Host creating room
	p1Id := "rt04pl01"
	postBodyPlayer1, _ := json.Marshal(map[string]string{
		"room_name": "Mystery Inc.",
		"host_id":   p1Id,
		"host_name": "Fred",
	})
	requestBody := bytes.NewBuffer(postBodyPlayer1)

	parsedResponse := response{}

	//Send create room request
	res, _ := http.Post(createRoomUrl, "application/json", requestBody)
	if err := json.NewDecoder(res.Body).Decode(&parsedResponse); err != nil {
		fmt.Printf("Error in decoding: \n%v\n", err)
	}
	fmt.Printf("\n%s\n", parsedResponse.Message)

	roomId := parsedResponse.Data["room_id"]

	////////////////////////////////////////////////////////
	p2Id := "rt04pl02"
	p3Id := "rt04pl03"
	p4Id := "rt04pl04"
	roomPlayers := map[string]*playerInfo{
		p1Id: {id: p1Id, name: "Fred", playerChan: make(chan *gameState), state: nil},
		p2Id: {id: p2Id, name: "Velma", playerChan: make(chan *gameState), state: nil},
		p3Id: {id: p3Id, name: "Daphne", playerChan: make(chan *gameState), state: nil},
		p4Id: {id: p4Id, name: "Shaggy", playerChan: make(chan *gameState), state: nil},
	}

	i := 0
	for _, joinP := range roomPlayers {
		if i == 0 {
			i++
			continue
		}

		func(p *playerInfo) {

			postBodyPlayer, _ := json.Marshal(map[string]string{
				"player_id":   p.id,
				"player_name": p.name,
			})
			requestBody = bytes.NewBuffer(postBodyPlayer)

			parsedJoinResponse := response{}

			joinRoomUrl := HOST + ":" + PORT + "/room/" + roomId + "/join"
			res, _ = http.Post(joinRoomUrl, "application/json", requestBody)

			if err := json.NewDecoder(res.Body).Decode(&parsedJoinResponse); err != nil {
				fmt.Printf("Error in decoding %v: \n%v\n", i, err)
			}
			fmt.Printf("\n%s\n", parsedJoinResponse.Message)
		}(joinP)
	}

	for _, p := range roomPlayers {
		go func(roomId, playerId string, playerChan chan<- *gameState) {
			resp, err := http.Get(fmt.Sprintf("http://localhost:8080/room/%s/%s/state", roomId, playerId))
			if err != nil {
				fmt.Println("Error connecting to SSE:", err)
				return
			}
			defer resp.Body.Close()

			fmt.Printf("\nConnected to SSE stream for player %s\n", playerId)

			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				line := scanner.Text()

				if len(line) > 6 && line[:5] == "data:" {
					//Parse JSON Payload
					data := line[5:]
					var newState gameState
					err := json.Unmarshal([]byte(data), &newState)
					if err != nil {
						fmt.Printf("Error parsing JSON for player %s: %v\n", playerId, err)
						continue
					}

					playerChan <- &newState
				}
			}

			if scanner.Err() != nil {
				fmt.Printf("Error reading SSE stream for player: %s", playerId)
			}

			close(playerChan)
		}(roomId, p.id, p.playerChan)
	}
	////////////////////////////////////////////////////////

	// Start room game
	startGameUrl := HOST + ":" + PORT + "/room/" + roomId + "/start"
	postBodyStartGame, _ := json.Marshal(map[string]string{"": ""})
	requestBody = bytes.NewBuffer(postBodyStartGame)

	res, _ = http.Post(startGameUrl, "application/json", requestBody)

	parsedJoinResponse := response{}
	if err := json.NewDecoder(res.Body).Decode(&parsedJoinResponse); err != nil {
		fmt.Printf("Error in decoding: \n%v\n", err)
	}
	fmt.Printf("\n%s\n", parsedJoinResponse.Message)

	// The game loop should include conditions to exit, based on game progress
	mainGameState := &gameState{}

	//Start Game
	func() {
		gs1 := <-roomPlayers[p1Id].playerChan
		mainGameState = gs1
		fmt.Printf("\n%v\n", mainGameState)

		gs2 := <-roomPlayers[p2Id].playerChan
		fmt.Printf("%v\n", gs2.Hand)

		gs3 := <-roomPlayers[p3Id].playerChan
		fmt.Printf("%v\n", gs3.Hand)

		gs4 := <-roomPlayers[p4Id].playerChan
		fmt.Printf("%v\n", gs4.Hand)
	}()

	firstPlayer := int(mainGameState.CurrPlayerTurn)
	dealer := int(mainGameState.Dealer)
	playerIds := mainGameState.Players

	// First player beg
	postUrl := HOST + ":" + PORT + "/room/" + roomId + "/" + playerIds[firstPlayer] + "/action"
	postBody, _ := json.Marshal(map[string]string{"action": "BEG", "card_played": ""})
	requestBody = bytes.NewBuffer(postBody)
	res, _ = http.Post(postUrl, "application/json", requestBody)
	parsedResponse = response{}
	if err := json.NewDecoder(res.Body).Decode(&parsedResponse); err != nil {
		fmt.Printf("Error in decoding: \n%v\n", err)
	}
	fmt.Printf("\n%s\n", parsedResponse.Message)

	//Beg
	func() {
		gs1 := <-roomPlayers[p1Id].playerChan
		mainGameState = gs1
		fmt.Printf("\n%v\n", mainGameState)

		gs2 := <-roomPlayers[p2Id].playerChan
		fmt.Printf("%v\n", gs2.Hand)

		gs3 := <-roomPlayers[p3Id].playerChan
		fmt.Printf("%v\n", gs3.Hand)

		gs4 := <-roomPlayers[p4Id].playerChan
		fmt.Printf("%v\n", gs4.Hand)
	}()

	// Dealer give one
	postUrl = HOST + ":" + PORT + "/room/" + roomId + "/" + playerIds[dealer] + "/action"
	postBody, _ = json.Marshal(map[string]string{"action": "GIVE_ONE", "card_played": ""})
	requestBody = bytes.NewBuffer(postBody)
	res, _ = http.Post(postUrl, "application/json", requestBody)
	parsedResponse = response{}
	if err := json.NewDecoder(res.Body).Decode(&parsedResponse); err != nil {
		fmt.Printf("Error in decoding: \n%v\n", err)
	}
	fmt.Printf("\n%s\n", parsedResponse.Message)

	//Give One
	func() {
		gs1 := <-roomPlayers[p1Id].playerChan
		mainGameState = gs1
		fmt.Printf("\n%v\n", mainGameState)

		gs2 := <-roomPlayers[p2Id].playerChan
		fmt.Printf("%v\n", gs2.Hand)

		gs3 := <-roomPlayers[p3Id].playerChan
		fmt.Printf("%v\n", gs3.Hand)

		gs4 := <-roomPlayers[p4Id].playerChan
		fmt.Printf("%v\n", gs4.Hand)
	}()
}

func playCardRequest(roomId, playerId string, cardPlayed string) {
	postUrl := HOST + ":" + PORT + "/room/" + roomId + "/" + playerId + "/action"
	postBody, _ := json.Marshal(map[string]string{"action": "PLAY_CARD", "card_played": cardPlayed})

	requestBody := bytes.NewBuffer(postBody)
	res, _ := http.Post(postUrl, "application/json", requestBody)
	parsedResponse := response{}
	if err := json.NewDecoder(res.Body).Decode(&parsedResponse); err != nil {
		fmt.Printf("Error in decoding: \n%v\n", err)
	}
	fmt.Printf("\n%s\n", parsedResponse.Message)
}

// -----------------------------------------------------------------------------
// Integration test helpers (used by TestBegGiveAndPlayOneTurn). These all fail
// loudly via t.Fatalf instead of silently printing an error and continuing.
// -----------------------------------------------------------------------------

// hostCreatedData mirrors the JSON shape returned inside `data` from
// POST /rooms. We capture host_id because the server generates it.
type hostCreatedData struct {
	RoomId   string `json:"room_id"`
	RoomName string `json:"room_name"`
	HostId   string `json:"host_id"`
	HostName string `json:"host_name"`
}

// joinedData mirrors the JSON shape returned inside `data` from
// POST /rooms/{id}/join. player_id is generated server-side.
type joinedData struct {
	RoomId     string `json:"room_id"`
	RoomName   string `json:"room_name"`
	PlayerId   string `json:"player_id"`
	PlayerName string `json:"player_name"`
}

// broadcastStatePlayer is the {pos, id, name} view of a gamePlayer that the
// server serializes inside gameState.Players.
type broadcastStatePlayer struct {
	Pos  int    `json:"pos"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

// broadcastState mirrors the gameState payload delivered over SSE. Fields
// typed as card on the server come through as plain strings because card
// has a custom MarshalJSON.
type broadcastState struct {
	Name       string                 `json:"name"`
	Position   int                    `json:"position"`
	RoomName   string                 `json:"room_name"`
	HostId     string                 `json:"host_id"`
	Hand       []string               `json:"hand"`
	ValidHand  []string               `json:"valid_hand"`
	Deck       int                    `json:"deck"`
	CurrTurn   int                    `json:"curr_turn"`
	Dealer     int                    `json:"dealer"`
	Players    []broadcastStatePlayer `json:"players"`
	Team1Score int                    `json:"team_1_score"`
	Team2Score int                    `json:"team_2_score"`
	Trump      string                 `json:"trump"`
	Lift       []string               `json:"lift"`
	PlayerBeg  bool                   `json:"player_beg"`
	RoundStart bool                   `json:"round_start"`
	GameStart  bool                   `json:"game_start"`
	PlayerStay bool                   `json:"player_stay"`
	Winner     string                 `json:"winner"`
}

// playerStream owns one SSE connection for a player. Ch receives every parsed
// state update; Last is the newest state the test has explicitly read.
type playerStream struct {
	Id   string
	Name string
	Ch   chan *broadcastState
	Last *broadcastState
}

func postCreateRoom(t *testing.T, roomName, hostName string) hostCreatedData {
	t.Helper()
	body, _ := json.Marshal(map[string]string{
		"room_name": roomName,
		"host_name": hostName,
	})
	res, err := http.Post(HOST+":"+PORT+"/rooms", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("POST /rooms: %v", err)
	}
	return decodeEnvelope[hostCreatedData](t, res)
}

func postJoinRoom(t *testing.T, roomId, playerName string) joinedData {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"player_name": playerName})
	url := HOST + ":" + PORT + "/rooms/" + roomId + "/join"
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	return decodeEnvelope[joinedData](t, res)
}

func postStartGame(t *testing.T, roomId string) {
	t.Helper()
	url := HOST + ":" + PORT + "/rooms/" + roomId + "/start"
	res, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	_ = decodeEnvelope[map[string]string](t, res)
}

func postAction(t *testing.T, roomId, playerId, action, cardPlayed string) {
	t.Helper()
	body, _ := json.Marshal(map[string]string{
		"action":      action,
		"card_played": cardPlayed,
	})
	url := HOST + ":" + PORT + "/rooms/" + roomId + "/" + playerId + "/action"
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("POST %s (%s): %v", url, action, err)
	}
	_ = decodeEnvelope[map[string]string](t, res)
}

// startSSEStream opens an SSE connection for a player and returns a stream
// that parses broadcasts onto Ch. It registers a cleanup to close the HTTP
// body when the test ends so the scanner goroutine exits cleanly.
func startSSEStream(t *testing.T, roomId, playerId, playerName string) *playerStream {
	t.Helper()

	url := fmt.Sprintf("%s:%s/rooms/%s/%s/state", HOST, PORT, roomId, playerId)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("SSE connect %s: %v", playerName, err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		t.Fatalf("SSE %s: expected 200, got %d", playerName, resp.StatusCode)
	}

	ps := &playerStream{
		Id:   playerId,
		Name: playerName,
		Ch:   make(chan *broadcastState, 64),
	}

	t.Cleanup(func() { resp.Body.Close() })

	go func() {
		defer close(ps.Ch)
		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1<<20)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data:") {
				continue
			}
			payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if payload == "" {
				continue
			}
			var s broadcastState
			if err := json.Unmarshal([]byte(payload), &s); err != nil {
				fmt.Printf("parse SSE payload for %s: %v\n", playerName, err)
				continue
			}
			ps.Ch <- &s
		}
	}()

	return ps
}

// awaitState blocks up to 5 seconds for at least one state, then coalesces
// any additional broadcasts that arrive within a short settle window. This
// prevents the test from operating on a broadcast that's already been
// superseded by a newer one still in flight.
func awaitState(t *testing.T, s *playerStream) *broadcastState {
	t.Helper()

	var state *broadcastState
	select {
	case st, ok := <-s.Ch:
		if !ok {
			t.Fatalf("SSE stream for %s closed unexpectedly", s.Name)
		}
		state = st
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for SSE state for %s", s.Name)
	}

	// Coalesce any near-simultaneous broadcasts so the test sees the latest.
	settle := time.After(100 * time.Millisecond)
	for {
		select {
		case st, ok := <-s.Ch:
			if !ok {
				s.Last = state
				return state
			}
			state = st
		case <-settle:
			s.Last = state
			return state
		}
	}
}

// syncAll waits for the next state on every stream and stores it on .Last.
func syncAll(t *testing.T, streams []*playerStream) {
	t.Helper()
	for _, s := range streams {
		awaitState(t, s)
	}
}

// TestBegGiveAndPlayOneTurn exercises the full bid + play path end-to-end:
//
//	4 players join -> SSE streams open -> game starts -> non-dealer begs ->
//	dealer gives one -> each player plays one card to complete a lift.
//
// It captures server-generated player ids from responses, matches the real
// broadcast payload shape, drains SSE backlog to the latest state after each
// action, and fails loudly on any transport / decode / protocol error.
func TestBegGiveAndPlayOneTurn(t *testing.T) {
	fmt.Println("Testing 4 players: bid (BEG + GIVE_ONE) and playing one lift")

	// ---- Create room and join the other three players -----------------------
	host := postCreateRoom(t, "Mystery Inc.", "Fred")
	roomId := host.RoomId

	velma := postJoinRoom(t, roomId, "Velma")
	daphne := postJoinRoom(t, roomId, "Daphne")
	shaggy := postJoinRoom(t, roomId, "Shaggy")

	// ---- Open an SSE stream per player -------------------------------------
	// Each SSE connect triggers `go room.broadcastState()` on the server, so
	// four backlog broadcasts will arrive before we do anything. We drain them
	// below with the first syncAll.
	fredStream := startSSEStream(t, roomId, host.HostId, "Fred")
	velmaStream := startSSEStream(t, roomId, velma.PlayerId, "Velma")
	daphneStream := startSSEStream(t, roomId, daphne.PlayerId, "Daphne")
	shaggyStream := startSSEStream(t, roomId, shaggy.PlayerId, "Shaggy")

	streams := []*playerStream{fredStream, velmaStream, daphneStream, shaggyStream}
	streamById := map[string]*playerStream{
		host.HostId:     fredStream,
		velma.PlayerId:  velmaStream,
		daphne.PlayerId: daphneStream,
		shaggy.PlayerId: shaggyStream,
	}

	// Consume the pre-game broadcasts so we operate on a clean baseline.
	syncAll(t, streams)

	// ---- Start the game -----------------------------------------------------
	postStartGame(t, roomId)
	syncAll(t, streams)

	state := fredStream.Last
	if !state.GameStart {
		t.Fatalf("expected game_start=true after /start, got %+v", state)
	}
	if len(state.Players) != 4 {
		t.Fatalf("expected 4 players in broadcast, got %d", len(state.Players))
	}
	if state.Trump == "" {
		t.Fatalf("expected a trump card after start, got empty")
	}

	currTurnIdx := state.CurrTurn
	dealerIdx := state.Dealer
	currPlayerId := state.Players[currTurnIdx].Id
	dealerId := state.Players[dealerIdx].Id

	fmt.Printf("\nAfter start: turn=%d(%s)  dealer=%d(%s)  trump=%s\n",
		currTurnIdx, state.Players[currTurnIdx].Name,
		dealerIdx, state.Players[dealerIdx].Name,
		state.Trump)

	// ---- Current (non-dealer) player begs ----------------------------------
	postAction(t, roomId, currPlayerId, "BEG", "")
	syncAll(t, streams)

	state = fredStream.Last
	if !state.PlayerBeg {
		t.Fatalf("expected player_beg=true after BEG, got %+v", state)
	}
	if state.RoundStart {
		t.Fatalf("round should not have started yet after BEG, got %+v", state)
	}

	// ---- Dealer awards the begging player's team a point -------------------
	postAction(t, roomId, dealerId, "GIVE_ONE", "")
	syncAll(t, streams)

	state = fredStream.Last
	if !state.RoundStart {
		t.Fatalf("expected round_start=true after GIVE_ONE, got %+v", state)
	}
	fmt.Printf("After bidding: team1=%d team2=%d trump=%s roundStart=%v\n",
		state.Team1Score, state.Team2Score, state.Trump, state.RoundStart)

	// ---- Play a full lift (four plays) -------------------------------------
	for i := 0; i < 4; i++ {
		turnIdx := state.CurrTurn
		if turnIdx < 0 || turnIdx >= len(state.Players) {
			t.Fatalf("invalid curr_turn %d (players: %v)", turnIdx, state.Players)
		}

		turnPlayerId := state.Players[turnIdx].Id
		turnStream := streamById[turnPlayerId]
		if turnStream == nil {
			t.Fatalf("no stream tracked for server-side player id %q", turnPlayerId)
		}
		if turnStream.Last == nil || len(turnStream.Last.ValidHand) == 0 {
			t.Fatalf("%s has no valid cards to play: %+v", turnStream.Name, turnStream.Last)
		}

		choice := turnStream.Last.ValidHand[0]
		fmt.Printf("play #%d: %s plays %s  (hand=%v valid=%v)\n",
			i+1, turnStream.Name, choice, turnStream.Last.Hand, turnStream.Last.ValidHand)

		postAction(t, roomId, turnPlayerId, "PLAY_CARD", choice)
		syncAll(t, streams)

		state = fredStream.Last
		fmt.Printf("  -> lift=%v  team1=%d  team2=%d  nextTurn=%d  winner=%s\n",
			state.Lift, state.Team1Score, state.Team2Score, state.CurrTurn, state.Winner)
	}

	// After a full lift the server clears the lift for the next trick (or, if
	// hands are empty, rolls into next round setup). Either outcome is fine;
	// we simply assert the test reached this point without hanging or erroring.
}
