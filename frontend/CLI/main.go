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
	for _, comm := range commands {
		switch comm {
		case "get", "remove", "rm":
			err := p.AddCommand(comm, parse.AddCommandOption("--host"), parse.AddCommandOption("-h"))
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
