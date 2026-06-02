package main

import (
	"log"
	"net/http"

	"github.com/Andhika-GIT/concurrent-cinema-booking/internal/api"
)

func main() {
	mux := http.NewServeMux()

	api.SetupRoutes(mux)

	err := http.ListenAndServe(":8089", mux)

	if err != nil {
		log.Fatal(err)
	}
}
