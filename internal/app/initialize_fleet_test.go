package app_test

import (
	"context"
	"errors"
	"testing"

	"theskyinflames/car-sharing/internal/app"
	"theskyinflames/car-sharing/internal/domain"

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
				Cars: []app.Car{
					{ID: 1, Seats: domain.CarCapacity5},
					{ID: 2, Seats: domain.CarCapacity4},
					{ID: 3, Seats: domain.CarCapacity6},
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

		cars := tc.cmd.(app.InitializeFleetCmd).Cars
		require.Equal(t, len(tc.evr.AddAllCalls()[0].Evs), len(cars))
	}
}
