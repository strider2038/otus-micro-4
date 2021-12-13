package main

import (
	"fmt"
	"log"
	"net/http"

	"echo-service/internal/di"
)

func main() {
	router, err := di.NewRouter()
	if err != nil {
		log.Fatal("failed to create router: ", err)
	}

	address := fmt.Sprintf(":8000")
	log.Println("starting HTTP server at", address)
	log.Fatal(http.ListenAndServe(address, router))
}
