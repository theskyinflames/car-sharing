package app_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"theskyinflames/car-sharing/internal/app"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandHandlerMultiMiddleware(t *testing.T) {
	assert := assert.New(t)

	var mwExecutionOrder []int
	commandHandlerMw := func(count int) app.CommandHandlerMiddleware {
		return func(h app.CommandHandler) app.CommandHandler {
			return &CommandHandlerMock{
				HandleFunc: func(ctx context.Context, cmd app.Command) error {
					_ = h.Handle(ctx, cmd)
					mwExecutionOrder = append(mwExecutionOrder, count)
					return nil
				},
			}
		}
	}

	multiMw := app.CommandHandlerMultiMiddleware(
		commandHandlerMw(1),
		commandHandlerMw(2),
		commandHandlerMw(3),
	)

	handlerExecutionCount := 0
	ch := &CommandHandlerMock{
		HandleFunc: func(context.Context, app.Command) error {
			handlerExecutionCount++
			return nil
		},
	}

	err := multiMw(ch).Handle(context.Background(), nil)
	assert.NoError(err)
	assert.Equal([]int{1, 2, 3}, mwExecutionOrder)
	assert.Equal(1, handlerExecutionCount)
}

func TestNewErrInvalidCommand(t *testing.T) {
	t.Run(`Given two command names,
	when the constructor is called,
	then it returns a InvalidCommandError containing the information for the command names`, func(t *testing.T) {
		const had = "had"
		const expected = "expected"
		err := app.NewInvalidCommandError(expected, had)
		require.Equal(t, fmt.Sprintf("invalid command, expected '%s' but found '%s'", expected, had), err.Error())
	})
}

func TestCommandHandlerErrorWrapperMiddleware(t *testing.T) {
	mw := app.CommandHandlerErrorWrapperMiddleware()
	randomErr := errors.New("random")

	cmd := &CommandMock{
		NameFunc: func() string { return "someCommand" },
	}

	tests := []struct {
		whenThen            string
		cmdErr, expectedErr error
	}{
		{
			whenThen:    "and a command handler that does not return error, then it does not return any error",
			cmdErr:      nil,
			expectedErr: nil,
		},
		{
			whenThen:    "and a command handler that not returns error, then it returns the error wrapped",
			cmdErr:      randomErr,
			expectedErr: randomErr,
		},
	}
	for _, tt := range tests {
		name := "Given a command handler error wrapper middleware " + tt.whenThen
		t.Run(name, func(t *testing.T) {
			require := require.New(t)

			ch := &CommandHandlerMock{
				HandleFunc: func(context.Context, app.Command) error {
					return tt.cmdErr
				},
			}

			err := mw(ch).Handle(context.Background(), cmd)
			require.ErrorIs(err, tt.expectedErr)
		})
	}
}
