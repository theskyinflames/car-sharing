package domain_test

import (
	"testing"

	"theskyinflames/car-sharing/internal/domain"
	"theskyinflames/car-sharing/internal/fixtures"
	"theskyinflames/car-sharing/internal/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewFleet(t *testing.T) {
	var (
		carID1 = uuid.New()
		carID2 = uuid.New()
		carID3 = uuid.New()

		gID1 = uuid.New()
	)
	t.Run(`Given an unordered array of evs, when it's called, then a fleet is returned`, func(t *testing.T) {
		cars := []domain.Car{
			fixtures.Car{ID: helpers.UUIDPtr(carID1), Capacity: helpers.CarCapacityPtr(domain.CarCapacity5)}.Build(),
			fixtures.Car{ID: helpers.UUIDPtr(carID2), Capacity: helpers.CarCapacityPtr(domain.CarCapacity4)}.Build(),
			fixtures.Car{ID: helpers.UUIDPtr(carID3), Capacity: helpers.CarCapacityPtr(domain.CarCapacity6)}.Build(),
		}
		wg := []domain.Group{
			fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(3)}.Build(),
			fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(1)}.Build(),
			fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(2)}.Build(),
		}
		fleet := domain.NewFleet(cars, wg)
		require.Equal(t, []domain.Car{
			fixtures.Car{ID: helpers.UUIDPtr(carID3), Capacity: helpers.CarCapacityPtr(domain.CarCapacity6)}.Build(),
			fixtures.Car{ID: helpers.UUIDPtr(carID1), Capacity: helpers.CarCapacityPtr(domain.CarCapacity5)}.Build(),
			fixtures.Car{ID: helpers.UUIDPtr(carID2), Capacity: helpers.CarCapacityPtr(domain.CarCapacity4)}.Build(),
		}, fleet.Cars())
		require.Equal(t, []domain.Group{
			fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(1)}.Build(),
			fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(2)}.Build(),
			fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(3)}.Build(),
		}, fleet.WaitingGroups())
	})
}

func TestFleetJourney(t *testing.T) {
	var (
		carID = uuid.New()
		gID   = uuid.New()
	)
	testCases := []struct {
		name              string
		f                 domain.Fleet
		g                 domain.Group
		expectedOnJourney bool
		expectedCarID     uuid.UUID
	}{
		{
			name: `Given a group that does not fit to any car, 
				when it's called, 
				then it keeps waiting`,
			f: fixtures.Fleet{
				Cars: []domain.Car{
					fixtures.Car{
						ID:       helpers.UUIDPtr(carID),
						Capacity: helpers.CarCapacityPtr(domain.CarCapacity4),
						Journeys: domain.Journeys{
							gID: fixtures.Group{People: helpers.IntPtr(4)}.Build(),
						},
					}.Build(),
				},
			}.Build(),
			g:                 fixtures.Group{}.Build(),
			expectedOnJourney: false,
		},
		{
			name: `Given a group that does fit to a car, 
				when it's called, 
				then the group on route`,
			f: fixtures.Fleet{
				Cars: []domain.Car{
					fixtures.Car{
						ID:       helpers.UUIDPtr(carID),
						Capacity: helpers.CarCapacityPtr(domain.CarCapacity4),
					}.Build(),
				},
			}.Build(),
			g:                 fixtures.Group{People: helpers.IntPtr(2)}.Build(),
			expectedOnJourney: true,
			expectedCarID:     carID,
		},
	}

	for _, tc := range testCases {
		g, car := tc.f.Journey(tc.g)
		require.Equal(t, tc.expectedOnJourney, g.IsOnJourney(), t.Name())
		require.Equal(t, tc.expectedCarID, car.ID(), t.Name())
	}
}

func TestFleetRebuildWaitingGroupsList(t *testing.T) {
	var (
		gID1 = uuid.New()
		gID2 = uuid.New()
		gID3 = uuid.New()
		gID4 = uuid.New()
	)
	testCases := []struct {
		name                    string
		fleet                   domain.Fleet
		car                     domain.Car
		expectedOnJourney       map[int]domain.Group
		expectedWaitingGroups   []domain.Group
		expectedCarAvailability int
	}{
		{
			name: `Given a list of waiting groups that does not fit to the freed car, when it's called, no groups pass from waiting to journey`,
			fleet: fixtures.Fleet{
				WaitingGroups: []domain.Group{
					fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(6)}.Build(),
				},
			}.Build(),
			car: fixtures.Car{Capacity: helpers.CarCapacityPtr(domain.CarCapacity4)}.Build(),
			expectedWaitingGroups: []domain.Group{
				fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(6)}.Build(),
			},
			expectedCarAvailability: 4,
		},
		{
			name: `Given a list of waiting groups that does fit to any car, when it's called, some groups pass from waiting to journey`,
			fleet: fixtures.Fleet{
				WaitingGroups: []domain.Group{
					fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(2)}.Build(),
					fixtures.Group{ID: helpers.UUIDPtr(gID2), People: helpers.IntPtr(3)}.Build(),
					fixtures.Group{ID: helpers.UUIDPtr(gID3), People: helpers.IntPtr(1)}.Build(),
					fixtures.Group{ID: helpers.UUIDPtr(gID4), People: helpers.IntPtr(6)}.Build(),
				},
			}.Build(),
			car: fixtures.Car{Capacity: helpers.CarCapacityPtr(domain.CarCapacity6)}.Build(),
			expectedOnJourney: map[int]domain.Group{
				1: fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(2)}.Build(),
				2: fixtures.Group{ID: helpers.UUIDPtr(gID2), People: helpers.IntPtr(3)}.Build(),
				3: fixtures.Group{ID: helpers.UUIDPtr(gID3), People: helpers.IntPtr(1)}.Build(),
			},
			expectedWaitingGroups: []domain.Group{
				fixtures.Group{ID: helpers.UUIDPtr(gID4), People: helpers.IntPtr(6)}.Build(),
			},
			expectedCarAvailability: 0,
		},
	}

	for _, tc := range testCases {
		onJourney, err := tc.fleet.RebuildWaitingGroupsList(&tc.car)
		require.NoError(t, err, t.Name())
		require.Equal(t, len(tc.expectedOnJourney), len(onJourney), t.Name())
		require.Equal(t, tc.expectedWaitingGroups, tc.fleet.WaitingGroups(), t.Name())
		require.Equal(t, tc.expectedCarAvailability, tc.car.Availability())
	}
}

func TestDropOff(t *testing.T) {
	var (
		gID1 = uuid.New()
		gID2 = uuid.New()
		gID3 = uuid.New()
		gID4 = uuid.New()
		gID5 = uuid.New()
	)
	t.Run(`Given a group that does not exist, when it's called, then it returns an error`, func(t *testing.T) {
		var (
			fleet           = fixtures.Fleet{}.Build()
			groupToDropOff  = fixtures.Group{}.Build()
			expectedErrFunc = func(t *testing.T, err error) {
				require.ErrorIs(t, err, domain.ErrNotFound)
			}
		)
		_, _, err := fleet.DropOff(groupToDropOff, nil)
		require.Equal(t, expectedErrFunc == nil, err == nil)
		if err != nil {
			expectedErrFunc(t, err)
		}
	})

	t.Run(`Given a group that is in waiting state,
	when it's called,
	then it's removed from waiting list and no error is returned`, func(t *testing.T) {
		var (
			fleet = fixtures.Fleet{
				Cars: []domain.Car{},
				WaitingGroups: []domain.Group{
					fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(2)}.Build(),
					fixtures.Group{ID: helpers.UUIDPtr(gID2), People: helpers.IntPtr(6)}.Build(),
				},
			}.Build()
			groupToDropOff        = fixtures.Group{ID: helpers.UUIDPtr(gID2), People: helpers.IntPtr(6)}.Build()
			expectedWaitingGroups = []domain.Group{
				fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(2)}.Build(),
			}
		)

		_, _, err := fleet.DropOff(groupToDropOff, nil)
		require.NoError(t, err)
		require.Equal(t, expectedWaitingGroups, fleet.WaitingGroups())
	})

	t.Run(`Given a group that is on journey state and no other waiting group can fit its car when it's dropped off,
	when it's called,
	then it's removed from the car and no error is returned`, func(t *testing.T) {
		g, _, car := dropOffSetup()
		var (
			fleet = fixtures.Fleet{
				WaitingGroups: []domain.Group{
					fixtures.Group{ID: helpers.UUIDPtr(gID3), People: helpers.IntPtr(6)}.Build(),
				},
				Cars: []domain.Car{car},
			}.Build()
			groupToDropOff        = g
			expectedWaitingGroups = []domain.Group{
				fixtures.Group{ID: helpers.UUIDPtr(gID3), People: helpers.IntPtr(6)}.Build(),
			}
			expectedEvs = []domain.Car{car}
		)
		resultEv, onJourney, err := fleet.DropOff(groupToDropOff, &car)
		require.NoError(t, err)
		require.Empty(t, onJourney)
		require.Len(t, resultEv.Journeys(), 1)
		require.Equal(t, expectedWaitingGroups, fleet.WaitingGroups(), t.Name())
		require.Equal(t, len(expectedEvs), len(fleet.Cars()), t.Name())
	})

	t.Run(`Given a group that is on journey state and another waiting group can fit its car when it's dropped off,
	when it's called,
	then it's removed from the ev, the other group is got on the car, and no error is returned`, func(t *testing.T) {
		g, g2, car := dropOffSetup()
		var (
			fleet = fixtures.Fleet{
				WaitingGroups: []domain.Group{
					fixtures.Group{ID: helpers.UUIDPtr(gID3), People: helpers.IntPtr(2)}.Build(),
					fixtures.Group{ID: helpers.UUIDPtr(gID4), People: helpers.IntPtr(2)}.Build(),
					fixtures.Group{ID: helpers.UUIDPtr(gID5), People: helpers.IntPtr(6)}.Build(),
				},
				Cars: []domain.Car{
					fixtures.Car{
						ID:       helpers.UUIDPtr(uuid.New()),
						Capacity: helpers.CarCapacityPtr(domain.CarCapacity6),
						Journeys: domain.Journeys{
							g.ID():  g,
							g2.ID(): g2,
						},
					}.Build(),
				},
			}.Build()
			groupToDropOff    = g
			expectedOnJourney = map[uuid.UUID]domain.Group{
				gID3: fixtures.Group{ID: helpers.UUIDPtr(gID3), People: helpers.IntPtr(2)}.Build(),
				gID4: fixtures.Group{ID: helpers.UUIDPtr(gID4), People: helpers.IntPtr(2)}.Build(),
			}
			expectedWaitingGroups = []domain.Group{
				fixtures.Group{ID: helpers.UUIDPtr(gID5), People: helpers.IntPtr(6)}.Build(),
			}
			expectedCars = []domain.Car{
				fixtures.Car{
					ID:       helpers.UUIDPtr(uuid.New()),
					Capacity: helpers.CarCapacityPtr(domain.CarCapacity4),
					Journeys: domain.Journeys{
						gID3: fixtures.Group{ID: helpers.UUIDPtr(gID3), People: helpers.IntPtr(2)}.Build(),
						gID4: fixtures.Group{ID: helpers.UUIDPtr(gID4), People: helpers.IntPtr(2)}.Build(),
					},
				}.Build(),
			}
		)

		resultCar, onJourney, err := fleet.DropOff(groupToDropOff, &car)
		require.NoError(t, err)
		require.Len(t, resultCar.Journeys(), 3)
		require.Equal(t, len(expectedOnJourney), len(onJourney), t.Name())
		require.Equal(t, expectedWaitingGroups, fleet.WaitingGroups(), t.Name())
		require.Equal(t, len(expectedCars), len(fleet.Cars()), t.Name())
	})
}

func dropOffSetup() (domain.Group, domain.Group, domain.Car) {
	var (
		gID1 = uuid.New()
		gID2 = uuid.New()
	)
	var (
		g   = fixtures.Group{ID: helpers.UUIDPtr(gID1), People: helpers.IntPtr(2)}.Build()
		g2  = fixtures.Group{ID: helpers.UUIDPtr(gID2), People: helpers.IntPtr(2)}.Build()
		car = fixtures.Car{
			ID:       helpers.UUIDPtr(uuid.New()),
			Capacity: helpers.CarCapacityPtr(domain.CarCapacity6),
			Journeys: domain.Journeys{
				g.ID():  g,
				g2.ID(): g2,
			},
		}.Build()
	)
	g.GetOn(&car)
	return g, g2, car
}
