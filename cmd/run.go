package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"theskyinflames/car-sharing/internal/app"
	"theskyinflames/car-sharing/internal/infra/api"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

// Run Starts the API server
func Run(ctx context.Context, srvPort string) {
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"POST", "PUT", "POST"},
	})
	r.Use(cors.Handler)
	r.Use(middleware.Logger)

	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log := log.New(os.Stdout, "car-sharing: ", os.O_APPEND)

	bus := app.BuildBus(log)

	r.Put("/cars", api.InitializeFleet(bus))
	r.Post("/journey", api.Journey(bus))
	r.Post("/dropoff", api.DropOff(bus))
	r.Post("/locate", api.Locate(bus))

	fmt.Printf("serving at port %s\n", srvPort)
	if err := http.ListenAndServe(srvPort, r); err != nil {
		fmt.Printf("something went wrong trying to start the server: %s\n", err.Error())
	}
}
