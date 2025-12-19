package main

import (
	"log"
	"os"

	"github.com/Airbag65/argparse"
)

func main() {
	p, err := InitParser()
	if err != nil {
		log.Fatal(err)
	}
	command, err := p.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	c := CreateCommand(command)
	c.Execute()
}

func InitParser() (*argparse.Parser, error){
	commands := []string{"status", "login", "signout", "signup", "add", "get", "list", "ls", "remove", "rm"}

	p := argparse.New()
	hostDesc := "Specify which host to direct the command at"
	hostFlag := argparse.NewFlag("--host", hostDesc, true)
	hFlag := argparse.NewFlag("-h", hostDesc, false)
	for _, comm := range commands {
		switch comm {
		case "get", "remove", "rm":
			err := p.AddCommand(comm, argparse.AddFlag(hostFlag), argparse.AddFlag(hFlag))
			if err != nil {
				return nil, err
			}
		default:
			err := p.AddCommand(comm)
			if err != nil {
				return nil, err
			}
		}
	}
	return p, nil
}
