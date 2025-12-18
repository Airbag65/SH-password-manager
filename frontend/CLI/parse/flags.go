package parse

import "slices"

func (p *Parser) commandHasFlags(c string) bool {
	for _, item := range p.commands {
		if item.CommandName == c {
			if len(item.CommandOptions) == 0 {
				return false
			}
		}
	}
	return true
}

func (p *Parser) commandHasForcedFlag(c string) (bool, []string) {
	forcedFlags := []string{}
	for _, item := range p.commands {
		hasFlags, flags := forcedFlagsHelper(item.CommandOptions)
		forcedFlags = flags
		if item.CommandName == c && hasFlags {
			return true, forcedFlags
		}
	}
	return false, nil
}

func forcedFlagsHelper(flags []Flag) (bool, []string) {
	fl := []string{}
	res := false
	for _, f := range flags {
		fl = append(fl, f.Flag)
		if f.Forced {
			res = true
		}
	}
	return res, fl
}

func (p *Parser) validateFlagUse(c []string) error {
	hasForced, flags := p.commandHasForcedFlag(c[1])
	if hasForced && len(c) < 3 {
		return &ParseError{
			What:       MissingFlags,
			WhichFlags: flags,
		}
	}
	if hasForced && len(c) < 4 {
		return &ParseError{
			What: MissingValue,
		}
	}
	return nil
}

func (c *command) getPossibleFlags() []string {
	flags := []string{}
	for _, flag := range c.CommandOptions {
		flags = append(flags, flag.Flag)
	}
	return flags
}

func (p *Parser) getFlags(command string, line []string) ([]string, error) {
	for _, com := range p.commands {
		if command == com.CommandName {
			if !slices.Contains(com.getPossibleFlags(), line[0]) {
				return nil, &ParseError{
					What: NoSuchFlag,
					WhichFlags: []string{line[0]},
				}
			}
		}
	}
	return line, nil	
}
