package speedrail

import (
	"context"
	"errors"
	"net/http"
)

// Speedrail is a composition of strategies that will be executed in order.
type Speedrail[C, M any] []Strategy[C, M]

// Plan will assemble a list of Strategy to StrategyList.
func Plan[C, M any](strategies ...Strategy[C, M]) Speedrail[C, M] {
	return strategies
}

// ErrNoStrategy is the error returned when no strategies are provided to execute.
var ErrNoStrategy = errors.New("no strategies to execute")

// ErrNoContextReturned is the error returned when no context is returned by a strategy.
var ErrNoContextReturned = errors.New("no context returned by strategy")

// Execute executes a list of strategies.
func (s Speedrail[C, M]) Execute(ctx context.Context, container C, model M) (context.Context, M, Error) {
	if s == nil {
		return ctx, model, NewError(ErrNoStrategy, http.StatusInternalServerError, "no strategies to execute")
	}

	for _, strategy := range s {
		var err Error
		ctx, model, err = strategy(ctx, container, model)
		if ctx == nil {
			return ctx, model, NewError(ErrNoContextReturned, http.StatusInternalServerError, "no context returned by strategy")
		}

		if err != nil {
			return ctx, model, err
		}
	}

	return ctx, model, nil
}
