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

	row, err := database.Query(fmt.Sprintf("SELECT * FROM user WHERE email = '%s';", userEmail))
	if err != nil {
		log.Fatalf("Could not retrieve user with email '%s'", userEmail)
		return nil
	}
	defer row.Close()
	user := DbEntryToUser(row)
	if user == nil {
		log.Println("User not found(GetUserWithEmail)")
	}
	return user
}

func GetUserWithAuthToken(authToken string) *User {
	database := CreateConnection()
	if database == nil {
		log.Fatal("Could not connect to './db/database.db'")
		return nil
	}
	defer database.Close()

	row, err := database.Query(fmt.Sprintf("SELECT * FROM user where auth_token = '%s';", authToken))
	if err != nil {
		log.Fatalf("Could not retrieve user with auth_token '%s'", authToken)
		return nil
	}
	defer row.Close()

	selectedUser := DbEntryToUser(row)
	if selectedUser == nil {
		log.Println("User not found(GetUserWithAuthToken)")
		return nil
	}

	if selectedUser.TokenExpiryDate <= time.Now().Unix() {
		log.Println("auth_token not valid")
		RemoveAuthToken(selectedUser.Email)
		return nil
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
    email = '%s';`, authToken, expiryDate, userEmail)
	statement, err := database.Prepare(updateUserQuery)
	if err != nil {
		log.Fatalf("Could not update AUTH on '%s'(SetAuthToken)", userEmail)
		return
	}
	statement.Exec()
	log.Printf("Updated AUTH on '%s'", userEmail)
}

func RemoveAuthToken(userEmail string) {
	database := CreateConnection()
	if database == nil {
		log.Fatal("Could not connect to './db/database.db'")
		return
	}
	defer database.Close()

	user := GetUserWithEmail(userEmail)
	if user == nil {
		log.Fatal("User not found(RemoveAuthToken)")
		return
	}
	removeTokenQuery := fmt.Sprintf(`UPDATE user
    SET auth_token = '',
        token_expiry_date = 0
    WHERE
    email = '%s';`, userEmail)

	statement, err := database.Prepare(removeTokenQuery)
	if err != nil {
		log.Fatalf("Could not remove AUTH token on '%s'", userEmail)
		return
	}
	statement.Exec()
	log.Printf("Removed AUTH token on '%s'", userEmail)
}

func CreateNewUser(email, password, name, surname string) {
	database := CreateConnection()
	if database == nil {
		log.Fatal("Could not connect to './db/database.db'")
		return
	}
	defer database.Close()

	createNewUserQuery := `INSERT INTO user(
        email, 
        password, 
        name, 
        surname,
        auth_token, 
        token_expiry_date) 
        VALUES(?,?,?,?,?,?);`
	log.Println("Inserting new user")
	statement, err := database.Prepare(createNewUserQuery)
	if err != nil {
		log.Fatalf("Error inserting new user '%s'", email)
		return
	}
	statement.Exec(email, password, name, surname, "", 0)
	log.Printf("Created user '%s %s' - '%s'", name, surname, email)
}
