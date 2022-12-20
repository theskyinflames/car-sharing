package domain

import (
	"errors"
	"sort"

	"github.com/google/uuid"
)

// Fleet is a domain service
type Fleet struct {
	cars          []Car
	waitingGroups []Group
}

// NewFleet is a constructor
func NewFleet(evs []Car, waitingGroups []Group) Fleet {
	fleet := Fleet{cars: evs, waitingGroups: waitingGroups}
	fleet.sort()
	return fleet
}

func (f Fleet) sort() {
	sort.SliceStable(f.cars, func(i, j int) bool { // order descending by ev availability
		return f.cars[i].Availability() > f.cars[j].Availability()
	})
	sort.SliceStable(f.waitingGroups, func(i, j int) bool { // order ascending by group size
		return f.waitingGroups[i].People() < f.waitingGroups[j].People()
	})
}

// Cars is a getter
func (f Fleet) Cars() []Car {
	return f.cars
}

// WaitingGroups is a getter
func (f Fleet) WaitingGroups() []Group {
	return f.waitingGroups
}

// Hydrate hydrates a Fleet. Used for testing
func (f *Fleet) Hydrate(cars []Car, waitingGroups []Group) {
	f.cars = cars
	f.waitingGroups = waitingGroups
	f.sort()
}

// Journey adds a new group to the EV and out it on journey state
func (f Fleet) Journey(g Group) (Group, Car) {
	for _, car := range f.cars {
		if err := car.GetOn(g); err == nil {
			g.GetOn(&car)
			return g, car
		}
	}
	return g, Car{}
}

// DropOff removes a group from the waiting list, or from its ev if it's on journey
func (f *Fleet) DropOff(g Group, car *Car) (*Car, Journeys, error) {
	if car == nil {
		// look for the group in the waiting list
		if removed := f.removeGroupsFromWaitingList(map[uuid.UUID]Group{g.ID(): g}); removed == 0 {
			return nil, nil, ErrNotFound
		}
		return nil, nil, nil
	}

	// At ch update on Journey Groups
	if err := car.DropOff(g.ID()); err != nil {
		return nil, nil, err
	}

	newJourneys, err := f.RebuildWaitingGroupsList(car)
	if err != nil {
		return nil, nil, err
	}
	return car, newJourneys, nil
}

// RebuildWaitingGroupsList is self-described
func (f *Fleet) RebuildWaitingGroupsList(car *Car) (newJourneys Journeys, err error) { // At ch update on journey groups
	newJourneys = make(Journeys)

	// Try to use the availability of the car to fit in it so many groups as it's possible
	for _, wg := range f.waitingGroups {
		if err := car.GetOn(wg); err != nil {
			if errors.Is(err, ErrNotFit) {
				continue
			}
			return nil, err
		}
		wg.GetOn(car)
		newJourneys[wg.id] = wg

		if car.Availability() == 0 {
			break
		}
	}

	if len(newJourneys) > 0 {
		f.removeGroupsFromWaitingList(newJourneys)
	}

	return newJourneys, nil
}

func (f *Fleet) removeGroupsFromWaitingList(toRemove map[uuid.UUID]Group) int {
	var removed int
	tmp := make(map[uuid.UUID]Group)
	for _, g := range f.waitingGroups {
		tmp[g.ID()] = g
	}
	f.waitingGroups = make([]Group, 0)
	for _, g := range tmp {
		_, ok := toRemove[g.ID()]
		if !ok {
			f.waitingGroups = append(f.waitingGroups, g)
			continue
		}
		removed++
	}
	return removed
}
