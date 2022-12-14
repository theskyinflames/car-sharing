// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package api

import "fmt"
import "reflect"
import "encoding/json"

type LocateRsJsonSeats int

var enumValues_LocateRsJsonSeats = []interface{}{
	4,
	5,
	6,
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *LocateRsJsonSeats) UnmarshalJSON(b []byte) error {
	var v int
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_LocateRsJsonSeats {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_LocateRsJsonSeats, v)
	}
	*j = LocateRsJsonSeats(v)
	return nil
}

// Schema definition to Initialize a fleet
type LocateRsJson struct {
	// car uuid
	Id string `json:"id"`

	// ev seats
	Seats LocateRsJsonSeats `json:"seats"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *LocateRsJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["id"]; !ok || v == nil {
		return fmt.Errorf("field id: required")
	}
	if v, ok := raw["seats"]; !ok || v == nil {
		return fmt.Errorf("field seats: required")
	}
	type Plain LocateRsJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = LocateRsJson(plain)
	return nil
}
