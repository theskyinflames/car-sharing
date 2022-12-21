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
	b, _ := json.Marshal(map[string]interface{}{
		"seats": car.Capacity().Int(),
	})
	return CarCreatedEvent{
		EventBasic: events.NewEventBasic(car.ID(), CarCreatedEventName, b),
	}
}

// GroupSetOnJourneyEventName is self-described
const GroupSetOnJourneyEventName = "group.is.on.journey"

// GroupSetOnJourneyEvent is an event
type GroupSetOnJourneyEvent struct {
	events.EventBasic
}

// NewGroupSetOnJourneyEvent is a constructor
func NewGroupSetOnJourneyEvent(g Group) GroupSetOnJourneyEvent {
	b, _ := json.Marshal(map[string]interface{}{
		"people": g.People(),
	})
	return GroupSetOnJourneyEvent{
		EventBasic: events.NewEventBasic(g.ID(), GroupSetOnJourneyEventName, b),
	}
}
