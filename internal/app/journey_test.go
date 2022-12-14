package app_test

import (
	"context"
	"errors"
	"testing"

	"theskyinflames/car-sharing/internal/app"
	"theskyinflames/car-sharing/internal/domain"
	"theskyinflames/car-sharing/internal/fixtures"
	"theskyinflames/car-sharing/internal/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
)

func TestJourney(t *testing.T) {
	var (
		randomErr = errors.New("")

		jID1 = uuid.New()
		jID2 = uuid.New()
		jID3 = uuid.New()
		jID4 = uuid.New()
	)

	testCases := []struct {
		name              string
		cmd               cqrs.Command
		gr                *GroupsRepositoryMock
		cr                *CarsRepositoryMock
		expectedOnJourney bool
		expectedErrFunc   func(*testing.T, error)
	}{
		{
			name: `Given an invalid command, when it's called, then an error is returned`,
			cmd:  newInvalidCommand(),
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorAs(t, err, &app.InvalidCommandError{})
			},
		},
		{
			name: `Given an gr repository that returns an error on FindGroupsWithoutEv method, when it's called, then an error is returned`,
			cmd:  app.JourneyCmd{},
			gr: &GroupsRepositoryMock{
				FindGroupsWithoutCarFunc: func(_ context.Context) ([]domain.Group, error) {
					return nil, randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given an gr repository that returns an error on Add method, when it's called, then an error is returned`,
			cmd: app.JourneyCmd{
				ID:     jID1,
				People: 2,
			},
			gr: &GroupsRepositoryMock{
				AddFunc: func(_ context.Context, _ domain.Group) error {
					return randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given an evr repository that returns an error on FindAll method, when it's called, then an error is returned`,
			cmd: app.JourneyCmd{
				ID:     jID1,
				People: 2,
			},
			gr: &GroupsRepositoryMock{},
			cr: &CarsRepositoryMock{
				FindAllFunc: func(_ context.Context) ([]domain.Car, error) {
					return nil, randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given group that gets on a ev and a groups repository that returns an error on Update method, 
				when it's called, then no error is returned`,
			cmd: app.JourneyCmd{
				ID:     jID2,
				People: 1,
			},
			gr: &GroupsRepositoryMock{},
			cr: &CarsRepositoryMock{
				FindAllFunc: func(_ context.Context) ([]domain.Car, error) {
					return []domain.Car{
						fixtures.Car{
							ID:       helpers.UUIDPtr(uuid.New()),
							Capacity: helpers.CarCapacityPtr(domain.CarCapacity6),
							Journeys: domain.Journeys{},
						}.Build(),
					}, nil
				},
				UpdateFunc: func(_ context.Context, _ domain.Car) error {
					return randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given group that gets on a ev and a evs repository that returns an error on Update method, 
				when it's called, then no error is returned`,
			cmd: app.JourneyCmd{
				ID:     jID3,
				People: 1,
			},
			gr: &GroupsRepositoryMock{
				UpdateFunc: func(_ context.Context, _ domain.Group) error {
					return randomErr
				},
			},
			cr: &CarsRepositoryMock{
				FindAllFunc: func(_ context.Context) ([]domain.Car, error) {
					return []domain.Car{
						fixtures.Car{
							ID:       helpers.UUIDPtr(uuid.New()),
							Capacity: helpers.CarCapacityPtr(domain.CarCapacity6),
							Journeys: domain.Journeys{},
						}.Build(),
					}, nil
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given group that gets on a ev, when it's called, then no error is returned`,
			cmd: app.JourneyCmd{
				ID:     jID4,
				People: 1,
			},
			gr: &GroupsRepositoryMock{},
			cr: &CarsRepositoryMock{
				FindAllFunc: func(_ context.Context) ([]domain.Car, error) {
					return []domain.Car{
						fixtures.Car{
							ID:       helpers.UUIDPtr(uuid.New()),
							Capacity: helpers.CarCapacityPtr(domain.CarCapacity6),
							Journeys: domain.Journeys{},
						}.Build(),
					}, nil
				},
			},
			expectedOnJourney: true,
		},
		{
			name: `Given group that does not get on a ev, when it's called, then no error is returned`,
			cmd: app.JourneyCmd{
				ID:     jID1,
				People: 1,
			},
			gr:                &GroupsRepositoryMock{},
			cr:                &CarsRepositoryMock{},
			expectedOnJourney: false,
		},
	}

	for _, tc := range testCases {
		ch := app.NewJourney(tc.gr, tc.cr)
		_, err := ch.Handle(context.Background(), tc.cmd)
		require.Equal(t, tc.expectedErrFunc == nil, err == nil)
		if err != nil {
			tc.expectedErrFunc(t, err)
			continue
		}

		require.Len(t, tc.gr.FindGroupsWithoutCarCalls(), 1)
		require.Len(t, tc.gr.AddCalls(), 1)
		require.Equal(t, tc.gr.AddCalls()[0].G.ID(), tc.cmd.(app.JourneyCmd).ID)
		require.Equal(t, tc.gr.AddCalls()[0].G.People(), tc.cmd.(app.JourneyCmd).People)
		require.Len(t, tc.cr.FindAllCalls(), 1)
		if tc.expectedOnJourney {
			require.Len(t, tc.gr.UpdateCalls(), 1)
			require.Len(t, tc.cr.UpdateCalls(), 1)
		}
	}
}
