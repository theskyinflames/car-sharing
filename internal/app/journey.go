package app

import (
	"context"

	"theskyinflames/car-sharing/internal/domain"
)

// JourneyCmd is a command
type JourneyCmd struct {
	Group domain.Group
}

// JourneyName is self-described
var JourneyName = "journey"

// Name implements the Command interface
func (cmd JourneyCmd) Name() string {
	return JourneyName
}

// Journey is a command handler
type Journey struct {
	gr  GroupsRepository
	evr CarsRepository
}

// NewJourney is a constructor
func NewJourney(gr GroupsRepository, evr CarsRepository) Journey {
	return Journey{gr: gr, evr: evr}
}

// Handle implements CommandHandler interface
func (ch Journey) Handle(ctx context.Context, cmd Command) error {
	co, ok := cmd.(JourneyCmd)
	if !ok {
		return NewInvalidCommandError(JourneyName, cmd.Name())
	}

	wg, err := ch.gr.FindGroupsWithoutCar(ctx)
	if err != nil {
		return err
	}

	if err := ch.gr.Add(ctx, co.Group); err != nil {
		return err
	}

	evs, err := ch.evr.FindAll(ctx)
	if err != nil {
		return err
	}

	fleet := domain.NewFleet(evs, wg)
	g, ev := fleet.Journey(co.Group) // try to get the group on a ev

	if !g.IsOnJourney() { // if the g is not in journey, there is not ev to be updated. Otherwise, its list of groups is updated
		return nil
	}

	if err := ch.gr.Update(ctx, g); err != nil {
		return err
	}

	return ch.evr.Update(ctx, ev)
}
