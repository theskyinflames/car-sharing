package app

import (
	"context"
	"fmt"
	"log"
)

// InvalidQueryError is self described
type InvalidQueryError struct {
	Expected string
	Had      string
}

const errMsgInvalidQuery = "invalid query, expected '%s' but found '%s'"

// Error implements the error.Error interface
func (ewq InvalidQueryError) Error() string {
	return fmt.Sprintf(errMsgInvalidQuery, ewq.Expected, ewq.Had)
}

// NewInvalidQueryError is a constructor
func NewInvalidQueryError(expected, had string) InvalidQueryError {
	return InvalidQueryError{Expected: expected, Had: had}
}

// Query is the interface to identify the DTO for a given query by name.
type Query interface {
	Name() string
}

// QueryName is string to identify a given query when it has not input parameters.
type QueryName string

// Name implements Query interface
func (qn QueryName) Name() string {
	return string(qn)
}

// QueryResult is a generic query result type.
type QueryResult interface{}

// QueryHandler is the interface for handling queries.
type QueryHandler interface {
	Handle(ctx context.Context, query Query) (QueryResult, error)
}

type queryHandlerFunc func(ctx context.Context, query Query) (QueryResult, error)

func (f queryHandlerFunc) Handle(ctx context.Context, query Query) (QueryResult, error) {
	return f(ctx, query)
}

// QueryHandlerMiddleware is a type for decorating QueryHandlers
type QueryHandlerMiddleware func(h QueryHandler) QueryHandler

// QueryHandlerMultiMiddleware applies a sequence of middlewares to a given query handler.
func QueryHandlerMultiMiddleware(middlewares ...QueryHandlerMiddleware) QueryHandlerMiddleware {
	return func(h QueryHandler) QueryHandler {
		handler := h
		for _, m := range middlewares {
			handler = m(handler)
		}
		return queryHandlerFunc(handler.Handle)
	}
}

// QueryHandlerErrorWrapperMiddleware wraps query handler errors with query name, Ex:
// query name: getUser
// error: not found
// result: getUser: not found
func QueryHandlerErrorWrapperMiddleware() QueryHandlerMiddleware {
	return func(h QueryHandler) QueryHandler {
		return queryHandlerFunc(func(ctx context.Context, q Query) (QueryResult, error) {
			result, err := h.Handle(ctx, q)
			if err != nil {
				log.Printf("ERR: %s: %s", q.Name(), err.Error())
			}

			return result, err
		})
	}
}
