package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Airbag65/argparse"
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

	p := argparse.New()
	hostDesc := "Specify which host to direct the command at"
	hostFlag := argparse.NewFlag("--host", hostDesc, true)
	hFlag := argparse.NewFlag("-h", hostDesc, false)
	for _, comm := range commands {
		switch comm {
		case "get", "remove", "rm":
			err := p.AddCommand(comm, argparse.AddFlag(hostFlag), argparse.AddFlag(hFlag))
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

	command, err := p.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", command)
}
