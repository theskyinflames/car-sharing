package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"theskyinflames/car-sharing/internal/app"
	"theskyinflames/car-sharing/internal/domain"
	"theskyinflames/car-sharing/internal/infra/repository"

	"github.com/google/uuid"
	"github.com/theskyinflames/cqrs-eda/pkg/bus"
)

// InitializeFleet is the HTTP handler to initialize the fleet
func InitializeFleet(commandBus bus.Bus) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkHeader(r, "Content-Type", "application/json") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var rq CarsRqJson
		if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var cars []app.Car
		for _, car := range rq {
			carID, err := uuid.Parse(car.Id)
			if err != nil {
				http.Error(w, "invalid car uuid", http.StatusBadRequest)
				return
			}
			seats, err := domain.ParseCarCapacityFromInt(int(car.Seats))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			cars = append(cars, app.Car{ID: carID, Seats: seats})
		}

		cmd := app.InitializeFleetCmd{Cars: cars}
		if _, err := commandBus.Dispatch(r.Context(), cmd); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Journey is the HTTP handler to add a new group
func Journey(commandBus bus.Bus) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkHeader(r, "Content-Type", "application/json") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var rq JourneyRqJson
		if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		gID, err := uuid.Parse(rq.Id)
		if err != nil { // this restriction is in the JSON schema, but the Go JSON schema compiler does not implements it yet.
			http.Error(w, "minimum id value is 1", http.StatusBadRequest)
			return
		}

		cmd := app.JourneyCmd{
			ID:     gID,
			People: int(rq.People),
		}
		if _, err := commandBus.Dispatch(r.Context(), cmd); err != nil {
			switch {
			case errors.Is(err, repository.ErrPKConflict):
				w.WriteHeader(http.StatusBadRequest)
				return
			case errors.Is(err, domain.ErrWrongSize):
				w.WriteHeader(http.StatusBadRequest)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

// DropOff is the HTTP handler to drop off a group
func DropOff(commandBus bus.Bus) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkHeader(r, "Content-Type", "application/x-www-form-urlencoded") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id := r.FormValue("ID")
		gID, err := uuid.Parse(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cmd := app.DropOffCmd{
			GroupID: gID,
		}

		if _, err := commandBus.Dispatch(r.Context(), cmd); err != nil {
			if errors.Is(err, domain.ErrNotFound) || errors.Is(err, repository.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Locate is the HTTP handler to locate a group
func Locate(queryBus bus.Bus) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkHeader(r, "Content-Type", "application/x-www-form-urlencoded") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id := r.FormValue("ID")
		gID, err := uuid.Parse(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cmd := app.LocateQuery{
			GroupID: gID,
		}

		queryRs, err := queryBus.Dispatch(r.Context(), cmd)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) || errors.Is(err, repository.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		locateRs := queryRs.(app.LocateResponse)
		if !locateRs.IsInJourney {
			w.Header().Add("Accept", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Accept", "application/json")
		jsonRs := LocateRsJson{
			Id:    locateRs.Car.ID().String(),
			Seats: LocateRsJsonSeats(locateRs.Car.Capacity()),
		}
		b, _ := json.Marshal(jsonRs)
		if _, err := w.Write(b); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func checkHeader(r *http.Request, name string, expected string) bool {
	v, ok := r.Header[name]
	if !ok {
		return false
	}
	return v[0] == expected
}
