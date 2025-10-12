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
		fmt.Println(loginModel.GetValues())
	} else {
		fmt.Println("Already Authorized")
	}
}
