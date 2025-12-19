package main

type Command interface {
	Execute() error
}

type StatusCommand struct{}

type LoginCommand struct{}

type SignOutCommand struct{}

type SignUpCommand struct{}

type AddCommand struct{}

type ListCommand struct{}

type GetCommand struct {
	FlagExists bool
	FlagValue  string
}

type RemoveCommand struct {
	FlagExists bool
	FlagValue  string
}
