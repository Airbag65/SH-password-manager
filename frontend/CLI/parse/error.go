package parse

import (
	"fmt"
)

type ErrorType int

const (
	NotEnoughArguments ErrorType = iota
	InvalidCommand
	AlreadyAdded
)

type ParseError struct {
	What ErrorType
}

func (e *ParseError) Error() string {
	var message string
	switch e.What {
	case NotEnoughArguments:
		message = "Not enough arguments"
	case InvalidCommand:
		message = "Invalid Command"
	case AlreadyAdded:
		message = "Command already added"
	}
	return fmt.Sprintf("ParseError: %s", message)
}
