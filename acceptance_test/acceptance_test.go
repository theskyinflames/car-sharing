package acceptantce_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"theskyinflames/car-sharing/internal/infra/api"
	"theskyinflames/car-sharing/internal/infra/server"

	"github.com/stretchr/testify/require"
)

func TestAcceptanceTest(t *testing.T) {
	log := log.New(os.Stdout, "car-sharing: ", os.O_APPEND)
	bus := server.BuildBus(log)

	t.Run(`Given a car-sharing API `, func(t *testing.T) {
		t.Run(`when cars endpoint is called, then these cars are added`, func(t *testing.T) {
			rq := []api.Cars{
				{
					Id:    1,
					Seats: 4,
				},
				{
					Id:    2,
					Seats: 6,
				},
			}
			do(t, http.HandlerFunc(api.InitializeFleet(bus)), http.MethodPut, "/cars", rq, http.StatusOK, nil, nil)
		})

		t.Run(`when a journey with a group of 3 people is added, then it's got on the six-seat car`, func(t *testing.T) {
			rq := api.JourneyRqJson{
				Id:     1,
				People: 3,
			}
			do(t, http.HandlerFunc(api.Journey(bus)), http.MethodPut, "/journey", rq, http.StatusOK, nil, nil)
		})

		t.Run(`when the same journey tried to be added, then a 409 HTTP status is returned`, func(t *testing.T) {
			rq := api.JourneyRqJson{
				Id:     1,
				People: 3,
			}
			do(t, http.HandlerFunc(api.Journey(bus)), http.MethodPut, "/journey", rq, http.StatusConflict, nil, nil)
		})

		t.Run(`when a journey with a second group of 4 people is added, then it's got on the four-seat car`, func(t *testing.T) {
			rq := api.JourneyRqJson{
				Id:     2,
				People: 4,
			}
			do(t, http.HandlerFunc(api.Journey(bus)), http.MethodPut, "/journey", rq, http.StatusOK, nil, nil)
		})

		t.Run(`when the group with Id=1 and 3 people is located, then car with Id=2 and six seats is returned`, func(t *testing.T) {
			rq := api.LocateRqJson{
				Id: 1,
			}
			expectedRs := api.LocateRsJson{
				Id:    2,
				Seats: 6,
			}
			unmarshalRsFunc := func(t *testing.T, b []byte) any {
				var rs api.LocateRsJson
				require.NoError(t, rs.UnmarshalJSON(b))
				return rs
			}
			do(t, http.HandlerFunc(api.Locate(bus)), http.MethodPost, "/locate", rq, http.StatusOK, expectedRs, unmarshalRsFunc)
		})

		t.Run(`when the group of 4 people is located, then car with Id=1 and for seats is returned`, func(t *testing.T) {
			rq := api.LocateRqJson{
				Id: 2,
			}
			expectedRs := api.LocateRsJson{
				Id:    1,
				Seats: 4,
			}
			unmarshalRsFunc := func(t *testing.T, b []byte) any {
				var rs api.LocateRsJson
				require.NoError(t, rs.UnmarshalJSON(b))
				return rs
			}
			do(t, http.HandlerFunc(api.Locate(bus)), http.MethodPost, "/locate", rq, http.StatusOK, expectedRs, unmarshalRsFunc)
		})

		t.Run(`when two more groups are added, then the stay waiting for a car`, func(t *testing.T) {
			rqJourney := api.JourneyRqJson{
				Id:     3,
				People: 4,
			}
			do(t, http.HandlerFunc(api.Journey(bus)), http.MethodPut, "/journey", rqJourney, http.StatusOK, nil, nil)

			rqJourney = api.JourneyRqJson{
				Id:     4,
				People: 4,
			}
			do(t, http.HandlerFunc(api.Journey(bus)), http.MethodPut, "/journey", rqJourney, http.StatusOK, nil, nil)

			rqLocate := api.LocateRqJson{
				Id: 3,
			}
			do(t, http.HandlerFunc(api.Locate(bus)), http.MethodPost, "/locate", rqLocate, http.StatusNoContent, nil, nil)

			rqLocate = api.LocateRqJson{
				Id: 4,
			}
			do(t, http.HandlerFunc(api.Locate(bus)), http.MethodPost, "/locate", rqLocate, http.StatusNoContent, nil, nil)
		})

		t.Run(`when the group with Id=1 is dropped off, then the first waiting group (Id=3) is got on and the group with Id=4 keeps waiting`, func(t *testing.T) {
			rqDropOff := api.DropOffRqJson{
				Id: 1,
			}
			do(t, http.HandlerFunc(api.DropOff(bus)), http.MethodPost, "/dropoff", rqDropOff, http.StatusNoContent, nil, nil)

			rqLocate := api.LocateRqJson{
				Id: 3,
			}
			expectedRs := api.LocateRsJson{
				Id:    2,
				Seats: 6,
			}
			unmarshalRsFunc := func(t *testing.T, b []byte) any {
				var rs api.LocateRsJson
				require.NoError(t, rs.UnmarshalJSON(b))
				return rs
			}
			do(t, http.HandlerFunc(api.Locate(bus)), http.MethodPost, "/locate", rqLocate, http.StatusOK, expectedRs, unmarshalRsFunc)

			rqLocate = api.LocateRqJson{
				Id: 4,
			}
			do(t, http.HandlerFunc(api.Locate(bus)), http.MethodPost, "/locate", rqLocate, http.StatusNoContent, nil, nil)
		})
	})
}

type unmarshalRsFunc func(*testing.T, []byte) any

func do(t *testing.T, hnd http.HandlerFunc, method string, path string, rq any, statusCode int, rs any, unmarshalRsFunc unmarshalRsFunc) {
	// start the server
	srv := httptest.NewServer(http.HandlerFunc(hnd))

	// Create a new HTTP client
	client := &http.Client{}

	// Set the API endpoint URL and query parameters
	apiURL, err := url.Parse(srv.URL + path)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(rq)
	require.NoError(t, err)

	req, err := http.NewRequest(method, apiURL.String(), bytes.NewBuffer(b))
	require.NoError(t, err)
	req.Header.Add("Accept", "application/json")
	require.NoError(t, err)

	// Send the request and retrieve the response
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check the status code is what we expect.
	require.Equal(t, statusCode, resp.StatusCode)

	if rs != nil {
		// Check for the expected response
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		received := unmarshalRsFunc(t, body)
		require.Equal(t, rs, received)
	}
}
