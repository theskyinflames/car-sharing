package domain

import (
	"errors"
	"sort"
)

// Fleet is a domain service
type Fleet struct {
	evs           []Car
	waitingGroups []Group
}

// NewFleet is a constructor
func NewFleet(evs []Car, waitingGroups []Group) Fleet {
	fleet := Fleet{evs: evs, waitingGroups: waitingGroups}
	fleet.sort()
	return fleet
}

func (f Fleet) sort() {
	sort.SliceStable(f.evs, func(i, j int) bool { // order descending by ev availability
		return f.evs[i].Availability() > f.evs[j].Availability()
	})
	sort.SliceStable(f.waitingGroups, func(i, j int) bool { // order ascending by group size
		return f.waitingGroups[i].People() < f.waitingGroups[j].People()
	})
}

// Evs is a getter
func (f Fleet) Evs() []Car {
	return f.evs
}

// WaitingGroups is a getter
func (f Fleet) WaitingGroups() []Group {
	return f.waitingGroups
}

// Hydrate hydrates a Fleet. Used for testing
func (f *Fleet) Hydrate(evs []Car, waitingGroups []Group) {
	f.evs = evs
	f.waitingGroups = waitingGroups
	f.sort()
}

// Journey adds a new group to the EV and out it on journey state
func (f Fleet) Journey(g Group) (Group, Car) {
	for _, ev := range f.evs {
		if err := ev.GetOn(g); err == nil {
			g.GetOn(&ev)
			return g, ev
		}
	}
	return g, Car{}
}

// DropOff removes a group from the waiting list, or from its ev if it's on journey
func (f *Fleet) DropOff(g Group, ev *Car) (*Car, Journeys, error) {
	if ev == nil {
		// look for the group in the waiting list
		if removed := f.removeGroupsFromWaitingList(map[int]Group{g.ID(): g}); removed == 0 {
			return nil, nil, ErrNotFound
		}
		return nil, nil, nil
	}

	// At ch update on Journey Groups
	if err := ev.DropOff(g.ID()); err != nil {
		return nil, nil, err
	}

	newJourneys, err := f.RebuildWaitingGroupsList(ev)
	if err != nil {
		return nil, nil, err
	}
	return ev, newJourneys, nil
}

// RebuildWaitingGroupsList is self-described
func (f *Fleet) RebuildWaitingGroupsList(ev *Car) (newJourneys Journeys, err error) { // At ch update on journey groups
	newJourneys = make(Journeys)

	// Try to use the availability of the car to fit in it so many groups as it's possible
	for _, wg := range f.waitingGroups {
		if err := ev.GetOn(wg); err != nil {
			if errors.Is(err, ErrNotFit) {
				continue
			}
			return nil, err
		}
		wg.GetOn(ev)
		newJourneys[wg.id] = wg

		if ev.Availability() == 0 {
			break
		}
	}

	if len(newJourneys) > 0 {
		f.removeGroupsFromWaitingList(newJourneys)
	}

	return newJourneys, nil
}

func (f *Fleet) removeGroupsFromWaitingList(toRemove map[int]Group) int {
	var removed int
	tmp := make(map[int]Group)
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
