package domain_test

import (
	"testing"

	"theskyinflames/car-sharing/internal/domain"
	"theskyinflames/car-sharing/internal/fixtures"
	"theskyinflames/car-sharing/internal/helpers"

	"github.com/stretchr/testify/require"
)

func TestNewFleet(t *testing.T) {
	t.Run(`Given an unordered array of evs, when it's called, then a fleet is returned`, func(t *testing.T) {
		evs := []domain.Car{
			fixtures.Car{ID: helpers.IntPtr(1), Capacity: helpers.CarCapacityPtr(domain.CarCapacity5)}.Build(),
			fixtures.Car{ID: helpers.IntPtr(2), Capacity: helpers.CarCapacityPtr(domain.CarCapacity4)}.Build(),
			fixtures.Car{ID: helpers.IntPtr(3), Capacity: helpers.CarCapacityPtr(domain.CarCapacity6)}.Build(),
		}
		wg := []domain.Group{
			fixtures.Group{ID: helpers.IntPtr(1), People: helpers.IntPtr(3)}.Build(),
			fixtures.Group{ID: helpers.IntPtr(1), People: helpers.IntPtr(1)}.Build(),
			fixtures.Group{ID: helpers.IntPtr(1), People: helpers.IntPtr(2)}.Build(),
		}
		fleet := domain.NewFleet(evs, wg)
		require.Equal(t, []domain.Car{
			fixtures.Car{ID: helpers.IntPtr(3), Capacity: helpers.CarCapacityPtr(domain.CarCapacity6)}.Build(),
			fixtures.Car{ID: helpers.IntPtr(1), Capacity: helpers.CarCapacityPtr(domain.CarCapacity5)}.Build(),
			fixtures.Car{ID: helpers.IntPtr(2), Capacity: helpers.CarCapacityPtr(domain.CarCapacity4)}.Build(),
		}, fleet.Evs())
		require.Equal(t, []domain.Group{
			fixtures.Group{ID: helpers.IntPtr(1), People: helpers.IntPtr(1)}.Build(),
			fixtures.Group{ID: helpers.IntPtr(1), People: helpers.IntPtr(2)}.Build(),
			fixtures.Group{ID: helpers.IntPtr(1), People: helpers.IntPtr(3)}.Build(),
		}, fleet.WaitingGroups())
	})
}

func TestFleetJourney(t *testing.T) {
	testCases := []struct {
		name              string
		f                 domain.Fleet
		g                 domain.Group
		expectedOnJourney bool
		expectedEvID      int
	}{
		{
			name: `Given a group that does not fit to any ev, 
				when it's called, 
				then it keeps waiting`,
			f: fixtures.Fleet{
				Evs: []domain.Car{
					fixtures.Car{
						ID:       helpers.IntPtr(1),
						Capacity: helpers.CarCapacityPtr(domain.CarCapacity4),
						Journeys: domain.Journeys{
							1: fixtures.Group{People: helpers.IntPtr(4)}.Build(),
						},
					}.Build(),
				},
			}.Build(),
			g:                 fixtures.Group{}.Build(),
			expectedOnJourney: false,
		},
		{
			name: `Given a group that does fit to a ev, 
				when it's called, 
				then the group on route`,
			f: fixtures.Fleet{
				Evs: []domain.Car{
					fixtures.Car{
						ID:       helpers.IntPtr(1),
						Capacity: helpers.CarCapacityPtr(domain.CarCapacity4),
					}.Build(),
				},
			}.Build(),
			g:                 fixtures.Group{People: helpers.IntPtr(2)}.Build(),
			expectedOnJourney: true,
			expectedEvID:      1,
		},
	}

	for _, tc := range testCases {
		g, ev := tc.f.Journey(tc.g)
		require.Equal(t, tc.expectedOnJourney, g.IsOnJourney(), t.Name())
		require.Equal(t, tc.expectedEvID, ev.ID(), t.Name())
	}
}

func TestFleetRebuildWaitingGroupsList(t *testing.T) {
	testCases := []struct {
		name                        string
		fleet                       domain.Fleet
		freedEv                     domain.Car
		expectedOnJourney           map[int]domain.Group
		expectedWaitingGroups       []domain.Group
		freedEvExpectedAvailability int
	}{
		{
			name: `Given a list of waiting groups that does not fit to the freed ev, when it's called, no groups pass from waiting to journey`,
			fleet: fixtures.Fleet{
				WaitingGroups: []domain.Group{
					fixtures.Group{ID: helpers.IntPtr(1), People: helpers.IntPtr(6)}.Build(),
				},
			}.Build(),
			freedEv: fixtures.Car{Capacity: helpers.CarCapacityPtr(domain.CarCapacity4)}.Build(),
			expectedWaitingGroups: []domain.Group{
				fixtures.Group{ID: helpers.IntPtr(1), People: helpers.IntPtr(6)}.Build(),
			},
			freedEvExpectedAvailability: 4,
		},
		{
			name: `Given a list of waiting groups that does fit to any ev, when it's called, some groups pass from waiting to journey`,
			fleet: fixtures.Fleet{
				WaitingGroups: []domain.Group{
					fixtures.Group{ID: helpers.IntPtr(1), People: helpers.IntPtr(2)}.Build(),
					fixtures.Group{ID: helpers.IntPtr(2), People: helpers.IntPtr(3)}.Build(),
					fixtures.Group{ID: helpers.IntPtr(3), People: helpers.IntPtr(1)}.Build(),
					fixtures.Group{ID: helpers.IntPtr(4), People: helpers.IntPtr(6)}.Build(),
				},
			}.Build(),
			freedEv: fixtures.Car{Capacity: helpers.CarCapacityPtr(domain.CarCapacity6)}.Build(),
			expectedOnJourney: map[int]domain.Group{
				1: fixtures.Group{ID: helpers.IntPtr(1), People: helpers.IntPtr(2)}.Build(),
				2: fixtures.Group{ID: helpers.IntPtr(2), People: helpers.IntPtr(3)}.Build(),
				3: fixtures.Group{ID: helpers.IntPtr(3), People: helpers.IntPtr(1)}.Build(),
			},
			expectedWaitingGroups: []domain.Group{
				fixtures.Group{ID: helpers.IntPtr(4), People: helpers.IntPtr(6)}.Build(),
			},
			freedEvExpectedAvailability: 0,
		},
	}

	for _, tc := range testCases {
		onJourney, err := tc.fleet.RebuildWaitingGroupsList(&tc.freedEv)
		require.NoError(t, err, t.Name())
		require.Equal(t, len(tc.expectedOnJourney), len(onJourney), t.Name())
		require.Equal(t, tc.expectedWaitingGroups, tc.fleet.WaitingGroups(), t.Name())
		require.Equal(t, tc.freedEvExpectedAvailability, tc.freedEv.Availability())
	}
}

func TestDropOff(t *testing.T) {
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
				Evs: []domain.Car{},
				WaitingGroups: []domain.Group{
					fixtures.Group{ID: helpers.IntPtr(1), People: helpers.IntPtr(2)}.Build(),
					fixtures.Group{ID: helpers.IntPtr(2), People: helpers.IntPtr(6)}.Build(),
				},
			}.Build()
			groupToDropOff        = fixtures.Group{ID: helpers.IntPtr(2), People: helpers.IntPtr(6)}.Build()
			expectedWaitingGroups = []domain.Group{
				fixtures.Group{ID: helpers.IntPtr(1), People: helpers.IntPtr(2)}.Build(),
			}
		)

		_, _, err := fleet.DropOff(groupToDropOff, nil)
		require.NoError(t, err)
		require.Equal(t, expectedWaitingGroups, fleet.WaitingGroups())
	})

	t.Run(`Given a group that is on journey state and no other waiting group can fit its Ev when it's dropped off,
	when it's called,
	then it's removed from the ev and no error is returned`, func(t *testing.T) {
		g, _, ev := dropOffSetup()
		var (
			fleet = fixtures.Fleet{
				WaitingGroups: []domain.Group{
					fixtures.Group{ID: helpers.IntPtr(3), People: helpers.IntPtr(6)}.Build(),
				},
				Evs: []domain.Car{ev},
			}.Build()
			groupToDropOff        = g
			expectedWaitingGroups = []domain.Group{
				fixtures.Group{ID: helpers.IntPtr(3), People: helpers.IntPtr(6)}.Build(),
			}
			expectedEvs = []domain.Car{ev}
		)
		resultEv, onJourney, err := fleet.DropOff(groupToDropOff, &ev)
		require.NoError(t, err)
		require.Empty(t, onJourney)
		require.Len(t, resultEv.Journeys(), 1)
		require.Equal(t, expectedWaitingGroups, fleet.WaitingGroups(), t.Name())
		require.Equal(t, len(expectedEvs), len(fleet.Evs()), t.Name())
	})

	t.Run(`Given a group that is on journey state and a waiting group can fit its Ev when it's dropped off,
	when it's called,
	then it's removed from the ev, the other group is got on the ev, and no error is returned`, func(t *testing.T) {
		g, g2, ev := dropOffSetup()
		var (
			fleet = fixtures.Fleet{
				WaitingGroups: []domain.Group{
					fixtures.Group{ID: helpers.IntPtr(3), People: helpers.IntPtr(2)}.Build(),
					fixtures.Group{ID: helpers.IntPtr(4), People: helpers.IntPtr(2)}.Build(),
					fixtures.Group{ID: helpers.IntPtr(5), People: helpers.IntPtr(6)}.Build(),
				},
				Evs: []domain.Car{
					fixtures.Car{
						ID:       helpers.IntPtr(2),
						Capacity: helpers.CarCapacityPtr(domain.CarCapacity6),
						Journeys: domain.Journeys{
							g.ID():  g,
							g2.ID(): g2,
						},
					}.Build(),
				},
			}.Build()
			groupToDropOff    = g
			expectedOnJourney = map[int]domain.Group{
				3: fixtures.Group{ID: helpers.IntPtr(3), People: helpers.IntPtr(2)}.Build(),
				4: fixtures.Group{ID: helpers.IntPtr(4), People: helpers.IntPtr(2)}.Build(),
			}
			expectedWaitingGroups = []domain.Group{
				fixtures.Group{ID: helpers.IntPtr(5), People: helpers.IntPtr(6)}.Build(),
			}
			expectedEvs = []domain.Car{
				fixtures.Car{
					ID:       helpers.IntPtr(2),
					Capacity: helpers.CarCapacityPtr(domain.CarCapacity4),
					Journeys: domain.Journeys{
						3: fixtures.Group{ID: helpers.IntPtr(3), People: helpers.IntPtr(2)}.Build(),
						4: fixtures.Group{ID: helpers.IntPtr(4), People: helpers.IntPtr(2)}.Build(),
					},
				}.Build(),
			}
		)

		resultEv, onJourney, err := fleet.DropOff(groupToDropOff, &ev)
		require.NoError(t, err)
		require.Len(t, resultEv.Journeys(), 3)
		require.Equal(t, len(expectedOnJourney), len(onJourney), t.Name())
		require.Equal(t, expectedOnJourney[0].ID(), onJourney[0].ID())
		require.Equal(t, expectedWaitingGroups, fleet.WaitingGroups(), t.Name())
		require.Equal(t, len(expectedEvs), len(fleet.Evs()), t.Name())
	})
}

func dropOffSetup() (domain.Group, domain.Group, domain.Car) {
	var (
		g  = fixtures.Group{ID: helpers.IntPtr(1), People: helpers.IntPtr(2)}.Build()
		g2 = fixtures.Group{ID: helpers.IntPtr(2), People: helpers.IntPtr(2)}.Build()
		ev = fixtures.Car{
			ID:       helpers.IntPtr(2),
			Capacity: helpers.CarCapacityPtr(domain.CarCapacity6),
			Journeys: domain.Journeys{
				g.ID():  g,
				g2.ID(): g2,
			},
		}.Build()
	)
	g.GetOn(&ev)
	return g, g2, ev
}
