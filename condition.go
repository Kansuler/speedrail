package speedrail

// Condition is a function that will be a condition for strategies.
type Condition[M any] func(M) bool

// And receives conditions, if all of them are true condition is passed.
func And[M any](conditions ...Condition[M]) Condition[M] {
	return func(model M) bool {
		for _, condition := range conditions {
			if !condition(model) {
				return false
			}
		}

		return true
	}
}

// Or receives conditions, if one of them is true condition is passed.
func Or[M any](conditions ...Condition[M]) Condition[M] {
	return func(model M) bool {
		for _, condition := range conditions {
			if condition(model) {
				return true
			}
		}

		return false
	}
}

// Not will invert a condition
func Not[M any](condition Condition[M]) Condition[M] {
	return func(model M) bool {
		return !condition(model)
	}
}
