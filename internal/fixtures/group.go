package fixtures

import (
	"math/rand"
	"time"

	"theskyinflames/car-sharing/internal/domain"
)

// Group is a fixture
type Group struct {
	ID     *int
	People *int
	Car    *domain.Car
}

// Build is self-described
func (g Group) Build() domain.Group {
	rand.Seed(time.Now().UnixNano())
	id := rand.Int()
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
