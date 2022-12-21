package app

import (
	"context"

	"theskyinflames/car-sharing/internal/domain"

	"github.com/google/uuid"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
	"github.com/theskyinflames/cqrs-eda/pkg/events"
)

// Car is a DTO
type Car struct {
	ID    uuid.UUID
	Seats domain.CarCapacity
}

// InitializeFleetCmd is a Command
type InitializeFleetCmd struct {
	Cars []Car
}

// InitializeFleetName is self-described
var InitializeFleetName = "initialize.fleet"

// Name implements the Command interface
func (cmd InitializeFleetCmd) Name() string {
	return InitializeFleetName
}

// InitializeFleet is a command handler
type InitializeFleet struct {
	gr  GroupsRepository
	evr CarsRepository
}

// NewInitializeFleet is constructor
func NewInitializeFleet(gr GroupsRepository, evr CarsRepository) InitializeFleet {
	return InitializeFleet{gr: gr, evr: evr}
}

// Handle implements the CommandHandler constructor
func (ch InitializeFleet) Handle(ctx context.Context, cmd cqrs.Command) ([]events.Event, error) {
	co, ok := cmd.(InitializeFleetCmd)
	if !ok {
		return nil, NewInvalidCommandError(InitializeFleetName, cmd.Name())
	}

	var cars []domain.Car
	for _, car := range co.Cars {
		seats, err := domain.ParseCarCapacityFromInt(int(car.Seats))
		if err != nil {
			return nil, err
		}
		cars = append(cars, domain.NewCar(car.ID, seats))
	}
	if err := ch.gr.RemoveAll(ctx); err != nil {
		return nil, err
	}
	if err := ch.evr.RemoveAll(ctx); err != nil {
		return nil, err
	}
	if err := ch.evr.AddAll(ctx, cars); err != nil {
		return nil, err
	}

	var evs []events.Event
	for _, car := range cars {
		evs = append(evs, car.Events()...)
	}
	return evs, nil
}
