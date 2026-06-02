package main

import (
	"log"
	"net/http"

	"github.com/Andhika-GIT/concurrent-cinema-booking/internal/adapter/redis"
	"github.com/Andhika-GIT/concurrent-cinema-booking/internal/api"
	"github.com/Andhika-GIT/concurrent-cinema-booking/internal/booking"
	"github.com/Andhika-GIT/concurrent-cinema-booking/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	store := booking.NewRedisStore(redis.NewClient("localhost:63795"))
	svc := booking.NewService(store)

	handler := handler.NewBookingHandler(svc)

	api.SetupRoutes(mux, handler)

	err := http.ListenAndServe(":8089", mux)

	if err != nil {
		log.Fatal(err)
	}
}
