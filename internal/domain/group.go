package domain

import (
	"errors"

	"github.com/google/uuid"
)

// Group is an entity
type Group struct {
	id     uuid.UUID
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
	return Group{id: id, people: people}, nil
}

// ID is a getter
func (g Group) ID() uuid.UUID {
	return g.id
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
	g.id = id
	g.people = people
	g.car = car
}

// GetOn links a group to its EV
func (g *Group) GetOn(car *Car) {
	g.car = car
}

// IsOnJourney returns TRUE is the group is in a journey
func (g Group) IsOnJourney() bool {
	return g.car != nil
}
