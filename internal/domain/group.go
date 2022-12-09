package domain

import "errors"

// Group is an entity
type Group struct {
	id     int
	people int
	car    *Car
}

// ErrWrongSize is self-described
var ErrWrongSize = errors.New("wrong size, it has to be from 1 to 6")

// NewGroup is a constructor
func NewGroup(id, people int) (Group, error) {
	if people < 1 || people > 6 {
		return Group{}, ErrWrongSize
	}
	return Group{id: id, people: people}, nil
}

// ID is a getter
func (g Group) ID() int {
	return g.id
}

// People is a getter
func (g Group) People() int {
	return g.people
}

// Ev is a getter
func (g Group) Ev() *Car {
	return g.car
}

// Hydrate hydrates a group
func (g *Group) Hydrate(id, people int, car *Car) {
	g.id = id
	g.people = people
	g.car = car
}

// GetOn links a group to its EV
func (g *Group) GetOn(ev *Car) {
	g.car = ev
}

// IsOnJourney returns TRUE is the group is in a journey
func (g Group) IsOnJourney() bool {
	return g.car != nil
}
