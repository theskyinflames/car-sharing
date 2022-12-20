package repository

// By the moment, using an in-memory repository without transactions support.
// If I have time enough, I'll add Sqlite3, with upper.io as ORM, with transactions support, as an CH middleware

import (
	"context"
	"errors"
	"sync"

	"theskyinflames/car-sharing/internal/domain"

	"github.com/google/uuid"
)

// CarRepository is a repository
type CarRepository struct {
	cars map[uuid.UUID]domain.Car

	mux *sync.RWMutex
}

// NewCarRepository is a constructor
func NewCarRepository() CarRepository {
	return CarRepository{cars: make(map[uuid.UUID]domain.Car), mux: &sync.RWMutex{}}
}

// RemoveAll is a constructor
func (cr *CarRepository) RemoveAll(_ context.Context) error {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	cr.cars = make(map[uuid.UUID]domain.Car)
	return nil
}

// ErrPKConflict is self-described
var ErrPKConflict = errors.New("pk conflict")

// AddAll is self-described
func (cr CarRepository) AddAll(_ context.Context, evs []domain.Car) error {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	for _, ev := range evs {
		_, ok := cr.cars[ev.ID()]
		if ok {
			return ErrPKConflict
		}
		cr.cars[ev.ID()] = ev
	}
	return nil
}

// ErrNotFound is self-described
var ErrNotFound = errors.New("not found")

// Update is self-described
func (cr CarRepository) Update(_ context.Context, ev domain.Car) error {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	if _, ok := cr.cars[ev.ID()]; !ok {
		return ErrNotFound
	}
	cr.cars[ev.ID()] = ev
	return nil
}

// FindAll is self-described
func (cr CarRepository) FindAll(_ context.Context) ([]domain.Car, error) {
	cr.mux.RLock()
	defer cr.mux.RUnlock()

	evs := make([]domain.Car, 0)
	for _, ev := range cr.cars {
		evs = append(evs, ev)
	}
	return evs, nil
}

// FindByID is a finder
func (cr CarRepository) FindByID(_ context.Context, id uuid.UUID) (domain.Car, error) {
	cr.mux.RLock()
	defer cr.mux.RUnlock()

	for _, ev := range cr.cars {
		if ev.ID() == id {
			return ev, nil
		}
	}
	return domain.Car{}, ErrNotFound
}

// GroupsRepository is a repository
type GroupsRepository struct {
	groups map[uuid.UUID]domain.Group

	mux *sync.RWMutex
}

// NewGroupsRepository is a constructor
func NewGroupsRepository() GroupsRepository {
	return GroupsRepository{groups: make(map[uuid.UUID]domain.Group), mux: &sync.RWMutex{}}
}

// RemoveAll is self-described
func (gr *GroupsRepository) RemoveAll(_ context.Context) error {
	gr.mux.Lock()
	defer gr.mux.Unlock()

	gr.groups = make(map[uuid.UUID]domain.Group)
	return nil
}

// FindGroupsWithoutCar is a finder
func (gr GroupsRepository) FindGroupsWithoutCar(_ context.Context) ([]domain.Group, error) {
	gr.mux.RLock()
	defer gr.mux.RUnlock()

	var withoutEv []domain.Group
	for _, g := range gr.groups {
		if g.Car() == nil {
			withoutEv = append(withoutEv, g)
		}
	}
	return withoutEv, nil
}

// Update is self-described
func (gr GroupsRepository) Update(_ context.Context, g domain.Group) error {
	gr.mux.Lock()
	defer gr.mux.Unlock()

	if _, ok := gr.groups[g.ID()]; !ok {
		return ErrNotFound
	}
	gr.groups[g.ID()] = g
	return nil
}

// Add is self-described
func (gr GroupsRepository) Add(_ context.Context, g domain.Group) error {
	gr.mux.Lock()
	defer gr.mux.Unlock()

	_, ok := gr.groups[g.ID()]
	if ok {
		return ErrPKConflict
	}
	gr.groups[g.ID()] = g
	return nil
}

// FindByID is a finder
func (gr GroupsRepository) FindByID(_ context.Context, id uuid.UUID) (domain.Group, error) {
	gr.mux.RLock()
	defer gr.mux.RUnlock()

	for _, g := range gr.groups {
		if g.ID() == id {
			return g, nil
		}
	}
	return domain.Group{}, ErrNotFound
}

// RemoveByID is self-described
func (gr GroupsRepository) RemoveByID(_ context.Context, id uuid.UUID) error {
	gr.mux.Lock()
	defer gr.mux.Unlock()

	delete(gr.groups, id)
	return nil
}
