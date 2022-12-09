package app

import (
	"context"

	"theskyinflames/car-sharing/internal/domain"
)

//go:generate moq -stub -out zmock_app_repositories_test.go -pkg app_test . GroupsRepository CarsRepository

// GroupsRepository is self-described
type GroupsRepository interface {
	RemoveAll(ctx context.Context) error
	Add(ctx context.Context, g domain.Group) error
	Update(ctx context.Context, g domain.Group) error
	FindGroupsWithoutCar(ctx context.Context) ([]domain.Group, error)
	FindByID(ctx context.Context, ID int) (domain.Group, error)
	RemoveByID(ctx context.Context, ID int) error
}

// CarsRepository is self-described
type CarsRepository interface {
	RemoveAll(ctx context.Context) error
	Update(ctx context.Context, ev domain.Car) error
	AddAll(ctx context.Context, evs []domain.Car) error
	FindAll(ctx context.Context) ([]domain.Car, error)
	FindByID(ctx context.Context, ID int) (domain.Car, error)
}
