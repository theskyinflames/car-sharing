package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"theskyinflames/car-sharing/internal/app"
	"theskyinflames/car-sharing/internal/domain"
	"theskyinflames/car-sharing/internal/infra/repository"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"
)

// InitializeFleet is the HTTP handler to initialize the fleet
func InitializeFleet(commandBus bus.Bus) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var rq InitializeFleetRqJson
		if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var evs []domain.Car
		for _, ev := range rq {
			if ev.Id < 1 { // this restriction is in the JSON schema, but the Go JSON schema compiler does not implements it yet.
				http.Error(w, "minimum id value is 1", http.StatusBadRequest)
				return
			}
			capacity, err := domain.ParseCarCapacityFromInt(int(ev.Seats))
			if err != nil {
				if errors.Is(err, domain.ErrCapacityNotSupported) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
			}
			evs = append(evs, domain.NewCar(ev.Id, capacity))
		}

		cmd := app.InitializeFleetCmd{Cars: evs}
		if _, err := commandBus.Dispatch(r.Context(), cmd); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Journey is the HTTP handler to add a new group
func Journey(commandBus bus.Bus) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var rq JourneyRqJson
		if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if rq.Id < 1 { // this restriction is in the JSON schema, but the Go JSON schema compiler does not implements it yet.
			http.Error(w, "minimum id value is 1", http.StatusBadRequest)
			return
		}

		g, err := domain.NewGroup(rq.Id, int(rq.People))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cmd := app.JourneyCmd{
			Group: g,
		}
		if _, err := commandBus.Dispatch(r.Context(), cmd); err != nil {
			if errors.Is(err, repository.ErrPKConflict) {
				w.WriteHeader(http.StatusConflict)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// DropOff is the HTTP handler to drop off a group
func DropOff(commandBus bus.Bus) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var rq DropOffRqJson
		if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if rq.Id < 1 { // this restriction is in the JSON schema, but the Go JSON schema compiler does not implements it yet.
			http.Error(w, "minimum id value is 1", http.StatusBadRequest)
			return
		}

		cmd := app.DropOffCmd{
			GroupID: rq.Id,
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
		var rq DropOffRqJson
		if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if rq.Id < 1 { // this restriction is in the JSON schema, but the Go JSON schema compiler does not implements it yet.
			http.Error(w, "minimum id value is 1", http.StatusBadRequest)
			return
		}

		cmd := app.LocateQuery{
			GroupID: rq.Id,
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
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonRs := LocateRsJson{
			Id:    locateRs.Ev.ID(),
			Seats: LocateRsJsonSeats(locateRs.Ev.Capacity()),
		}
		b, _ := json.Marshal(jsonRs)
		if _, err := w.Write(b); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}
}
