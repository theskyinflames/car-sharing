package app_test

import (
	"context"
	"errors"
	"testing"
	"theskyinflames/car-sharing/internal/app"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type wrongQuery struct{}

func (wq wrongQuery) Name() string {
	return "wrongQuery"
}

func TestQueryMultiMiddleware(t *testing.T) {
	assert := assert.New(t)

	var mwExecutionOrder []int
	queryHandlerMw := func(count int) app.QueryHandlerMiddleware {
		return func(h app.QueryHandler) app.QueryHandler {
			return &QueryHandlerMock{
				HandleFunc: func(ctx context.Context, query app.Query) (app.QueryResult, error) {
					_, _ = h.Handle(ctx, query)
					mwExecutionOrder = append(mwExecutionOrder, count)
					return nil, nil
				},
			}
		}
	}

	multiMw := app.QueryHandlerMultiMiddleware(
		queryHandlerMw(1),
		queryHandlerMw(2),
		queryHandlerMw(3),
	)

	handlerExecutionCount := 0
	qh := &QueryHandlerMock{
		HandleFunc: func(context.Context, app.Query) (app.QueryResult, error) {
			handlerExecutionCount++
			return nil, nil
		},
	}

	_, err := multiMw(qh).Handle(context.Background(), nil)
	assert.NoError(err)
	assert.Equal([]int{1, 2, 3}, mwExecutionOrder)
	assert.Equal(1, handlerExecutionCount)
}

func TestQueryHandlerErrorWrapperMiddleware(t *testing.T) {
	mw := app.QueryHandlerErrorWrapperMiddleware()

	q := &QueryMock{
		NameFunc: func() string { return "someQuery" },
	}

	expectedResult := struct{}{}

	tests := []struct {
		whenThen          string
		qErr, expectedErr error
	}{
		{
			whenThen:    "and a query handler that does not return error, then it does not return any error",
			qErr:        nil,
			expectedErr: nil,
		},
		{
			whenThen:    "and a query handler that not returns error, then it returns the error wrapped",
			qErr:        errors.New("error"),
			expectedErr: errors.New("error"),
		},
	}
	for _, tt := range tests {
		name := "Given a query handler error wrapper middleware " + tt.whenThen
		t.Run(name, func(t *testing.T) {
			require := require.New(t)

			ch := &QueryHandlerMock{
				HandleFunc: func(context.Context, app.Query) (app.QueryResult, error) {
					return expectedResult, tt.qErr
				},
			}

			result, err := mw(ch).Handle(context.Background(), q)
			require.Equal(expectedResult, result)
			if tt.expectedErr != nil {
				require.Equal(tt.expectedErr.Error(), err.Error())
			}
		})
	}
}
