package parse

import (
	"fmt"
)

type ErrorType int

const (
	NotEnoughArguments ErrorType = iota
	InvalidCommand
	AlreadyAdded
	MissingFlags
	MissingValue
	NoSuchFlag
)

type ParseError struct {
	What       ErrorType
	WhichFlags []string
}

func (e *ParseError) Error() string {
	var message string
	switch e.What {
	case NotEnoughArguments:
		message = "Not enough arguments"
	case InvalidCommand:
		message = "Invalid command"
	case AlreadyAdded:
		message = "Command already added"
	case MissingFlags:
		message = fmt.Sprintf("Invalid or missing flags. Possible flags are: %+v", e.WhichFlags)
	case MissingValue:
		message = "Flag is missing value"
	case NoSuchFlag:
		message = fmt.Sprintf("No such flag (%s)", e.WhichFlags[0])
	}
	return fmt.Sprintf("ParseError: %s", message)
}
