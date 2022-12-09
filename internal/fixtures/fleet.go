package fixtures

import (
	"theskyinflames/car-sharing/internal/domain"
	"theskyinflames/car-sharing/internal/helpers"
)

// Fleet is a fixture
type Fleet struct {
	Evs           []domain.Car
	WaitingGroups []domain.Group
}

// Build is self-described
func (f Fleet) Build() domain.Fleet {
	evs := []domain.Car{
		Car{ID: helpers.IntPtr(1), Capacity: helpers.CarCapacityPtr(domain.CarCapacity4)}.Build(),
	}
	if f.Evs != nil {
		evs = f.Evs
	}
	var wg []domain.Group
	if f.WaitingGroups != nil {
		wg = f.WaitingGroups
	}

	var df domain.Fleet
	df.Hydrate(evs, wg)
	return df
}
