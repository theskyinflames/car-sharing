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

func TestDropOff(t *testing.T) {
	randomErr := errors.New("")
	testCases := []struct {
		name            string
		cmd             cqrs.Command
		gr              *GroupsRepositoryMock
		cr              *CarsRepositoryMock
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
			name: `Given groups repository that returns an error on FindById method,
				when it's called, then an error is returned`,
			cmd: app.DropOffCmd{},
			gr: &GroupsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ int) (domain.Group, error) {
					return domain.Group{}, randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given groups repository that returns an error on FindGroupsWithoutEv method,
				when it's called, then an error is returned`,
			cmd: app.DropOffCmd{},
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
			name: `Given evs repository that returns an error on FindById method,
				when it's called, then an error is returned`,
			cmd: app.DropOffCmd{},
			gr: &GroupsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ int) (domain.Group, error) {
					return fixtures.Group{
						Car: helpers.EvPtr(fixtures.Car{}.Build()),
					}.Build(), nil
				},
			},
			cr: &CarsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ int) (domain.Car, error) {
					return domain.Car{}, randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given evs repository that returns an error on Update method,
				when it's called, then an error is returned`,
			cmd: app.DropOffCmd{},
			gr: &GroupsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ int) (domain.Group, error) {
					return fixtures.Group{
						ID:  helpers.IntPtr(1),
						Car: helpers.EvPtr(fixtures.Car{}.Build()),
					}.Build(), nil
				},
			},
			cr: &CarsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ int) (domain.Car, error) {
					return fixtures.Car{
						Journeys: domain.Journeys{
							1: fixtures.Group{
								ID: helpers.IntPtr(1),
							}.Build(),
						},
					}.Build(), nil
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
			name: `Given groups repository that returns an error on Update method,
				when it's called, then an error is returned`,
			cmd: app.DropOffCmd{},
			gr: &GroupsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ int) (domain.Group, error) {
					return fixtures.Group{
						ID:  helpers.IntPtr(1),
						Car: helpers.EvPtr(fixtures.Car{}.Build()),
					}.Build(), nil
				},
				FindGroupsWithoutCarFunc: func(_ context.Context) ([]domain.Group, error) {
					return []domain.Group{
						fixtures.Group{ID: helpers.IntPtr(2), People: helpers.IntPtr(1)}.Build(),
					}, nil
				},
				UpdateFunc: func(_ context.Context, _ domain.Group) error {
					return randomErr
				},
			},
			cr: &CarsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ int) (domain.Car, error) {
					return fixtures.Car{
						Journeys: domain.Journeys{
							1: fixtures.Group{
								ID: helpers.IntPtr(1),
							}.Build(),
						},
					}.Build(), nil
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given a group to be dropped off, 
				when it's called, then no error is returned`,
			cmd: app.DropOffCmd{},
			gr: &GroupsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ int) (domain.Group, error) {
					return fixtures.Group{
						ID:  helpers.IntPtr(1),
						Car: helpers.EvPtr(fixtures.Car{}.Build()),
					}.Build(), nil
				},
				FindGroupsWithoutCarFunc: func(_ context.Context) ([]domain.Group, error) {
					return []domain.Group{
						fixtures.Group{ID: helpers.IntPtr(2), People: helpers.IntPtr(1)}.Build(),
					}, nil
				},
				UpdateFunc: func(_ context.Context, _ domain.Group) error {
					return nil
				},
			},
			cr: &CarsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ int) (domain.Car, error) {
					return fixtures.Car{
						Journeys: domain.Journeys{
							1: fixtures.Group{
								ID: helpers.IntPtr(1),
							}.Build(),
						},
					}.Build(), nil
				},
			},
		},
	}

	for _, tc := range testCases {
		ch := app.NewDropOff(tc.gr, tc.cr)
		_, err := ch.Handle(context.Background(), tc.cmd)
		require.Equal(t, tc.expectedErrFunc == nil, err == nil)
		if err != nil {
			tc.expectedErrFunc(t, err)
			continue
		}

		require.Len(t, tc.gr.FindByIDCalls(), 1)
		require.Len(t, tc.gr.FindGroupsWithoutCarCalls(), 1)
		require.Len(t, tc.cr.FindByIDCalls(), 1)
		require.Len(t, tc.cr.UpdateCalls(), 1)
		require.Len(t, tc.cr.UpdateCalls()[0].Ev.Journeys(), 1)
		require.Equal(t, 2, tc.cr.UpdateCalls()[0].Ev.Journeys()[2].ID())
		require.Len(t, tc.gr.UpdateCalls(), 1)
		require.Equal(t, 2, tc.gr.UpdateCalls()[0].G.ID())
		require.Len(t, tc.gr.RemoveByIDCalls(), 1)
	}
}
