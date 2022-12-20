package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"theskyinflames/car-sharing/internal/app"
	"theskyinflames/car-sharing/internal/fixtures"
	"theskyinflames/car-sharing/internal/infra/api"
	"theskyinflames/car-sharing/internal/infra/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/cqrs-eda/pkg/bus"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
	"github.com/theskyinflames/cqrs-eda/pkg/helpers"
)

type dispatchableMock struct {
	name string
}

func (d dispatchableMock) Name() string {
	return d.name
}

func TestInitializeFleet(t *testing.T) {
	testCases := []struct {
		name           string
		rq             api.CarsRqJson
		headers        map[string]string
		ch             *CommandHandlerMock
		expectedStatus int
	}{
		{
			name: `Given an initialize fleet endpoint,
			when it's called without "Content-type: application/json" header,
			then a 400 HTTP status is returned`,
			rq: api.CarsRqJson{
				{Seats: 6},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an initialize fleet endpoint,
			when it's called with a wrong rq without id,
			then a 400 HTTP status is returned`,
			rq: api.CarsRqJson{
				{Seats: 5},
			},
			headers:        map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an initialize fleet endpoint,
			when it's called with a wrong rq with a wrong id,
			then a 400 HTTP status is returned`,
			rq: api.CarsRqJson{
				{Id: "wrongID", Seats: 5},
			},
			headers:        map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an initialize fleet endpoint,
			when it's called with a wrong rq with a not allowed number of seats,
			then a 400 HTTP status is returned`,
			rq: api.CarsRqJson{
				{Id: uuid.New().String(), Seats: 3},
			},
			headers:        map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an initialize fleet endpoint with a ch that returns an error, 
			when it's called ,
			then a 500 HTTP status is returned`,
			rq: api.CarsRqJson{
				{Id: uuid.New().String(), Seats: 5},
			},
			headers: map[string]string{"Content-Type": "application/json"},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]cqrs.Event, error) {
					return nil, errors.New("")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: `Given an initialize fleet endpoint,
			when it's called with a right rq,
			then a 200 HTTP status is returned`,
			rq:      api.CarsRqJson{{Id: uuid.New().String(), Seats: 5}},
			headers: map[string]string{"Content-Type": "application/json"},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]cqrs.Event, error) {
					return nil, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: `Given an initialize fleet endpoint,
			when it's called with a empty rq,
			then a 200 HTTP status is returned`,
			rq:      api.CarsRqJson{},
			headers: map[string]string{"Content-Type": "application/json"},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]cqrs.Event, error) {
					return nil, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		s, err := json.Marshal(tc.rq)
		require.NoError(t, err)
		reqBody := strings.NewReader(string(s))

		bus := bus.New()
		bus.Register(app.InitializeFleetName, helpers.BusChHandler(tc.ch))

		hnd := api.InitializeFleet(bus)
		r := httptest.NewRequest("", "/cars", reqBody)
		for h, v := range tc.headers {
			r.Header.Add(h, v)
		}
		w := httptest.NewRecorder()
		hnd(w, r)
		require.Equal(t, tc.expectedStatus, w.Code, tc.name)
	}
}

func TestJourney(t *testing.T) {
	gID := uuid.New().String()
	testCases := []struct {
		name           string
		rq             api.JourneyRqJson
		headers        map[string]string
		ch             *CommandHandlerMock
		expectedStatus int
	}{
		{
			name: `Given a journey endpoint,
			when it's called without "Content-type: application/json" header,
			then a 400 HTTP status is returned`,
			rq:             api.JourneyRqJson{People: 6},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given a journey endpoint,
			when it's called with a wrong rq without id,
			then a 400 HTTP status is returned`,
			rq:             api.JourneyRqJson{People: 5},
			headers:        map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an journey endpoint,
			when it's called with a wrong rq with a wrong id,
			then a 400 HTTP status is returned`,
			rq:             api.JourneyRqJson{Id: "wrongID", People: 5},
			headers:        map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an journey endpoint,
			when it's called with a wrong rq with a not allowed group size,
			then a 400 HTTP status is returned`,
			rq:             api.JourneyRqJson{Id: gID, People: 10},
			headers:        map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an journey endpoint with a ch that returns an error, 
			when it's called ,
			then a 500 HTTP status is returned`,
			rq:      api.JourneyRqJson{Id: gID, People: 5},
			headers: map[string]string{"Content-Type": "application/json"},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]cqrs.Event, error) {
					return nil, errors.New("")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: `Given an journey endpoint,
			when it's called with a empty rq,
			then a 400 HTTP status is returned`,
			rq:      api.JourneyRqJson{},
			headers: map[string]string{"Content-Type": "application/json"},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]cqrs.Event, error) {
					return nil, nil
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given an journey endpoint,
			when it's called with a right rq,
			then a 200 HTTP status is returned`,
			rq:      api.JourneyRqJson{Id: gID, People: 5},
			headers: map[string]string{"Content-Type": "application/json"},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]cqrs.Event, error) {
					return nil, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		s, err := json.Marshal(tc.rq)
		require.NoError(t, err)
		reqBody := strings.NewReader(string(s))

		bus := bus.New()
		bus.Register(app.JourneyName, helpers.BusChHandler(tc.ch))

		hnd := api.Journey(bus)
		r := httptest.NewRequest("", "/journey", reqBody)
		for h, v := range tc.headers {
			r.Header.Add(h, v)
		}
		w := httptest.NewRecorder()
		hnd(w, r)
		require.Equal(t, tc.expectedStatus, w.Code, tc.name)
	}
}

func TestDropOff(t *testing.T) {
	gID := uuid.New().String()
	testCases := []struct {
		name           string
		headers        map[string]string
		ch             *CommandHandlerMock
		expectedStatus int
	}{
		{
			name: `Given a journey endpoint,
			when it's called without "Content-type: application/x-www-form-urlencoded" header,
			then a 400 HTTP status is returned`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given a drop off endpoint with a ch that returns an error other than not found, 
			when it's called ,
			then a 500 HTTP status is returned`,
			headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]cqrs.Event, error) {
					return nil, errors.New("")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: `Given a drop off endpoint with a ch that returns an not found error, 
			when it's called ,
			then a 404 HTTP status is returned`,
			headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]cqrs.Event, error) {
					return nil, repository.ErrNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: `Given a drop off endpoint,
			when it's called with a right rq,
			then a 200 HTTP status is returned`,
			headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
			ch: &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]cqrs.Event, error) {
					return nil, nil
				},
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		params := url.Values{}
		params.Set("ID", gID)

		bus := bus.New()
		bus.Register(app.DropOffName, helpers.BusChHandler(tc.ch))

		hnd := api.DropOff(bus)
		r := httptest.NewRequest(http.MethodPost, "/dropoff", bytes.NewBufferString(params.Encode()))
		for h, v := range tc.headers {
			r.Header.Add(h, v)
		}
		w := httptest.NewRecorder()
		hnd(w, r)
		require.Equal(t, tc.expectedStatus, w.Code)
	}
}

func TestLocale(t *testing.T) {
	gID := uuid.New().String()
	car := fixtures.Car{}.Build()
	testCases := []struct {
		name           string
		headers        map[string]string
		expectedRs     *api.LocateRsJson
		qh             *QueryHandlerMock
		expectedStatus int
	}{
		{
			name: `Given a locale endpoint,
			when it's called without "Content-type: application/x-www-form-urlencoded" header,
			then a 400 HTTP status is returned`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: `Given a locale endpoint with a qh that returns an error other than not found,
			when it's called ,
			then a 500 HTTP status is returned`,
			headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
			qh: &QueryHandlerMock{
				HandleFunc: func(ctx context.Context, query cqrs.Query) (cqrs.QueryResult, error) {
					return nil, errors.New("")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: `Given locale endpoint with a ch that returns an not found error,
			when it's called ,
			then a 404 HTTP status is returned`,
			headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
			qh: &QueryHandlerMock{
				HandleFunc: func(ctx context.Context, query cqrs.Query) (cqrs.QueryResult, error) {
					return nil, repository.ErrNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: `Given locale endpoint and a waiting group,
			when it's called with a right rq,
			then a 200 HTTP status is returned`,
			headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
			qh: &QueryHandlerMock{
				HandleFunc: func(ctx context.Context, query cqrs.Query) (cqrs.QueryResult, error) {
					return app.LocateResponse{}, nil
				},
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: `Given locale endpoint and an journey group,
			when it's called with a right rq,
			then a 200 HTTP status is returned`,
			headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
			qh: &QueryHandlerMock{
				HandleFunc: func(ctx context.Context, query cqrs.Query) (cqrs.QueryResult, error) {
					return app.LocateResponse{
						IsInJourney: true,
						Car:         car,
					}, nil
				},
			},
			expectedRs: &api.LocateRsJson{
				Id:    car.ID().String(),
				Seats: api.LocateRsJsonSeats(car.Capacity()),
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		params := url.Values{}
		params.Set("ID", gID)

		bus := bus.New()
		bus.Register(app.LocateName, helpers.BusQhHandler(tc.qh))

		hnd := api.Locate(bus)
		r := httptest.NewRequest(http.MethodPost, "/locale", bytes.NewBufferString(params.Encode()))
		for h, v := range tc.headers {
			r.Header.Add(h, v)
		}
		w := httptest.NewRecorder()
		hnd(w, r)

		require.Equal(t, tc.expectedStatus, w.Code, tc.name)
		if tc.expectedRs == nil {
			continue
		}
		buff := &bytes.Buffer{}
		buff.ReadFrom(w.Body)

		var rs api.LocateRsJson
		require.NoError(t, json.Unmarshal(buff.Bytes(), &rs))
		require.Equal(t, *tc.expectedRs, rs)

		require.Equal(t, "application/json", w.Header().Get("Accept"))
	}
}
