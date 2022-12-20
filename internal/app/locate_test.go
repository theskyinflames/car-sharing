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

func newInvalidQuery() cqrs.Query {
	return &QueryMock{
		NameFunc: func() string {
			return "invalid_query"
		},
	}
}

func TestLocate(t *testing.T) {
	randomErr := errors.New("")
	testCases := []struct {
		name               string
		q                  cqrs.Query
		gr                 *GroupsRepositoryMock
		evr                *CarsRepositoryMock
		expectedCallsToEvr int
		expectedRs         app.LocateResponse
		expectedErrFunc    func(*testing.T, error)
	}{
		{
			name: `Given an invalid query, when it's called, then an error is returned`,
			q:    newInvalidQuery(),
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorAs(t, err, &app.InvalidQueryError{})
			},
		},
		{
			name: `Given a groups repository that returns an error on FindById method,
				when it's called, then an error is returned`,
			q: app.LocateQuery{},
			gr: &GroupsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ uuid.UUID) (domain.Group, error) {
					return domain.Group{}, randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given a evs repository that returns an error on FindById method, 
				when it's called, then an error is returned`,
			q: app.LocateQuery{},
			gr: &GroupsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ uuid.UUID) (domain.Group, error) {
					return fixtures.Group{
						Car: helpers.EvPtr(
							fixtures.Car{}.Build(),
						),
					}.Build(), nil
				},
			},
			evr: &CarsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ uuid.UUID) (domain.Car, error) {
					return domain.Car{}, randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given an waiting group, 
				when it's called, then it's returned`,
			q: app.LocateQuery{},
			gr: &GroupsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ uuid.UUID) (domain.Group, error) {
					return fixtures.Group{}.Build(), nil
				},
			},
			expectedRs: app.LocateResponse{},
		},
		{
			name: `Given an on journey,
				when it's called, then it's ev is returned`,
			q: app.LocateQuery{},
			gr: &GroupsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ uuid.UUID) (domain.Group, error) {
					return fixtures.Group{
						Car: helpers.EvPtr(
							fixtures.Car{}.Build(),
						),
					}.Build(), nil
				},
			},
			evr: &CarsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ uuid.UUID) (domain.Car, error) {
					return fixtures.Car{}.Build(), nil
				},
			},
			expectedCallsToEvr: 1,
			expectedRs: app.LocateResponse{
				IsInJourney: true,
				Car:         fixtures.Car{}.Build(),
			},
		},
	}

	for _, tc := range testCases {
		qh := app.NewLocate(tc.gr, tc.evr)
		_, err := qh.Handle(context.Background(), tc.q)
		require.Equal(t, tc.expectedErrFunc == nil, err == nil)
		if err != nil {
			tc.expectedErrFunc(t, err)
			continue
		}

		require.Len(t, tc.gr.FindByIDCalls(), 1)
		if tc.expectedCallsToEvr > 0 {
			require.Len(t, tc.evr.FindByIDCalls(), tc.expectedCallsToEvr)
		}
	}
}
