package main

import (
	"fmt"
	art "pwd-manager-tui/artistics"
	"pwd-manager-tui/auth"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Print("\033[H\033[2J")
	if !auth.ValidTokenExists() {
		startScreenModel := auth.NewStartScreenModel(new(int))
		startScreen := tea.NewProgram(startScreenModel, tea.WithAltScreen())
		startScreen.Run()
		switch startScreenModel.GetValue(){
		case 0:
			loginModel := auth.NewLoginModel()
			loginScreen := tea.NewProgram(loginModel, tea.WithAltScreen())
			loginScreen.Run()
			loginRes, err := auth.Login(loginModel.GetValues()[0], loginModel.GetValues()[1])
			if err != nil {
				fmt.Printf("Could not login: %v", err)
			}
			fmt.Printf("You are now logged in as %s %s\n", loginRes.Name, loginRes.Surname)
		case 1:
			signUpModel := auth.NewSignUpModel()
			signUpScreen := tea.NewProgram(signUpModel, tea.WithAltScreen())
			signUpScreen.Run()
		}
	} 
	fmt.Println("Authorized")

	fmt.Println(art.LoadTitle())
}
