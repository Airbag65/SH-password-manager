package db

import (
	"fmt"
	"log"
	"time"
)

func GetUserWithEmail(userEmail string) *User {
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
	for row.Next() {
		row.Scan(&selectedUser.Id, &selectedUser.Email, &selectedUser.Password, &selectedUser.Name, &selectedUser.Surname, &selectedUser.AuthToken, &selectedUser.TokenExpiryDate)
	}
	return selectedUser
}

func SetAuthToken(userEmail, authToken string) {
    database := CreateConnection()
    if database == nil {
		log.Fatal("Could not connect to './db/database.db'")
		return
	}
	defer database.Close()

    expiryDate := time.Now().AddDate(0, 1, 0).Unix()
    updateUserQuery := fmt.Sprintf(`UPDATE user
    SET auth_token = '%s',
        token_expiry_date = %d
    WHERE 
    email = '%s'`, authToken, expiryDate, userEmail)
    statement, err := database.Prepare(updateUserQuery)
    if err != nil {
        log.Fatalf("Could not update AUTH on '%s'", userEmail)
        return
    }
    statement.Exec()
    log.Printf("Updated AUTH on '%s'", userEmail)
}
