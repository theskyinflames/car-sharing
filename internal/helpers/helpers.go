package helpers

import (
	"theskyinflames/car-sharing/internal/domain"

	"github.com/google/uuid"
)

// IntPtr is a helper
func IntPtr(i int) *int {
	return &i
}

// CarCapacityPtr is a helper
func CarCapacityPtr(c domain.CarCapacity) *domain.CarCapacity {
	return &c
}

// BoolPtr is a helper
func BoolPtr(b bool) *bool {
	return &b
}

// EvPtr is a helper
func EvPtr(ev domain.Car) *domain.Car {
	return &ev
}

// UUIDPtr is a helper
func UUIDPtr(uuid uuid.UUID) *uuid.UUID {
	return &uuid
}
