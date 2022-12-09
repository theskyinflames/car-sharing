// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package app_test

import (
	"context"
	"theskyinflames/car-sharing/internal/app"
	"sync"
)

// Ensure, that CommandHandlerMock does implement app.CommandHandler.
// If this is not the case, regenerate this file with moq.
var _ app.CommandHandler = &CommandHandlerMock{}

// CommandHandlerMock is a mock implementation of app.CommandHandler.
//
// 	func TestSomethingThatUsesCommandHandler(t *testing.T) {
//
// 		// make and configure a mocked app.CommandHandler
// 		mockedCommandHandler := &CommandHandlerMock{
// 			HandleFunc: func(ctx context.Context, cmd app.Command) error {
// 				panic("mock out the Handle method")
// 			},
// 		}
//
// 		// use mockedCommandHandler in code that requires app.CommandHandler
// 		// and then make assertions.
//
// 	}
type CommandHandlerMock struct {
	// HandleFunc mocks the Handle method.
	HandleFunc func(ctx context.Context, cmd app.Command) error

	// calls tracks calls to the methods.
	calls struct {
		// Handle holds details about calls to the Handle method.
		Handle []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Cmd is the cmd argument value.
			Cmd app.Command
		}
	}
	lockHandle sync.RWMutex
}

// Handle calls HandleFunc.
func (mock *CommandHandlerMock) Handle(ctx context.Context, cmd app.Command) error {
	if mock.HandleFunc == nil {
		panic("CommandHandlerMock.HandleFunc: method is nil but CommandHandler.Handle was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Cmd app.Command
	}{
		Ctx: ctx,
		Cmd: cmd,
	}
	mock.lockHandle.Lock()
	mock.calls.Handle = append(mock.calls.Handle, callInfo)
	mock.lockHandle.Unlock()
	return mock.HandleFunc(ctx, cmd)
}

// HandleCalls gets all the calls that were made to Handle.
// Check the length with:
//     len(mockedCommandHandler.HandleCalls())
func (mock *CommandHandlerMock) HandleCalls() []struct {
	Ctx context.Context
	Cmd app.Command
} {
	var calls []struct {
		Ctx context.Context
		Cmd app.Command
	}
	mock.lockHandle.RLock()
	calls = mock.calls.Handle
	mock.lockHandle.RUnlock()
	return calls
}

// Ensure, that CommandMock does implement app.Command.
// If this is not the case, regenerate this file with moq.
var _ app.Command = &CommandMock{}

// CommandMock is a mock implementation of app.Command.
//
// 	func TestSomethingThatUsesCommand(t *testing.T) {
//
// 		// make and configure a mocked app.Command
// 		mockedCommand := &CommandMock{
// 			NameFunc: func() string {
// 				panic("mock out the Name method")
// 			},
// 		}
//
// 		// use mockedCommand in code that requires app.Command
// 		// and then make assertions.
//
// 	}
type CommandMock struct {
	// NameFunc mocks the Name method.
	NameFunc func() string

	// calls tracks calls to the methods.
	calls struct {
		// Name holds details about calls to the Name method.
		Name []struct {
		}
	}
	lockName sync.RWMutex
}

// Name calls NameFunc.
func (mock *CommandMock) Name() string {
	if mock.NameFunc == nil {
		panic("CommandMock.NameFunc: method is nil but Command.Name was just called")
	}
	callInfo := struct {
	}{}
	mock.lockName.Lock()
	mock.calls.Name = append(mock.calls.Name, callInfo)
	mock.lockName.Unlock()
	return mock.NameFunc()
}

// NameCalls gets all the calls that were made to Name.
// Check the length with:
//     len(mockedCommand.NameCalls())
func (mock *CommandMock) NameCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockName.RLock()
	calls = mock.calls.Name
	mock.lockName.RUnlock()
	return calls
}

// Ensure, that QueryHandlerMock does implement app.QueryHandler.
// If this is not the case, regenerate this file with moq.
var _ app.QueryHandler = &QueryHandlerMock{}

// QueryHandlerMock is a mock implementation of app.QueryHandler.
//
// 	func TestSomethingThatUsesQueryHandler(t *testing.T) {
//
// 		// make and configure a mocked app.QueryHandler
// 		mockedQueryHandler := &QueryHandlerMock{
// 			HandleFunc: func(ctx context.Context, query app.Query) (app.QueryResult, error) {
// 				panic("mock out the Handle method")
// 			},
// 		}
//
// 		// use mockedQueryHandler in code that requires app.QueryHandler
// 		// and then make assertions.
//
// 	}
type QueryHandlerMock struct {
	// HandleFunc mocks the Handle method.
	HandleFunc func(ctx context.Context, query app.Query) (app.QueryResult, error)

	// calls tracks calls to the methods.
	calls struct {
		// Handle holds details about calls to the Handle method.
		Handle []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Query is the query argument value.
			Query app.Query
		}
	}
	lockHandle sync.RWMutex
}

// Handle calls HandleFunc.
func (mock *QueryHandlerMock) Handle(ctx context.Context, query app.Query) (app.QueryResult, error) {
	if mock.HandleFunc == nil {
		panic("QueryHandlerMock.HandleFunc: method is nil but QueryHandler.Handle was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Query app.Query
	}{
		Ctx:   ctx,
		Query: query,
	}
	mock.lockHandle.Lock()
	mock.calls.Handle = append(mock.calls.Handle, callInfo)
	mock.lockHandle.Unlock()
	return mock.HandleFunc(ctx, query)
}

// HandleCalls gets all the calls that were made to Handle.
// Check the length with:
//     len(mockedQueryHandler.HandleCalls())
func (mock *QueryHandlerMock) HandleCalls() []struct {
	Ctx   context.Context
	Query app.Query
} {
	var calls []struct {
		Ctx   context.Context
		Query app.Query
	}
	mock.lockHandle.RLock()
	calls = mock.calls.Handle
	mock.lockHandle.RUnlock()
	return calls
}

// Ensure, that QueryMock does implement app.Query.
// If this is not the case, regenerate this file with moq.
var _ app.Query = &QueryMock{}

// QueryMock is a mock implementation of app.Query.
//
// 	func TestSomethingThatUsesQuery(t *testing.T) {
//
// 		// make and configure a mocked app.Query
// 		mockedQuery := &QueryMock{
// 			NameFunc: func() string {
// 				panic("mock out the Name method")
// 			},
// 		}
//
// 		// use mockedQuery in code that requires app.Query
// 		// and then make assertions.
//
// 	}
type QueryMock struct {
	// NameFunc mocks the Name method.
	NameFunc func() string

	// calls tracks calls to the methods.
	calls struct {
		// Name holds details about calls to the Name method.
		Name []struct {
		}
	}
	lockName sync.RWMutex
}

// Name calls NameFunc.
func (mock *QueryMock) Name() string {
	if mock.NameFunc == nil {
		panic("QueryMock.NameFunc: method is nil but Query.Name was just called")
	}
	callInfo := struct {
	}{}
	mock.lockName.Lock()
	mock.calls.Name = append(mock.calls.Name, callInfo)
	mock.lockName.Unlock()
	return mock.NameFunc()
}

// NameCalls gets all the calls that were made to Name.
// Check the length with:
//     len(mockedQuery.NameCalls())
func (mock *QueryMock) NameCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockName.RLock()
	calls = mock.calls.Name
	mock.lockName.RUnlock()
	return calls
}

// Ensure, that QueryResultMock does implement app.QueryResult.
// If this is not the case, regenerate this file with moq.
var _ app.QueryResult = &QueryResultMock{}

// QueryResultMock is a mock implementation of app.QueryResult.
//
// 	func TestSomethingThatUsesQueryResult(t *testing.T) {
//
// 		// make and configure a mocked app.QueryResult
// 		mockedQueryResult := &QueryResultMock{
// 		}
//
// 		// use mockedQueryResult in code that requires app.QueryResult
// 		// and then make assertions.
//
// 	}
type QueryResultMock struct {
	// calls tracks calls to the methods.
	calls struct {
	}
}
