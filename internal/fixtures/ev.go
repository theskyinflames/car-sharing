package fixtures

import (
	"math/rand"
	"time"

	"theskyinflames/car-sharing/internal/domain"
)

// Car is fixture
type Car struct {
	ID       *int
	Capacity *domain.CarCapacity
	Journeys domain.Journeys
}

// Build is self-described
func (e Car) Build() domain.Car {
	rand.Seed(time.Now().UnixNano())
	id := rand.Int()
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
