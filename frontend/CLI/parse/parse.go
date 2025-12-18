package parse

func AddFlag(f *Flag) commandOpt {
	return func(c *command) {
		for _, flag := range c.CommandOptions {
			if f.Flag == flag.Flag {
				return
			}
		}
		c.CommandOptions = append(c.CommandOptions, *f)
	}
}

func (p *Parser) AddCommand(command string, opts ...commandOpt) error {
	for _, comm := range p.commands {
		if comm.CommandName == command {
			return &ParseError{
				What: AlreadyAdded,
			}
		}
	}
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
	result := &ParsedCommand{
		Command: "",
		Option: "",
		Parameter: "",
	}
	if !p.isValidCommand(c[1]) {
		return nil, &ParseError{
			What: InvalidCommand,
		}
	}
	result.Command = c[1]
	if !p.commandHasFlags(c[1]) {
		return result, nil
	}
	if err := p.validateFlagUse(c); err != nil {
		return nil, err
	}
	flags, err := p.getFlags(result.Command, c[2:])
	if err != nil {
		return nil, err
	}

	result.Option = flags[0]
	result.Parameter = flags[1]
	return result, nil
}

func (p *Parser) isValidCommand(comm string) bool {
	for _, commmand := range p.commands {
		if comm == commmand.CommandName {
			return true
		}
	}
	return false
}
