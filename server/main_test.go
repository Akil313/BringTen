package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

const HOST = "http://localhost"
const PORT = "8080"

type response struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
	Error   string            `json:"error"`
}

func TestOnePlayerCreateRoom(t *testing.T) {
	fmt.Println("Testing One Player Creating a Room")
	posturl := HOST + ":" + PORT + "/room"

	postBodyPlayer1, _ := json.Marshal(map[string]string{
		"host_id":   "rt01pl01",
		"host_name": "Arin",
		"room_id":   "rt01",
		"room_name": "Game Grumps",
	})
	responseBody := bytes.NewBuffer(postBodyPlayer1)

	res, _ := http.Post(posturl, "application/json", responseBody)

	var response struct {
		RoomId   string `json:"room_id"`
		RoomName string `json:"room_name"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		fmt.Printf("Error in decoding: \n%v\n", err)
		return
	}
	fmt.Printf("RoomID: %v\tRoom Name: %v\n", response.RoomId, response.RoomName)
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

func TestBegGiveAndPlayOneTurn(t *testing.T) {

	type gameState struct {
		Name           string   `json:"name"`
		Hand           []string `json:"hand"`
		ValidHand      []string `json:"valid_hand"`
		CurrPlayerTurn int      `json:"curr_turn"`
		Dealer         int      `json:"dealer"`
		Players        []string `json:"players"`
		Team1Score     int      `json:"team_1_score"`
		Team2Score     int      `json:"team_2_score"`
		Trump          string   `json:"trump"`
		Lift           []string `json:"lift"`
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

	for _, joinP := range roomPlayers {
		if joinP.id == p1Id {
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
				fmt.Printf("Error in decoding %v: \n%v\n", joinP.id, err)
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
		roomPlayers[p1Id].state = gs1
		fmt.Printf("\n%v\n", mainGameState)

		gs2 := <-roomPlayers[p2Id].playerChan
		roomPlayers[p2Id].state = gs2
		fmt.Printf("%v\n", gs2.Hand)

		gs3 := <-roomPlayers[p3Id].playerChan
		roomPlayers[p3Id].state = gs3
		fmt.Printf("%v\n", gs3.Hand)

		gs4 := <-roomPlayers[p4Id].playerChan
		roomPlayers[p4Id].state = gs4
		fmt.Printf("%v\n", gs4.Hand)
	}()

	currPlayerIdx := int(mainGameState.CurrPlayerTurn)
	dealerIdx := int(mainGameState.Dealer)
	playerIds := mainGameState.Players

	// First player beg
	postUrl := HOST + ":" + PORT + "/room/" + roomId + "/" + playerIds[currPlayerIdx] + "/action"
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
		roomPlayers[p1Id].state = gs1
		fmt.Printf("\n%v\n", mainGameState)

		gs2 := <-roomPlayers[p2Id].playerChan
		roomPlayers[p2Id].state = gs2
		fmt.Printf("%v\n", gs2.Hand)

		gs3 := <-roomPlayers[p3Id].playerChan
		roomPlayers[p3Id].state = gs3
		fmt.Printf("%v\n", gs3.Hand)

		gs4 := <-roomPlayers[p4Id].playerChan
		roomPlayers[p4Id].state = gs4
		fmt.Printf("%v\n", gs4.Hand)
	}()

	// Dealer give one
	postUrl = HOST + ":" + PORT + "/room/" + roomId + "/" + playerIds[dealerIdx] + "/action"
	postBody, _ = json.Marshal(map[string]string{"action": "GO_AGAIN", "card_played": ""})
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
		roomPlayers[p1Id].state = gs1
		fmt.Printf("\n%v\n", mainGameState)

		gs2 := <-roomPlayers[p2Id].playerChan
		roomPlayers[p2Id].state = gs2
		fmt.Printf("%v\n", gs2.Hand)

		gs3 := <-roomPlayers[p3Id].playerChan
		roomPlayers[p3Id].state = gs3
		fmt.Printf("%v\n", gs3.Hand)

		gs4 := <-roomPlayers[p4Id].playerChan
		roomPlayers[p4Id].state = gs4
		fmt.Printf("%v\n", gs4.Hand)
	}()

	for i := 0; i < 4; i++ {
		// Player 1 Play
		playCardRequest(roomId, playerIds[currPlayerIdx], roomPlayers[playerIds[currPlayerIdx]].state.ValidHand[0])

		//First Play
		func() {
			gs1 := <-roomPlayers[p1Id].playerChan
			mainGameState = gs1
			roomPlayers[p1Id].state = gs1
			fmt.Printf("\nTeam1:%v    Trump:%v    Team2:%v    CurrPlayer:%v    Winner:%v\n", gs1.Team1Score, gs1.Trump, gs1.Team2Score, gs1.CurrPlayerTurn, gs1.Winner)
			fmt.Printf("Lift:%v", gs1.Lift)
			fmt.Printf("\nHand:%v\tValid Hand:%v\n", gs1.Hand, gs1.ValidHand)

			gs2 := <-roomPlayers[p2Id].playerChan
			roomPlayers[p2Id].state = gs2
			fmt.Printf("\nHand:%v\tValid Hand:%v\n", gs2.Hand, gs2.ValidHand)

			gs3 := <-roomPlayers[p3Id].playerChan
			roomPlayers[p3Id].state = gs3
			fmt.Printf("\nHand:%v\tValid Hand:%v\n", gs3.Hand, gs3.ValidHand)

			gs4 := <-roomPlayers[p4Id].playerChan
			roomPlayers[p4Id].state = gs4
			fmt.Printf("\nHand:%v\tValid Hand:%v\n", gs4.Hand, gs4.ValidHand)
		}()

		currPlayerIdx = mainGameState.CurrPlayerTurn
	}
}
