package app

import (
	"context"

	"theskyinflames/car-sharing/internal/domain"

	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
)

// InitializeFleetCmd is a Command
type InitializeFleetCmd struct {
	Cars []domain.Car
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
func (ch InitializeFleet) Handle(ctx context.Context, cmd cqrs.Command) ([]cqrs.Event, error) {
	co, ok := cmd.(InitializeFleetCmd)
	if !ok {
		return nil, NewInvalidCommandError(InitializeFleetName, cmd.Name())
	}

	if err := ch.gr.RemoveAll(ctx); err != nil {
		return nil, err
	}
	if err := ch.evr.RemoveAll(ctx); err != nil {
		return nil, err
	}

	return nil, ch.evr.AddAll(ctx, co.Cars)
}
