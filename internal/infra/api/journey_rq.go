// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package api

import "fmt"
import "reflect"
import "encoding/json"

type JourneyRqJsonPeople int

var enumValues_JourneyRqJsonPeople = []interface{}{
	1,
	2,
	3,
	4,
	5,
	6,
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *JourneyRqJsonPeople) UnmarshalJSON(b []byte) error {
	var v int
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_JourneyRqJsonPeople {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_JourneyRqJsonPeople, v)
	}
	*j = JourneyRqJsonPeople(v)
	return nil
}

// Schema definition to add a group for a journey
type JourneyRqJson struct {
	// group id
	Id string `json:"id"`

	// group size. Allowed from 1 to 6
	People JourneyRqJsonPeople `json:"people"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *JourneyRqJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["id"]; !ok || v == nil {
		return fmt.Errorf("field id: required")
	}
	if v, ok := raw["people"]; !ok || v == nil {
		return fmt.Errorf("field people: required")
	}
	type Plain JourneyRqJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = JourneyRqJson(plain)
	return nil
}
