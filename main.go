package main

import (
	"fmt"
	"net/http"
	"testjunior/handlers"
)

func main() {

	http.HandleFunc("/receive", handlers.RecieveHandler)

	http.HandleFunc("/refresh", handlers.RefreshHandler)

	fmt.Println("Server is listening...")
	err := http.ListenAndServe("localhost:8181", nil)
	if err != nil {
		return
	}
}
