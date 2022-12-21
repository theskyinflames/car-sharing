package app

import "github.com/theskyinflames/cqrs-eda/pkg/cqrs"

func BuildChEventsMiddleware(eventsBus cqrs.Bus) cqrs.CommandHandlerMiddleware {
	return nil
}
