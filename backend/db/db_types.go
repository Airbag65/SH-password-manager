package db

import (
	"database/sql"
	"fmt"
)

type User struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Id              int    `json:"id"`
	AuthToken       string `json:"auth_token"`
	TokenExpiryDate int64  `json:"token_expiry_date"`
}

func (user *User) ToString() string {
	return fmt.Sprintf("Email: %s\nPassword: %s\nName: %s\nSurname: %s\nId: %d\nAuthToken: %s\nTokenExpiryDate: %d",
		user.Email,
		user.Password,
		user.Name,
		user.Surname,
		user.Id,
		user.AuthToken,
		user.TokenExpiryDate)
}

func DbEntryToUser(row *sql.Rows) *User {
	selectedUser := &User{}
	for row.Next() {
		row.Scan(&selectedUser.Id, &selectedUser.Email, &selectedUser.Password, &selectedUser.Name, &selectedUser.Surname, &selectedUser.AuthToken, &selectedUser.TokenExpiryDate)
	}
	if selectedUser.Name == "" {
		return nil
	}
	return selectedUser
}

func DbEntryToHostNames(rows *sql.Rows) []string {
	hostNames := []string{}
	for rows.Next() {
		name := new(string)
		rows.Scan(&name)
		hostNames = append(hostNames, *name)
	}
	return hostNames
}

func DbEntryToPassword(row *sql.Rows) string {
	var password string
	for row.Next() {
		row.Scan(&password)
	}

	return password
}
