package app

import (
	"theskyinflames/car-sharing/internal/infra/repository"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
	"github.com/theskyinflames/cqrs-eda/pkg/helpers"
)

// BuildBus returns the command/query bus
func BuildBus(log cqrs.Logger) bus.Bus {
	gr := repository.NewGroupsRepository()
	evr := repository.NewCarRepository()

	initializeFleetCh := cqrs.ChErrMw(log)(NewInitializeFleet(&gr, &evr))
	journeyCh := cqrs.ChErrMw(log)(NewJourney(&gr, &evr))
	dropOffCh := cqrs.ChErrMw(log)(NewDropOff(&gr, &evr))
	localeQh := cqrs.QhErrMw(log)(NewLocate(&gr, &evr))

	bus := bus.New()
	bus.Register(InitializeFleetName, helpers.BusChHandler(initializeFleetCh))
	bus.Register(JourneyName, helpers.BusChHandler(journeyCh))
	bus.Register(DropOffName, helpers.BusChHandler(dropOffCh))
	bus.Register(LocateName, helpers.BusQhHandler(localeQh))
	return bus
}
