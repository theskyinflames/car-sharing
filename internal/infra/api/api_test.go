package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"theskyinflames/car-sharing/internal/app"
	"theskyinflames/car-sharing/internal/fixtures"
	"theskyinflames/car-sharing/internal/infra/api"
	"theskyinflames/car-sharing/internal/infra/repository"

	"github.com/stretchr/testify/require"
)

func TestInitializeFleet(t *testing.T) {
	testCases := []struct {
		name           string
		rq             api.InitializeFleetRqJson
		ch             *CommandHandlerMock
		expectedStatus int
	}{
		{
			name: `Given an initialize fleet endpoint,
			when it's called with a wrong rq without id,
			then a 400 HTTP status is returned`,
			rq: api.InitializeFleetRqJson{
				{Seats: 5},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an initialize fleet endpoint,
			when it's called with a wrong rq with an id < 1,
			then a 400 HTTP status is returned`,
			rq: api.InitializeFleetRqJson{
				{Id: -1, Seats: 5},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an initialize fleet endpoint,
			when it's called with a wrong rq with a not allowed number of seats,
			then a 400 HTTP status is returned`,
			rq: api.InitializeFleetRqJson{
				{Id: 1, Seats: 3},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an initialize fleet endpoint with a ch that returns an error, 
			when it's called ,
			then a 500 HTTP status is returned`,
			rq: api.InitializeFleetRqJson{
				{Id: 1, Seats: 5},
			},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ app.Command) error {
					return errors.New("")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: `Given an initialize fleet endpoint,
			when it's called with a right rq,
			then a 200 HTTP status is returned`,
			rq: api.InitializeFleetRqJson{{Id: 1, Seats: 5}},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ app.Command) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: `Given an initialize fleet endpoint,
			when it's called with a empty rq,
			then a 200 HTTP status is returned`,
			rq: api.InitializeFleetRqJson{},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ app.Command) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		s, err := json.Marshal(tc.rq)
		require.NoError(t, err)
		reqBody := strings.NewReader(string(s))

		hnd := api.InitializeFleet(tc.ch)
		r := httptest.NewRequest("", "/evs", reqBody)
		w := httptest.NewRecorder()
		hnd(w, r)
		require.Equal(t, tc.expectedStatus, w.Code)
	}
}

func TestJourney(t *testing.T) {
	testCases := []struct {
		name           string
		rq             api.JourneyRqJson
		ch             *CommandHandlerMock
		expectedStatus int
	}{
		{
			name: `Given an journey endpoint,
			when it's called with a wrong rq without id,
			then a 400 HTTP status is returned`,
			rq:             api.JourneyRqJson{People: 5},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an journey endpoint,
			when it's called with a wrong rq with an id < 1,
			then a 400 HTTP status is returned`,
			rq:             api.JourneyRqJson{Id: -1, People: 5},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an journey endpoint,
			when it's called with a wrong rq with a not allowed group size,
			then a 400 HTTP status is returned`,
			rq:             api.JourneyRqJson{Id: 1, People: 10},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an journey endpoint with a ch that returns an error, 
			when it's called ,
			then a 500 HTTP status is returned`,
			rq: api.JourneyRqJson{Id: 1, People: 5},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ app.Command) error {
					return errors.New("")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: `Given an journey endpoint,
			when it's called with a empty rq,
			then a 400 HTTP status is returned`,
			rq: api.JourneyRqJson{},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ app.Command) error {
					return nil
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an journey endpoint,
			when it's called with a right rq,
			then a 200 HTTP status is returned`,
			rq: api.JourneyRqJson{Id: 1, People: 5},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ app.Command) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		s, err := json.Marshal(tc.rq)
		require.NoError(t, err)
		reqBody := strings.NewReader(string(s))

		hnd := api.Journey(tc.ch)
		r := httptest.NewRequest("", "/journey", reqBody)
		w := httptest.NewRecorder()
		hnd(w, r)
		require.Equal(t, tc.expectedStatus, w.Code)
	}
}

func TestDropOff(t *testing.T) {
	testCases := []struct {
		name           string
		rq             api.DropOffRqJson
		ch             *CommandHandlerMock
		expectedStatus int
	}{
		{
			name: `Given a drop off endpoint,
			when it's called with a wrong rq without id,
			then a 400 HTTP status is returned`,
			rq:             api.DropOffRqJson{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given a drop off endpoint,
			when it's called with a wrong rq with an id < 1,
			then a 400 HTTP status is returned`,
			rq:             api.DropOffRqJson{Id: -1},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given a drop off endpoint with a ch that returns an error other than not found, 
			when it's called ,
			then a 500 HTTP status is returned`,
			rq: api.DropOffRqJson{Id: 1},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ app.Command) error {
					return errors.New("")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: `Given a drop off endpoint with a ch that returns an not found error, 
			when it's called ,
			then a 404 HTTP status is returned`,
			rq: api.DropOffRqJson{Id: 1},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ app.Command) error {
					return repository.ErrNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: `Given a drop off endpoint,
			when it's called with a right rq,
			then a 200 HTTP status is returned`,
			rq: api.DropOffRqJson{Id: 1},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ app.Command) error {
					return nil
				},
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		s, err := json.Marshal(tc.rq)
		require.NoError(t, err)
		reqBody := strings.NewReader(string(s))

		hnd := api.DropOff(tc.ch)
		r := httptest.NewRequest("", "/dropoff", reqBody)
		w := httptest.NewRecorder()
		hnd(w, r)
		require.Equal(t, tc.expectedStatus, w.Code)
	}
}

func TestLocale(t *testing.T) {
	ev := fixtures.Car{}.Build()
	testCases := []struct {
		name           string
		rq             api.LocateRqJson
		expectedRs     *api.LocateRsJson
		qh             *QueryHandlerMock
		expectedStatus int
	}{
		{
			name: `Given locale endpoint,
			when it's called with a wrong rq without id,
			then a 400 HTTP status is returned`,
			rq:             api.LocateRqJson{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given locale endpoint,
			when it's called with a wrong rq with an id < 1,
			then a 400 HTTP status is returned`,
			rq:             api.LocateRqJson{Id: -1},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given a locale endpoint with a qh that returns an error other than not found,
			when it's called ,
			then a 500 HTTP status is returned`,
			rq: api.LocateRqJson{Id: 1},
			qh: &QueryHandlerMock{
				HandleFunc: func(ctx context.Context, query app.Query) (app.QueryResult, error) {
					return nil, errors.New("")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: `Given locale endpoint with a ch that returns an not found error,
			when it's called ,
			then a 404 HTTP status is returned`,
			rq: api.LocateRqJson{Id: 1},
			qh: &QueryHandlerMock{
				HandleFunc: func(ctx context.Context, query app.Query) (app.QueryResult, error) {
					return nil, repository.ErrNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: `Given locale endpoint and a waiting group,
			when it's called with a right rq,
			then a 200 HTTP status is returned`,
			rq: api.LocateRqJson{Id: 1},
			qh: &QueryHandlerMock{
				HandleFunc: func(ctx context.Context, query app.Query) (app.QueryResult, error) {
					return app.LocateResponse{}, nil
				},
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: `Given locale endpoint and an journey group,
			when it's called with a right rq,
			then a 200 HTTP status is returned`,
			rq: api.LocateRqJson{Id: 1},
			qh: &QueryHandlerMock{
				HandleFunc: func(ctx context.Context, query app.Query) (app.QueryResult, error) {
					return app.LocateResponse{
						IsInJourney: true,
						Ev:          ev,
					}, nil
				},
			},
			expectedRs: &api.LocateRsJson{
				Id:    ev.ID(),
				Seats: api.LocateRsJsonSeats(ev.Capacity()),
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		s, err := json.Marshal(tc.rq)
		require.NoError(t, err)
		reqBody := strings.NewReader(string(s))

		hnd := api.Locate(tc.qh)
		r := httptest.NewRequest("", "/locale", reqBody)
		w := httptest.NewRecorder()
		hnd(w, r)

		require.Equal(t, tc.expectedStatus, w.Code)
		if tc.expectedRs == nil {
			continue
		}
		buff := &bytes.Buffer{}
		buff.ReadFrom(w.Body)

		var rs api.LocateRsJson
		require.NoError(t, json.Unmarshal(buff.Bytes(), &rs))
		require.Equal(t, *tc.expectedRs, rs)
	}
}
