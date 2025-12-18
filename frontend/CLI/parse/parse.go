package parse

import (
	"fmt"
)


type ParsedCommand struct {
	Command   string
	Option    *string
	Parameter *string
}

type command struct {
	CommandName string
	CommandOptions []string
}

type Parser struct{
	commands []command
}


func New() *Parser {
	return &Parser{}
}

type commandOpt func(*command)


func defaultCommand (cn string) command {
	return command{
		CommandName: cn,
		CommandOptions: []string{},
	}	
}

func AddCommandOption(on string) commandOpt {
	return func(c *command) {
		c.CommandOptions = append(c.CommandOptions, on)
	}
}

func (p *Parser) AddCommand(command string, opts ...commandOpt) error {
	newCommand := defaultCommand(command)	
	for _, fn := range opts {
		fn(&newCommand)
	}

	p.commands = append(p.commands, newCommand)
	
	return nil
}

func (p *Parser) Parse(c []string) (*ParsedCommand, error) {
	if len(c) < 2 {
		return nil, &ParseError{
			What: NotEnoughArguments,
		}
	}

	fmt.Printf("%+v\n", p.commands)
	
	return &ParsedCommand{
		Command: "status",
		Option: nil,
		Parameter: nil,
	}, nil
}

