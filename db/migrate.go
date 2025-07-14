package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func runMigration() {
    log.Println("Migrating database")
	os.Remove("./database.db")

	file, err := os.Create("./database.db")
	if err != nil {
		return
	}
	file.Close()
	database, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return
	}
	defer database.Close()

    createTables(database)
}

func createTables(db *sql.DB) {
	CreateTablesQuery := `CREATE TABLE user (
		"email" TEXT NOT NULL,
        "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
        "password" TEXT NOT NULL,
        "name" TEXT NOT NULL,
        "surname" TEXT NOT NULL
    );
    CREATE TABLE password (
        "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        "userId" integer NOT NULL,
        "password": TEXT NOT NULL,
        "goesTo": TEXT NOT NULL,
        FOREIGN KEY(userId) REFERENCES user(id)
    );`
    log.Println("Creating tables")
    
    statement, err := db.Prepare(CreateTablesQuery)
    if err != nil {
        log.Fatalf("Error creating tables: %s", err)
        return
    }
    statement.Exec()
    log.Println("Tables created")
}
