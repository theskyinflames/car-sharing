package fixtures

import (
	"theskyinflames/car-sharing/internal/domain"

	"github.com/google/uuid"
)

// Car is fixture
type Car struct {
	ID       *uuid.UUID
	Capacity *domain.CarCapacity
	Journeys domain.Journeys
}

// Build is self-described
func (e Car) Build() domain.Car {
	id := uuid.New()
	if e.ID != nil {
		id = *e.ID
	}
	capacity := domain.CarCapacity4
	if e.Capacity != nil {
		capacity = *e.Capacity
	}
	journeys := make(domain.Journeys)
	if e.Journeys != nil {
		journeys = e.Journeys
	}
	dev := domain.Car{}
	dev.Hydrate(id, capacity, journeys)
	return dev
}
