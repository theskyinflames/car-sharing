package fixtures

import (
	"theskyinflames/car-sharing/internal/domain"

	"github.com/google/uuid"
)

// Group is a fixture
type Group struct {
	ID     *uuid.UUID
	People *int
	Car    *domain.Car
}

// Build is self-described
func (g Group) Build() domain.Group {
	id := uuid.New()
	if g.ID != nil {
		id = *g.ID
	}
	people := 1
	if g.People != nil {
		people = *g.People
	}
	var car *domain.Car
	if g.Car != nil {
		car = g.Car
	}
	dg := domain.Group{}
	dg.Hydrate(id, people, car)
	return dg
}
