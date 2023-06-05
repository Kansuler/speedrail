package speedrail

import "context"

// Strategy is a function that will be executed.
type Strategy[C, M any] func(context.Context, C, M) (context.Context, M, Error)

// If executes a strategy if the condition is true.
func If[C, M any](condition Condition[M], onTrue Strategy[C, M]) Strategy[C, M] {
	return func(ctx context.Context, container C, model M) (context.Context, M, Error) {
		if condition(model) {
			return onTrue(ctx, container, model)
		}

		return ctx, model, nil
	}
}

// IfElse executes a strategy if the condition is true, otherwise execute another strategy.
func IfElse[C, M any](condition Condition[M], onTrue Strategy[C, M], onFalse Strategy[C, M]) Strategy[C, M] {
	return func(ctx context.Context, container C, model M) (context.Context, M, Error) {
		if condition(model) {
			return onTrue(ctx, container, model)
		}

		return onFalse(ctx, container, model)
	}
}

// Merge executes all strategies and will not stop on error, but merge all errors together and then return any error.
func Merge[C, M any](strategies ...Strategy[C, M]) Strategy[C, M] {
	return func(ctx context.Context, container C, model M) (context.Context, M, Error) {
		var resultErr Error
		for _, strategy := range strategies {
			var err Error
			ctx, model, err = strategy(ctx, container, model)
			if err == nil {
				continue
			}

			if resultErr == nil {
				resultErr = err
				continue
			}

			resultErr = resultErr.Merge(err)
		}

		return ctx, model, resultErr
	}
}

// Group is a helper function that makes it easier to read strategies logically grouped together. They are executed in
// order. If an error is returned, the execution of the strategies will stop and error returned.
func Group[C, M any](strategies ...Strategy[C, M]) Strategy[C, M] {
	return func(ctx context.Context, container C, model M) (context.Context, M, Error) {
		for _, strategy := range strategies {
			var err Error
			ctx, model, err = strategy(ctx, container, model)
			if err != nil {
				return ctx, model, err
			}
		}

		return ctx, model, nil
	}
}

// ThrowError will return a defined error.
func ThrowError[C, M any](err Error) Strategy[C, M] {
	return func(ctx context.Context, container C, model M) (context.Context, M, Error) {
		return ctx, model, err
	}
}
