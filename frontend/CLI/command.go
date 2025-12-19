package main

import (
	"fmt"

	"github.com/Airbag65/argparse"
)

func CreateCommand(pc *argparse.ParsedCommand) Command {
	switch pc.Command {
	case "status":
		return &StatusCommand{}
	case "login":
		return &LoginCommand{}
	case "signout":
		return &SignOutCommand{}
	case "signup":
		return &SignUpCommand{}
	case "add":
		return &AddCommand{}
	case "get":
		return &GetCommand{
			FlagExists: pc.Option != "",
			FlagValue: pc.Parameter,
		}
	case "list", "ls":
		return &ListCommand{}
	case "remove", "rm":
		return &RemoveCommand{
			FlagExists: pc.Option != "",
			FlagValue: pc.Parameter,
		}
	}
	return nil
}

func (c *StatusCommand) Execute() error {
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *LoginCommand) Execute() error {
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *SignOutCommand) Execute() error {
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *SignUpCommand) Execute() error {
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *AddCommand) Execute() error {
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *ListCommand) Execute() error {
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *GetCommand) Execute() error {
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *RemoveCommand) Execute() error {
	fmt.Printf("%+v\n", c)
	return nil
}
