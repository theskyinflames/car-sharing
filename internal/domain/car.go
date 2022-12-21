package domain

import (
	"errors"

	"github.com/google/uuid"
	"github.com/theskyinflames/cqrs-eda/pkg/ddd"
)

// CarCapacity is the number of seats that en EV has
type CarCapacity int

// Int is self-described
func (evc CarCapacity) Int() int {
	return int(evc)
}

// Allowed EV capacities
const (
	CarCapacity4 CarCapacity = 4
	CarCapacity5 CarCapacity = 5
	CarCapacity6 CarCapacity = 6
)

// ErrCapacityNotSupported is self-described
var ErrCapacityNotSupported = errors.New("capacity not supported")

// ParseCarCapacityFromInt returns the int value for an EvCapacity
func ParseCarCapacityFromInt(c int) (CarCapacity, error) {
	switch c {
	case 4:
		return CarCapacity4, nil
	case 5:
		return CarCapacity5, nil
	case 6:
		return CarCapacity6, nil
	default:
		return 0, ErrCapacityNotSupported
	}
}

// Journeys is the set of groups that are on journey in the Ev
type Journeys map[uuid.UUID]Group

// Car is an entity
type Car struct {
	ddd.AggregateBasic

	capacity CarCapacity
	journeys Journeys
}

// NewCar is a constructor
func NewCar(id uuid.UUID, capacity CarCapacity) Car {
	car := Car{
		AggregateBasic: ddd.NewAggregateBasic(id),
		capacity:       capacity,
		journeys:       make(Journeys),
	}

	car.AggregateBasic.RecordEvent(NewCarCreatedEvent(car))

	return car
}

// ID is a getter
func (e Car) ID() uuid.UUID {
	return e.AggregateBasic.ID()
}

// Capacity is a getter
func (e Car) Capacity() CarCapacity {
	return e.capacity
}

// Journeys is a getter
func (e Car) Journeys() Journeys {
	return e.journeys
}

// Hydrate hydrates an EV
func (e *Car) Hydrate(id uuid.UUID, capacity CarCapacity, journeys Journeys) {
	e.AggregateBasic = ddd.NewAggregateBasic(id)
	e.capacity = capacity
	e.journeys = journeys
}

// Availability returns the amount of available seats
func (e Car) Availability() int {
	var currentLoad int
	for _, g := range e.journeys {
		currentLoad += g.people
	}
	return e.capacity.Int() - currentLoad
}

// ErrNotFit is self-described
var ErrNotFit = errors.New("not fit")

// GetOn is self-described
func (e *Car) GetOn(g Group) error {
	if e.Availability() < g.people {
		return ErrNotFit
	}
	e.journeys[g.ID()] = g
	return nil
}

// ErrNotFound is self-described
var ErrNotFound = errors.New("not found")

// DropOff drops off an on journey group from the EV
func (e *Car) DropOff(id uuid.UUID) error {
	_, ok := e.journeys[id]
	if !ok {
		return ErrNotFound
	}
	delete(e.journeys, id)
	return nil
}
