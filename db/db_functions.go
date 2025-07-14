package db

import (
	"fmt"
	"log"
)

func GetUser(userEmail string) *User {
	database := CreateConnection()
	if database == nil {
		log.Fatal("Could not connect to './db/database.db'")
		return nil
	}
	defer database.Close()

	row, err := database.Query(fmt.Sprintf("SELECT * FROM user WHERE email = '%s'", userEmail))
	if err != nil {
		log.Fatalf("Could not retrieve user with email '%s'", userEmail)
		return nil
	}
	defer row.Close()
    selectedUser := &User{}
    for row.Next(){
        row.Scan(&selectedUser.Id, &selectedUser.Email, &selectedUser.Password, &selectedUser.Name, &selectedUser.Surname, &selectedUser.AuthToken)
    }
	return selectedUser
}
