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

func greetPlayer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("player has joined")
	io.WriteString(w, "Welcome to BringTen!\n")
}

func main() {

	mux := http.NewServeMux() // Routes request to approriate handler
	mux.HandleFunc("/", greetPlayer)

	err := http.ListenAndServe(":8080", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
