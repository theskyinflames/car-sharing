package app_test

import (
	"context"
	"errors"
	"testing"

	"theskyinflames/car-sharing/internal/app"
	"theskyinflames/car-sharing/internal/domain"
	"theskyinflames/car-sharing/internal/fixtures"
	"theskyinflames/car-sharing/internal/helpers"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
)

func newInvalidCommand() cqrs.Command {
	return &CommandMock{
		NameFunc: func() string {
			return "invalid_command"
		},
	}
}

func TestInitializeFleet(t *testing.T) {
	randomErr := errors.New("")

	testCases := []struct {
		name            string
		cmd             cqrs.Command
		gr              *GroupsRepositoryMock
		evr             *CarsRepositoryMock
		expectedErrFunc func(*testing.T, error)
	}{
		{
			name: `Given an invalid command, when it's called, then an error is returned`,
			cmd:  newInvalidCommand(),
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorAs(t, err, &app.InvalidCommandError{})
			},
		},
		{
			name: `Given a groups repository that returns an error on RemoveAll method,
				when it's called,
				then an error is returned`,
			cmd: app.InitializeFleetCmd{},
			gr: &GroupsRepositoryMock{
				RemoveAllFunc: func(_ context.Context) error {
					return randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given a evs repository that returns an error on RemoveAll method,
				when it's called,
				then an error is returned`,
			cmd: app.InitializeFleetCmd{},
			gr:  &GroupsRepositoryMock{},
			evr: &CarsRepositoryMock{
				RemoveAllFunc: func(_ context.Context) error {
					return randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given a evs repository that returns an error on AddAll method,
				when it's called,
				then an error is returned`,
			cmd: app.InitializeFleetCmd{},
			gr:  &GroupsRepositoryMock{},
			evr: &CarsRepositoryMock{
				AddAllFunc: func(_ context.Context, _ []domain.Car) error {
					return randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given a list of Evs to be added,
				when it's called,
				then no error is returned`,
			cmd: app.InitializeFleetCmd{
				Cars: []domain.Car{
					fixtures.Car{ID: helpers.IntPtr(1), Capacity: helpers.CarCapacityPtr(domain.CarCapacity5)}.Build(),
					fixtures.Car{ID: helpers.IntPtr(2), Capacity: helpers.CarCapacityPtr(domain.CarCapacity4)}.Build(),
					fixtures.Car{ID: helpers.IntPtr(3), Capacity: helpers.CarCapacityPtr(domain.CarCapacity6)}.Build(),
				},
			},
			gr:  &GroupsRepositoryMock{},
			evr: &CarsRepositoryMock{},
		},
	}

	for _, tc := range testCases {
		ch := app.NewInitializeFleet(tc.gr, tc.evr)
		_, err := ch.Handle(context.Background(), tc.cmd)
		require.Equal(t, tc.expectedErrFunc == nil, err == nil)
		if err != nil {
			tc.expectedErrFunc(t, err)
			continue
		}

		require.Len(t, tc.gr.RemoveAllCalls(), 1)
		require.Len(t, tc.evr.RemoveAllCalls(), 1)
		require.Equal(t, tc.cmd.(app.InitializeFleetCmd).Cars, tc.evr.AddAllCalls()[0].Evs)
	}
}
