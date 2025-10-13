package main

import (
	"fmt"
	"pwd-manager-tui/auth"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if !auth.ValidTokenExists() {
		loginModel := auth.NewLoginModel()
		loginScreen := tea.NewProgram(loginModel, tea.WithAltScreen())
		loginScreen.Run()
		loginRes, err := auth.Login(loginModel.GetValues()[0], loginModel.GetValues()[1])
		if err != nil {
			fmt.Printf("Could not login: %v", err)
		}
		fmt.Println(loginRes)

	} else {
		fmt.Println("Already Authorized")
	}
}
