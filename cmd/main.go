package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"theskyinflames/car-sharing/internal/app"
	"theskyinflames/car-sharing/internal/infra/api"
	"theskyinflames/car-sharing/internal/infra/repository"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
	"github.com/theskyinflames/cqrs-eda/pkg/helpers"

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

	log := log.New(os.Stdout, "car-sharing: ", os.O_APPEND)

	bus := buildBus(log)

	r.Put("/cars", api.InitializeFleet(bus))
	r.Post("/journey", api.Journey(bus))
	r.Post("/dropoff", api.DropOff(bus))
	r.Post("/locale", api.Locate(bus))

	fmt.Printf("serving at port %s\n", srvPort)
	if err := http.ListenAndServe(srvPort, r); err != nil {
		fmt.Printf("something went wrong trying to start the server: %s\n", err.Error())
	}
}

func buildBus(log cqrs.Logger) bus.Bus {
	gr := repository.NewGroupsRepository()
	evr := repository.NewCarRepository()

	initializeFleetCh := cqrs.ChErrMw(log)(app.NewInitializeFleet(&gr, &evr))
	journeyCh := cqrs.ChErrMw(log)(app.NewJourney(&gr, &evr))
	dropOffCh := cqrs.ChErrMw(log)(app.NewDropOff(&gr, &evr))
	localeQh := cqrs.QhErrMw(log)(app.NewLocate(&gr, &evr))

	bus := bus.New()
	bus.Register(app.InitializeFleetName, helpers.BusChHandler(initializeFleetCh))
	bus.Register(app.JourneyName, helpers.BusChHandler(journeyCh))
	bus.Register(app.DropOffName, helpers.BusChHandler(dropOffCh))
	bus.Register(app.LocateName, helpers.BusQhHandler(localeQh))
	return bus
}
