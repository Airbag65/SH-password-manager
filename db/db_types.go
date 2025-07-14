package db

import "fmt"

type User struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Id        int    `json:"id"`
	AuthToken string `json:"auth_token"`
}

func (user *User) ToString() string {
    return fmt.Sprintf("Email: %s\nPassword: %s\nName: %s\nSurname: %s\nId: %d\nAuthToken: %s", user.Email, user.Password, user.Name, user.Surname, user.Id, user.AuthToken)
}
