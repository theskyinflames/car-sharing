package acceptantce_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"theskyinflames/car-sharing/cmd/service"
	"theskyinflames/car-sharing/internal/app"
	"theskyinflames/car-sharing/internal/infra/api"

	"github.com/stretchr/testify/require"
)

const srvPort = ":8080"

func TestAcceptanceTest(t *testing.T) {
	log := log.New(os.Stdout, "car-sharing: ", os.O_APPEND)
	commandBus := app.BuildCommandBus(log)

	ctx, cancel := context.WithCancel(context.Background())
	go service.Run(ctx, srvPort)
	defer cancel()

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

			b, err := json.Marshal(rq)
			require.NoError(t, err)

			do(t, doCmd{
				http.HandlerFunc(api.InitializeFleet(commandBus)),
				http.MethodPut,
				"/cars",
				b,
				map[string]string{"Content-Type": "application/json"},
				http.StatusOK,
				nil,
				nil,
			})
		})

		t.Run(`when a journey with a group of 3 people is added, then it's got on the six-seat car`, func(t *testing.T) {
			rq := api.JourneyRqJson{
				Id:     1,
				People: 3,
			}
			b, err := json.Marshal(rq)
			require.NoError(t, err)
			do(t, doCmd{
				http.HandlerFunc(api.Journey(commandBus)),
				http.MethodPost,
				"/journey",
				b,
				map[string]string{"Content-Type": "application/json"},
				http.StatusOK,
				nil,
				nil,
			})
		})

		t.Run(`when the same journey tried to be added, then a 409 HTTP status is returned`, func(t *testing.T) {
			rq := api.JourneyRqJson{
				Id:     1,
				People: 3,
			}
			b, err := json.Marshal(rq)
			require.NoError(t, err)
			do(t, doCmd{
				http.HandlerFunc(api.Journey(commandBus)),
				http.MethodPost,
				"/journey",
				b,
				map[string]string{"Content-Type": "application/json"},
				http.StatusBadRequest,
				nil,
				nil,
			})
		})

		t.Run(`when a journey with a second group of 4 people is added, then it's got on the four-seat car`, func(t *testing.T) {
			rq := api.JourneyRqJson{
				Id:     2,
				People: 4,
			}
			b, err := json.Marshal(rq)
			require.NoError(t, err)
			do(t, doCmd{
				http.HandlerFunc(api.Journey(commandBus)),
				http.MethodPost,
				"/journey",
				b,
				map[string]string{"Content-Type": "application/json"},
				http.StatusOK,
				nil,
				nil,
			})
		})

		t.Run(`when the group with Id=1 and 3 people is located, then car with Id=2 and six seats is returned`, func(t *testing.T) {
			expectedRs := api.LocateRsJson{
				Id:    2,
				Seats: 6,
			}
			unmarshalRsFunc := func(t *testing.T, b []byte) any {
				var rs api.LocateRsJson
				require.NoError(t, rs.UnmarshalJSON(b))
				return rs
			}
			do(t, doCmd{
				http.HandlerFunc(api.Locate(commandBus)),
				http.MethodPost,
				"/locate",
				[]byte("ID=1"),
				map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
				http.StatusOK,
				expectedRs,
				unmarshalRsFunc,
			})
		})

		t.Run(`when the group of 4 people is located, then car with Id=1 and for seats is returned`, func(t *testing.T) {
			expectedRs := api.LocateRsJson{
				Id:    1,
				Seats: 4,
			}
			unmarshalRsFunc := func(t *testing.T, b []byte) any {
				var rs api.LocateRsJson
				require.NoError(t, rs.UnmarshalJSON(b))
				return rs
			}
			do(t, doCmd{
				http.HandlerFunc(api.Locate(commandBus)),
				http.MethodPost,
				"/locate",
				[]byte("ID=2"),
				map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
				http.StatusOK,
				expectedRs,
				unmarshalRsFunc,
			})
		})

		t.Run(`when two more groups are added, then the stay waiting for a car`, func(t *testing.T) {
			rqJourney := api.JourneyRqJson{
				Id:     3,
				People: 4,
			}
			b, err := json.Marshal(rqJourney)
			require.NoError(t, err)
			do(t, doCmd{
				http.HandlerFunc(api.Journey(commandBus)),
				http.MethodPost,
				"/journey",
				b,
				map[string]string{"Content-Type": "application/json"},
				http.StatusOK,
				nil,
				nil,
			})

			rqJourney = api.JourneyRqJson{
				Id:     4,
				People: 4,
			}
			b, err = json.Marshal(rqJourney)
			require.NoError(t, err)
			do(t, doCmd{
				http.HandlerFunc(api.Journey(commandBus)),
				http.MethodPost,
				"/journey",
				b,
				map[string]string{"Content-Type": "application/json"},
				http.StatusOK,
				nil,
				nil,
			})

			do(t, doCmd{
				http.HandlerFunc(api.Locate(commandBus)),
				http.MethodPost,
				"/locate",
				[]byte("ID=3"),
				map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
				http.StatusNoContent,
				nil,
				nil,
			})

			do(t, doCmd{
				http.HandlerFunc(api.Locate(commandBus)),
				http.MethodPost,
				"/locate",
				[]byte("ID=4"),
				map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
				http.StatusNoContent,
				nil,
				nil,
			})
		})

		t.Run(`when the group with Id=1 is dropped off, then the first waiting group (Id=3) is got on and the group with Id=4 keeps waiting`, func(t *testing.T) {
			do(t, doCmd{
				http.HandlerFunc(api.DropOff(commandBus)),
				http.MethodPost,
				"/dropoff",
				[]byte("ID=1"),
				map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
				http.StatusNoContent,
				nil,
				nil,
			})

			expectedRs := api.LocateRsJson{
				Id:    2,
				Seats: 6,
			}
			unmarshalRsFunc := func(t *testing.T, b []byte) any {
				var rs api.LocateRsJson
				require.NoError(t, rs.UnmarshalJSON(b))
				return rs
			}
			do(t, doCmd{
				http.HandlerFunc(api.Locate(commandBus)),
				http.MethodPost,
				"/locate",
				[]byte("ID=3"),
				map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
				http.StatusOK,
				expectedRs,
				unmarshalRsFunc,
			})

			do(t, doCmd{
				http.HandlerFunc(api.Locate(commandBus)),
				http.MethodPost,
				"/locate",
				[]byte("ID=4"),
				map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
				http.StatusNoContent,
				nil,
				nil,
			})
		})
	})
}

type (
	doCmd struct {
		hnd             http.HandlerFunc
		method          string
		path            string
		rq              []byte
		rqHeaders       map[string]string
		statusCode      int
		rs              any
		unmarshalRsFunc unmarshalRsFunc
	}
	unmarshalRsFunc func(*testing.T, []byte) any
)

func do(t *testing.T, doCmd doCmd) {
	// start the server
	// srv := httptest.NewServer(http.HandlerFunc(doCmd.hnd))

	// Create a new HTTP client
	client := &http.Client{}

	// Set the API endpoint URL and query parameters
	apiURL, err := url.Parse("http://localhost" + srvPort + doCmd.path)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(doCmd.method, apiURL.String(), bytes.NewBuffer(doCmd.rq))
	require.NoError(t, err)
	for h, v := range doCmd.rqHeaders {
		req.Header.Add(h, v)
	}
	require.NoError(t, err)

	// Send the request and retrieve the response
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check the status code is what we expect.
	require.Equal(t, doCmd.statusCode, resp.StatusCode)

	if doCmd.rs != nil {
		// Check for the expected response
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		received := doCmd.unmarshalRsFunc(t, body)
		require.Equal(t, doCmd.rs, received)
	}
}
