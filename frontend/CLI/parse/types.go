package parse

type ParsedCommand struct {
	Command   string
	Option    string
	Parameter string
}

type command struct {
	CommandName    string
	CommandOptions []Flag
}

type Parser struct {
	commands []command
}

type commandOpt func(*command)

type Flag struct {
	Flag        string
	Description string
	Forced      bool
}

func NewFlag(flag, desc string, must bool) *Flag {
	return &Flag{
		Flag:        flag,
		Description: desc,
		Forced:      must,
	}
}

func New() *Parser {
	return &Parser{}
}

func defaultCommand(cn string) command {
	return command{
		CommandName:    cn,
		CommandOptions: []Flag{},
	}
}
