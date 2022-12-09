package app

import (
	"context"
	"fmt"
	"log"
)

//go:generate moq -out ./zmock_cqrs_test.go -pkg app_test . CommandHandler Command QueryHandler Query QueryResult

// Command is the interface for identifying commands by name.
type Command interface {
	Name() string
}

// CommandHandler is the interface for implementing the logic that mutates our domain.
// It receives the context for both termination and transactionality injection, and a command DTO for the input data.
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

// CommandHandlerFunc is a function that implements CommandHandler interface.
type CommandHandlerFunc func(ctx context.Context, cmd Command) error

// Handle is the CommandHandler interface implementation.
func (f CommandHandlerFunc) Handle(ctx context.Context, cmd Command) error {
	return f(ctx, cmd)
}

// CommandHandlerMiddleware is a type for decorating CommandHandlers
type CommandHandlerMiddleware func(h CommandHandler) CommandHandler

// CommandHandlerMultiMiddleware applies a sequence of middlewares to a given command handler.
func CommandHandlerMultiMiddleware(middlewares ...CommandHandlerMiddleware) CommandHandlerMiddleware {
	return func(h CommandHandler) CommandHandler {
		handler := h
		for _, m := range middlewares {
			handler = m(handler)
		}
		return CommandHandlerFunc(handler.Handle)
	}
}

// InvalidCommandError should be returned by the implementations of the interface when the handler does not receive the needed command.
type InvalidCommandError struct {
	expected string
	had      string
}

// NewInvalidCommandError is a constructor
func NewInvalidCommandError(expected string, had string) InvalidCommandError {
	return InvalidCommandError{expected: expected, had: had}
}

const errMsgInvalidCommand = "invalid command, expected '%s' but found '%s'"

func (e InvalidCommandError) Error() string {
	return fmt.Sprintf(errMsgInvalidCommand, e.expected, e.had)
}

// CommandHandlerErrorWrapperMiddleware wraps command handler errors with command name, Ex:
// command name: drop off
// error: not found
// result: drop off: not found
func CommandHandlerErrorWrapperMiddleware() CommandHandlerMiddleware {
	return func(h CommandHandler) CommandHandler {
		return CommandHandlerFunc(func(ctx context.Context, cmd Command) error {
			err := h.Handle(ctx, cmd)
			if err != nil {
				log.Printf("ERR: %s: %s", cmd.Name(), err.Error())
			}
			return err
		})
	}
}
