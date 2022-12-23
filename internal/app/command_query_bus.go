package app

import (
	"theskyinflames/car-sharing/internal/infra/repository"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
	"github.com/theskyinflames/cqrs-eda/pkg/helpers"
)

// BuildCommandQueryBus returns the command/query bus
func BuildCommandQueryBus(log cqrs.Logger, eventsBus bus.Bus) bus.Bus {
	gr := repository.NewGroupsRepository()
	evr := repository.NewCarRepository()

	chMw := cqrs.CommandHandlerMultiMiddleware(
		cqrs.ChEventMw(eventsBus),
		cqrs.ChErrMw(log),
	)

	initializeFleetCh := chMw(NewInitializeFleet(&gr, &evr))
	journeyCh := chMw(NewJourney(&gr, &evr))
	dropOffCh := chMw(NewDropOff(&gr, &evr))

	localeQh := cqrs.QhErrMw(log)(NewLocate(&gr, &evr))

	bus := bus.New()
	bus.Register(InitializeFleetName, helpers.BusChHandler(initializeFleetCh))
	bus.Register(JourneyName, helpers.BusChHandler(journeyCh))
	bus.Register(DropOffName, helpers.BusChHandler(dropOffCh))
	bus.Register(LocateName, helpers.BusQhHandler(localeQh))
	return bus
}
