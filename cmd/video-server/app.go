package main

import (
	"log"
	"net/http"

	routes "github.com/personal-work/video_server/internal/delivery"
)

func main() {

	routes.Init()

	// serve
	log.Printf("Server up on port 9001")
	log.Fatal(http.ListenAndServe(":9001", nil))
}
