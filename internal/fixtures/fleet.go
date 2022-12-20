package fixtures

import (
	"theskyinflames/car-sharing/internal/domain"
	"theskyinflames/car-sharing/internal/helpers"

	"github.com/google/uuid"
)

// Fleet is a fixture
type Fleet struct {
	Cars          []domain.Car
	WaitingGroups []domain.Group
}

// Build is self-described
func (f Fleet) Build() domain.Fleet {
	evs := []domain.Car{
		Car{ID: helpers.UUIDPtr(uuid.New()), Capacity: helpers.CarCapacityPtr(domain.CarCapacity4)}.Build(),
	}
	if f.Cars != nil {
		evs = f.Cars
	}
	var wg []domain.Group
	if f.WaitingGroups != nil {
		wg = f.WaitingGroups
	}

	var df domain.Fleet
	df.Hydrate(evs, wg)
	return df
}
