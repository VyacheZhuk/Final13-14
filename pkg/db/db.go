package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(255) NOT NULL DEFAULT ""
);

CREATE INDEX idx_date ON scheduler(date);
`

var db *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if install {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("failed to create schema: %v", err)
		}
		fmt.Println("Database schema created successfully")
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
