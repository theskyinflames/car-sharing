package main

import (
	"fmt"
	"net/http"

	"theskyinflames/car-sharing/internal/app"
	"theskyinflames/car-sharing/internal/infra/api"
	"theskyinflames/car-sharing/internal/infra/repository"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

const srvPort = ":80"

func main() {
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

	gr := repository.NewGroupsRepository()
	evr := repository.NewCarRepository()

	initializeFleetCh := app.CommandHandlerErrorWrapperMiddleware()(app.NewInitializeFleet(&gr, &evr))
	r.Put("/cars", api.InitializeFleet(initializeFleetCh))

	journeyCh := app.CommandHandlerErrorWrapperMiddleware()(app.NewJourney(&gr, &evr))
	r.Post("/journey", api.Journey(journeyCh))

	dropOffCh := app.CommandHandlerErrorWrapperMiddleware()(app.NewDropOff(&gr, &evr))
	r.Post("/dropoff", api.DropOff(dropOffCh))

	localeQh := app.QueryHandlerErrorWrapperMiddleware()(app.NewLocate(&gr, &evr))
	r.Post("/locale", api.Locate(localeQh))

	fmt.Printf("serving at port %s\n", srvPort)
	if err := http.ListenAndServe(srvPort, r); err != nil {
		fmt.Printf("something went wrong trying to start the server: %s\n", err.Error())
	}
}
