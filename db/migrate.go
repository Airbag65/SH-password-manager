package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func CreateConnection() *sql.DB {
	database, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		log.Fatal("Could not connect to database")
		return nil
	}
	log.Println("Connected to './db/database.db'")
	return database
}

func Migrate() {
	log.Println("Migrating database")
	os.Remove("./db/database.db")

	file, err := os.Create("./db/database.db")
	if err != nil {
		return
	}
	file.Close()
	database := CreateConnection()
	if database == nil {
		return
	}
	defer database.Close()
	createTables(database)
    insertDefault(database)
	log.Println("Migration finished")
}

func insertDefault(db *sql.DB) {
    InsertDefaultUserQuery := `INSERT INTO user(email, password, name, surname, auth_token) VALUES(?,?,?,?,?);`

    log.Println("Inserting default user")
    statement, err := db.Prepare(InsertDefaultUserQuery)
    if err != nil {
        log.Fatalf("Error inserting default user: %s", err)
        return
    }
    statement.Exec("normananton03@gmail.com", "VerySecurePassword", "Anton", "Norman", uuid.New().String())
}

func createTables(db *sql.DB) {
	CreateUsersTablesQuery := `CREATE TABLE user (
        "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"email" TEXT NOT NULL,
        "password" TEXT NOT NULL,
        "name" TEXT NOT NULL,
        "surname" TEXT NOT NULL,
        "auth_token" TEXT
    );`
	CreatePasswordsTableQuery := `CREATE TABLE password (
        "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        "userId" integer NOT NULL,
        "password" TEXT NOT NULL,
        "goesTo" TEXT NOT NULL,
        FOREIGN KEY(userId) REFERENCES user(id)
    );`

	log.Println("Creating 'user' tables")

	statement, err := db.Prepare(CreateUsersTablesQuery)
	if err != nil {
		log.Fatalf("Error creating 'user' table: %s", err)
		return
	}
	statement.Exec()

	log.Println("Creating 'password' tables")

	statement, err = db.Prepare(CreatePasswordsTableQuery)
	if err != nil {
		log.Fatalf("Error creating 'password' table: %s", err)
		return
	}
	statement.Exec()
    log.Println("All tables created")

}
