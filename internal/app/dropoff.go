package app

import (
	"context"

	"theskyinflames/car-sharing/internal/domain"
)

// DropOffCmd is a command
type DropOffCmd struct {
	GroupID int
}

// DropOffName is self-described
var DropOffName = "dropOff.group"

// Name implements the Command interface
func (cmd DropOffCmd) Name() string {
	return DropOffName
}

// DropOff is a command handler
type DropOff struct {
	gr  GroupsRepository
	evr CarsRepository
}

// NewDropOff is a constructor
func NewDropOff(gr GroupsRepository, evr CarsRepository) DropOff {
	return DropOff{gr: gr, evr: evr}
}

// Handle implements CommandHandler interface
func (ch DropOff) Handle(ctx context.Context, cmd Command) error {
	co, ok := cmd.(DropOffCmd)
	if !ok {
		return NewInvalidCommandError(DropOffName, cmd.Name())
	}

	g, err := ch.gr.FindByID(ctx, co.GroupID)
	if err != nil {
		return err
	}

	wg, err := ch.gr.FindGroupsWithoutCar(ctx)
	if err != nil {
		return err
	}

	var ev *domain.Car
	if g.Ev() != nil {
		gev, err := ch.evr.FindByID(ctx, g.Ev().ID())
		if err != nil {
			return err
		}
		ev = &gev
	}

	// here we don't need all the list of evs. We only need the dropping group ev
	fleet := domain.NewFleet(nil, wg)
	resultEv, onJourney, err := fleet.DropOff(g, ev)
	if err != nil {
		return err
	}

	if resultEv != nil {
		if err := ch.evr.Update(ctx, *resultEv); err != nil {
			return err
		}
	}

	for _, oj := range onJourney {
		if err := ch.gr.Update(ctx, oj); err != nil {
			return err
		}
	}

	return ch.gr.RemoveByID(ctx, g.ID())
}
