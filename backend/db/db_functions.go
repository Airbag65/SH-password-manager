package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Store struct {
	db *sql.DB
}

func (s *Store) Init() error {
	database, err := sql.Open("sqlite3", "./db/.db")
	if err != nil {
		log.Fatal("Could not connect to database")
		return err
	}
	log.Println("Connected to './db/.db'")
	s.db = database
	return nil
}

func (s *Store) GetUserWithEmail(userEmail string) *User {
	row, err := s.db.Query(fmt.Sprintf("SELECT * FROM user WHERE email = '%s';", userEmail))
	if err != nil {
		log.Fatalf("Could not retrieve user with email '%s'", userEmail)
		return nil
	}

	defer func() {
		if err := row.Close(); err != nil {
			log.Fatal("Could not close database row")
		}
	}()

	user := DbEntryToUser(row)
	if user == nil {
		log.Println("User not found(GetUserWithEmail)")
	}
	return user
}

func (s *Store) GetUserWithAuthToken(authToken string) *User {
	row, err := s.db.Query(fmt.Sprintf("SELECT * FROM user where auth_token = '%s';", authToken))
	if err != nil {
		log.Fatalf("Could not retrieve user with auth_token '%s'", authToken)
		return nil
	}

	defer func() {
		if err := row.Close(); err != nil {
			log.Fatal("Could not close database row")
		}
	}()

	selectedUser := DbEntryToUser(row)
	if selectedUser == nil {
		log.Println("User not found(GetUserWithAuthToken)")
		return nil
	}

	if selectedUser.TokenExpiryDate <= time.Now().Unix() {
		log.Println("auth_token not valid")
		s.RemoveAuthToken(selectedUser.Email)
		return nil
	}

	return selectedUser
}

func (s *Store) SetAuthToken(userEmail, authToken string) {
	expiryDate := time.Now().AddDate(0, 1, 0).Unix()
	updateUserQuery := fmt.Sprintf(`UPDATE user
    SET auth_token = '%s',
        token_expiry_date = %d
    WHERE 
    email = '%s';`, authToken, expiryDate, userEmail)
	statement, err := s.db.Prepare(updateUserQuery)
	if err != nil {
		log.Fatalf("Could not update AUTH on '%s'(SetAuthToken)", userEmail)
		return
	}
	statement.Exec()
	log.Printf("Updated AUTH on '%s'", userEmail)
}

func (s *Store) RemoveAuthToken(userEmail string) {
	user := s.GetUserWithEmail(userEmail)
	if user == nil {
		log.Fatal("User not found(RemoveAuthToken)")
		return
	}
	removeTokenQuery := fmt.Sprintf(`UPDATE user
    SET auth_token = '',
        token_expiry_date = 0
    WHERE
    email = '%s';`, userEmail)

	statement, err := s.db.Prepare(removeTokenQuery)
	if err != nil {
		log.Fatalf("Could not remove AUTH token on '%s'", userEmail)
		return
	}
	statement.Exec()
	log.Printf("Removed AUTH token on '%s'", userEmail)
}

func (s *Store) CreateNewUser(email, password, name, surname string) {
	createNewUserQuery := `INSERT INTO user(
        email, 
        password, 
        name, 
        surname,
        auth_token, 
        token_expiry_date) 
        VALUES(?,?,?,?,?,?);`
	log.Println("Inserting new user")
	statement, err := s.db.Prepare(createNewUserQuery)
	if err != nil {
		log.Fatalf("Error inserting new user '%s'", email)
		return
	}
	statement.Exec(email, password, name, surname, "", 0)
	log.Printf("Created user '%s %s' - '%s'", name, surname, email)
}

func (s *Store) AddNewPassord(userId int, password, hostName string) error {
	insertNewPasswordQuery := `INSERT INTO password(
		user_id,
		password,
		host_name)
		VALUES(?,?,?);`
	log.Printf("Inserting new password for user '%d'", userId)
	statement, err := s.db.Prepare(insertNewPasswordQuery)
	if err != nil {
		log.Fatalf("Error inserting new password for user '%d' - %v", userId, err)
		return err
	}

	_, err = statement.Exec(userId, password, hostName)
	if err != nil {
		return err
	}
	log.Printf("Inserted new password for user '%d'", userId)
	return nil
}

func (s *Store) GetHostNames(userId int) []string {
	getHostsQuery := fmt.Sprintf("SELECT host_name FROM password where user_id = %d", userId)

	row, err := s.db.Query(getHostsQuery)
	if err != nil {
		log.Printf("Could not get host names for user: %d - %v", userId, err)
		return []string{}
	}

	return DbEntryToHostNames(row)
}

func (s *Store) GetPassword(userId int, hostname string) (string, error) {
	getPasswordQuery := fmt.Sprintf(`SELECT password 
		FROM password 
		WHERE 
			user_id = %d
		AND 
			host_name = '%s';`, userId, hostname)
	
	row, err := s.db.Query(getPasswordQuery)
	if err != nil {
		log.Fatalf("Could not get host names for user: %d - %v", userId, err)
		return "", err
	}

	return DbEntryToPassword(row), nil
}

func (s *Store) RemovePassword(userId int, hostname string) error {
	removePasswordQuery := fmt.Sprintf(`REMOVE FROM password
		WHERE
			user_id = %d
		AND
			host_name = '%s';`, userId, hostname)
	statement, err := s.db.Prepare(removePasswordQuery)	
	if err != nil {
		log.Fatalf("Could not remove password (userId: %d - hostname: %s)", userId, hostname)
		return err
	}

	_, err = statement.Exec()
	return err
}
