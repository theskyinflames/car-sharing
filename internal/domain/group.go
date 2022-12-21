package domain

import (
	"errors"

	"github.com/google/uuid"
	"github.com/theskyinflames/cqrs-eda/pkg/ddd"
)

// Group is an entity
type Group struct {
	ddd.AggregateBasic

	people int
	car    *Car
}

// ErrWrongSize is self-described
var ErrWrongSize = errors.New("wrong size, it has to be from 1 to 6")

// NewGroup is a constructor
func NewGroup(id uuid.UUID, people int) (Group, error) {
	if people < 1 || people > 6 {
		return Group{}, ErrWrongSize
	}
	return Group{AggregateBasic: ddd.NewAggregateBasic(id), people: people}, nil
}

// ID is a getter
func (g Group) ID() uuid.UUID {
	return g.AggregateBasic.ID()
}

// People is a getter
func (g Group) People() int {
	return g.people
}

// Car is a getter
func (g Group) Car() *Car {
	return g.car
}

// Hydrate hydrates a group
func (g *Group) Hydrate(id uuid.UUID, people int, car *Car) {
	g.AggregateBasic = ddd.NewAggregateBasic(id)
	g.people = people
	g.car = car
}

// GetOn links a group to its EV
func (g *Group) GetOn(car *Car) {
	g.car = car

	g.RecordEvent(NewGroupSetOnJourneyEvent(*g))
}

// DropOff drops off the group from its car
func (g *Group) DropOff() {
	g.car = nil

	g.RecordEvent(NewGroupDroppedOff(*g))
}

// IsOnJourney returns TRUE is the group is in a journey
func (g Group) IsOnJourney() bool {
	return g.car != nil
}
