package app

import (
	"context"

	"theskyinflames/car-sharing/internal/domain"

	"github.com/google/uuid"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
	"github.com/theskyinflames/cqrs-eda/pkg/events"
)

// JourneyCmd is a command
type JourneyCmd struct {
	ID     uuid.UUID
	People int
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
func (ch Journey) Handle(ctx context.Context, cmd cqrs.Command) ([]events.Event, error) {
	co, ok := cmd.(JourneyCmd)
	if !ok {
		return nil, NewInvalidCommandError(JourneyName, cmd.Name())
	}

	wg, err := ch.gr.FindGroupsWithoutCar(ctx)
	if err != nil {
		return nil, err
	}

	g, err := domain.NewGroup(co.ID, co.People)
	if err != nil {
		return nil, err
	}

	if err := ch.gr.Add(ctx, g); err != nil {
		return nil, err
	}

	evs, err := ch.evr.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	fleet := domain.NewFleet(evs, wg)
	g, ev := fleet.Journey(g) // try to get the group on a ev

	if !g.IsOnJourney() { // if the g is not in journey, there is not ev to be updated. Otherwise, its list of groups is updated
		return nil, nil
	}

	if err := ch.gr.Update(ctx, g); err != nil {
		return nil, err
	}

	if err := ch.evr.Update(ctx, ev); err != nil {
		return nil, err
	}

	return g.Events(), nil
}
