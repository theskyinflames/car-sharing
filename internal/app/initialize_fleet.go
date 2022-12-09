package app

import (
	"context"

	"theskyinflames/car-sharing/internal/domain"
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
func (ch InitializeFleet) Handle(ctx context.Context, cmd Command) error {
	co, ok := cmd.(InitializeFleetCmd)
	if !ok {
		return NewInvalidCommandError(InitializeFleetName, cmd.Name())
	}

	if err := ch.gr.RemoveAll(ctx); err != nil {
		return err
	}
	if err := ch.evr.RemoveAll(ctx); err != nil {
		return err
	}

	return ch.evr.AddAll(ctx, co.Cars)
}
