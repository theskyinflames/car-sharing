package domain

import (
	"encoding/json"

	"github.com/theskyinflames/cqrs-eda/pkg/events"
)

// CarCreatedEventName is self-described
const CarCreatedEventName = "car.added"

// CarCreatedEvent is an event
type CarCreatedEvent struct {
	events.EventBasic
}

// NewCarCreatedEvent is a constructor
func NewCarCreatedEvent(car Car) CarCreatedEvent {
	b, _ := json.Marshal(map[string]int{
		"seats": car.Capacity().Int(),
	})
	return CarCreatedEvent{
		EventBasic: events.NewEventBasic(car.ID(), CarCreatedEventName, b),
	}
}
