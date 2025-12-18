package main

import (
	"fmt"
	"os"
	"passport-cli/parse"
)

func main() {
	commands := []string{
		"status",
		"login",
		"signout",
		"signup",
		"add",
		"get",
		"list",
		"ls",
		"remove",
		"rm",
	}

	p := parse.New()
	hostDesc := "Specify which host to direct the command at"
	hostFlag := parse.NewFlag("--host", hostDesc, true)
	hFlag := parse.NewFlag("-h", hostDesc, false)
	for _, comm := range commands {
		switch comm {
		case "get", "remove", "rm":
			err := p.AddCommand(comm, parse.AddFlag(hostFlag), parse.AddFlag(hFlag))
			if err != nil {
				fmt.Println(err)
			}
		default:
			err := p.AddCommand(comm)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	_, err := p.Parse(os.Args)
	if err != nil {
		fmt.Println(err)
	}

}
