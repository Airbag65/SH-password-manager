package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	args := os.Args

	if len(args) > 1 {
		switch args[1] {
		case "l", "login":
			loggedIn := Login(reader)
			if loggedIn {
				fmt.Println("You are now logged in")
			} else {
				fmt.Println("Login failed")
			}
		case "cls", "check-login-status":
			CheckLoginStatus()
		default:
			fmt.Println("Use program correctly instead you dumdum")

		}
	} else {
		fmt.Println("Use program correctly instead you dumdum")
	}
	// email, err := GetInput("Enter email", reader)
	// if err != nil {
	// 	panic("Fix the GetInput function immediatley")
	// }
	// fmt.Println(email)
	//
	// password, err := GetObscuredInput("Enter password")
	// if err != nil {
	// 	panic("Fix the GetObscuredInput function immediatley")
	// }
	// fmt.Println(password)

}
