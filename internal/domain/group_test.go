package domain_test

import (
	"testing"

	"theskyinflames/car-sharing/internal/domain"
	"theskyinflames/car-sharing/internal/fixtures"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewGroup(t *testing.T) {
	id := uuid.New()
	testCases := []struct {
		name            string
		id              uuid.UUID
		people          int
		expectedErrFunc func(*testing.T, error)
	}{
		{
			name: `Given an empty group, when it's called then an error is returned`,
			id:   id,
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, domain.ErrWrongSize)
			},
		},
		{
			name:   `Given an oversized group, when it's called then an error is returned`,
			id:     id,
			people: 7,
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, domain.ErrWrongSize)
			},
		},
		{
			name:   `Given a group, when it's called then no error is returned`,
			id:     id,
			people: 3,
		},
	}

	for _, tc := range testCases {
		g, err := domain.NewGroup(tc.id, tc.people)
		require.Equal(t, tc.expectedErrFunc == nil, err == nil)
		if err != nil {
			tc.expectedErrFunc(t, err)
			continue
		}
		require.Equal(t, tc.id, g.ID())
		require.Equal(t, tc.people, g.People())
	}
}

func TestGroupIsOnJourney(t *testing.T) {
	testCases := []struct {
		name     string
		g        domain.Group
		expected bool
	}{
		{
			name:     `Given a group without ev assigned, when it's called, then it's not on journey`,
			g:        fixtures.Group{}.Build(),
			expected: false,
		},
		{
			name:     `Given a group with ev assigned, when it's called, then it's on journey`,
			g:        fixtures.Group{}.Build(),
			expected: false,
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expected, tc.g.IsOnJourney())
	}
}
