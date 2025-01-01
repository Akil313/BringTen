package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	// "regexp"
	"testing"
)

const HOST = "http://localhost"
const PORT = "8080"

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
		"player_id":   "rt03pl02",
		"player_name": "Scary",
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
		"player_id":   "rt03pl03",
		"player_name": "Posh",
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
		"player_id":   "rt03pl04",
		"player_name": "Sporty",
	})
	requestBody = bytes.NewBuffer(postBodyPlayer4)

	parsedJoinResponse = joinRoomRes{}

	res, _ = http.Post(joinRoomUrl, "application/json", requestBody)

	if err := json.NewDecoder(res.Body).Decode(&parsedJoinResponse); err != nil {
		fmt.Printf("Error in decoding: \n%v\n", err)
	}
	fmt.Printf("%+v\n", parsedJoinResponse)

	// Start room game

	startGameUrl := HOST + ":" + PORT + "/room/" + roomId + "/start"
	postBodyStartGame, _ := json.Marshal(map[string]string{"": ""})
	requestBody = bytes.NewBuffer(postBodyStartGame)

	res, _ = http.Post(startGameUrl, "application/json", requestBody)

	return
}
