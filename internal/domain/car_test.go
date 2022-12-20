package domain_test

import (
	"testing"

	"theskyinflames/car-sharing/internal/domain"
	"theskyinflames/car-sharing/internal/fixtures"
	"theskyinflames/car-sharing/internal/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewCar(t *testing.T) {
	t.Run(`Given an ID and a car capacity, when it's called, then an Car is returned`, func(t *testing.T) {
		var (
			capacity = domain.CarCapacity5
			id       = uuid.New()
			car      = domain.NewCar(id, capacity)
		)
		require.Equal(t, id, car.ID())
		require.Equal(t, capacity, car.Capacity())
		require.Equal(t, capacity.Int(), car.Availability())
	})
}

func TestCarAvailability(t *testing.T) {
	var (
		id1 = uuid.New()
		id2 = uuid.New()
	)
	t.Run(`Given a Car, when it's called, then its availability is returned`, func(t *testing.T) {
		car := fixtures.Car{
			ID:       helpers.UUIDPtr(uuid.New()),
			Capacity: helpers.CarCapacityPtr(domain.CarCapacity6),
			Journeys: domain.Journeys{
				id1: fixtures.Group{ID: helpers.UUIDPtr(id1), People: helpers.IntPtr(2)}.Build(),
				id2: fixtures.Group{ID: helpers.UUIDPtr(id2), People: helpers.IntPtr(3)}.Build(),
			},
		}.Build()
		require.Equal(t, 1, car.Availability())
	})
}

func TestCarGetOn(t *testing.T) {
	var (
		gID1 = uuid.New()
		gID2 = uuid.New()
	)
	toGetOn := fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(3)}.Build()
	testCases := []struct {
		name                 string
		car                  domain.Car
		g                    domain.Group
		expectedAvailability int
		expectedErrFunc      func(*testing.T, error)
	}{
		{
			name: `Given a Car,
				when a group whose size exceeds its availability tries to get on it,
				then an error is returned`,
			car: fixtures.Car{
				Capacity: helpers.CarCapacityPtr(domain.CarCapacity6),
				Journeys: domain.Journeys{
					gID1: fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(2)}.Build(),
					gID2: fixtures.Group{ID: helpers.UUIDPtr(gID2), People: helpers.IntPtr(3)}.Build(),
				},
			}.Build(),
			g: toGetOn,
			expectedErrFunc: func(t *testing.T, err error) {
				require.Error(t, err, domain.ErrNotFit)
			},
		},
		{
			name: `Given a Car,
				when a group whose size does not exceeds its availability tries to get on it,
				then group is added and no error is returned`,
			car: fixtures.Car{
				Capacity: helpers.CarCapacityPtr(domain.CarCapacity6),
				Journeys: domain.Journeys{
					gID1: fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(2)}.Build(),
				},
			}.Build(),
			g:                    fixtures.Group{ID: helpers.UUIDPtr(gID2), People: helpers.IntPtr(3)}.Build(),
			expectedAvailability: 1,
		},
	}

	for _, tc := range testCases {
		err := tc.car.GetOn(tc.g)
		require.Equal(t, tc.expectedErrFunc == nil, err == nil, tc.name)
		if err != nil {
			tc.expectedErrFunc(t, err)
			continue
		}
		require.Equal(t, tc.expectedAvailability, tc.car.Availability(), tc.name)
		require.Len(t, tc.car.Journeys(), 2)
		for _, g := range map[uuid.UUID]domain.Group{
			gID1:         fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(2)}.Build(),
			toGetOn.ID(): toGetOn,
		} {
			_, ok := tc.car.Journeys()[g.ID()]
			require.True(t, ok)
		}
	}
}

func TestCarDropOff(t *testing.T) {
	var (
		gID1 = uuid.New()
		gID2 = uuid.New()
		gID3 = uuid.New()
	)
	testCases := []struct {
		name            string
		car             domain.Car
		gID             uuid.UUID
		expectedErrFunc func(*testing.T, error)
	}{
		{
			name: `Given a Car,
				when it's tried to drop off a unknown group,
				then an error is returned`,
			car: fixtures.Car{
				Capacity: helpers.CarCapacityPtr(domain.CarCapacity6),
				Journeys: domain.Journeys{
					gID1: fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(2)}.Build(),
					gID2: fixtures.Group{ID: helpers.UUIDPtr(gID2), People: helpers.IntPtr(3)}.Build(),
				},
			}.Build(),
			gID: gID3,
			expectedErrFunc: func(t *testing.T, err error) {
				require.Error(t, err, domain.ErrNotFound)
			},
		},
		{
			name: `Given a Car,
				when it's tried to drop off a known group,
				then group is dropped off and no error is returned`,
			car: fixtures.Car{
				Capacity: helpers.CarCapacityPtr(domain.CarCapacity6),
				Journeys: domain.Journeys{
					gID1: fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(2)}.Build(),
				},
			}.Build(),
			gID: gID1,
		},
	}

	for _, tc := range testCases {
		err := tc.car.DropOff(tc.gID)
		require.Equal(t, tc.expectedErrFunc == nil, err == nil, tc.name)
		if err != nil {
			tc.expectedErrFunc(t, err)
			continue
		}
		_, ok := tc.car.Journeys()[tc.gID]
		require.False(t, ok, tc.name)
	}
}
