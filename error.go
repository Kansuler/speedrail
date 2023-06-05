package speedrail

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Error is a speedrail error interface that can be used to wrap any error
type Error interface {
	error
	json.Marshaler
	Trail() Trail
	Merge(Error) Error
	StatusCode() int
}

// defaultError is the default error struct for speedrail.
type defaultError struct {
	trail           Trail
	outgoingMessage string
	statusCode      int
}

// Type check that default error implements Error interface
var _ Error = defaultError{}

func (e defaultError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"error": e.Error()})
}

type ErrorWithTrail struct {
	StrategyName string
	Error        error
}

// Trail is a map of errors, with a custom marshaler.
type Trail []ErrorWithTrail

func (t Trail) MarshalJSON() ([]byte, error) {
	result := map[string]string{}
	for index, err := range t {
		if strings.Contains(err.StrategyName, "/") && strings.LastIndex(err.StrategyName, "/") < len(err.StrategyName) {
			result[fmt.Sprintf("[%d]%s", index+1, err.StrategyName[strings.LastIndex(err.StrategyName, "/")+1:])] = err.Error.Error()
			continue
		}

		result[fmt.Sprintf("[%d]%s", index+1, err.StrategyName)] = err.Error.Error()

	}

	return json.Marshal(result)
}

// Trail returns the error trail
func (e defaultError) Trail() Trail {
	return e.trail
}

// Merge will merge two errors together.
func (e defaultError) Merge(err Error) Error {
	e.trail = append(e.trail, err.Trail()...)
	if e.StatusCode() < err.StatusCode() {
		e.statusCode = err.StatusCode()
	}

	e.outgoingMessage = strings.Join([]string{e.outgoingMessage, err.Error()}, "; ")
	return e
}

// StatusCode returns the status code of the error.
func (e defaultError) StatusCode() int {
	return e.statusCode
}

// Error returns the error message.
func (e defaultError) Error() string {
	return e.outgoingMessage
}

// Unwrap will return the underlying error.
func (e defaultError) Unwrap() error {
	for _, err := range e.trail {
		return err.Error
	}

	return nil
}

// Is checks if the error is of the given type.
func (e defaultError) Is(target error) bool {
	for _, err := range e.trail {
		if errors.Is(err.Error, target) {
			return true
		}
	}

	return false
}

// As checks if the error can be cast as this type.
func (e defaultError) As(target any) bool {
	for _, err := range e.trail {
		if errors.As(err.Error, target) {
			return true
		}
	}

	return false
}

// NewError will return a default error struct.
func NewError(err error, statusCode int, outgoingMessage string) Error {
	// Create err if it is nil.
	if err == nil {
		err = errors.New(outgoingMessage)
	}

	strategyName := "unknown"
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		strategyName = details.Name()
	}
	return defaultError{
		trail: []ErrorWithTrail{
			{
				StrategyName: strategyName,
				Error:        err,
			},
		},
		statusCode:      statusCode,
		outgoingMessage: outgoingMessage,
	}
}
