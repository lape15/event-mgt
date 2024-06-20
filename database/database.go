package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Database struct {
	db *sql.DB
}

func (db Database) CloseDb() {
	db.db.Close()
}

func (db *Database) QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	row := db.db.QueryRow(query, args...)
	return row, nil
}

func (db *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.db.Query(query, args...)
	return rows, err
}

func (db *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := db.db.Exec(query, args...)
	return result, err
}

func (d *Database) SetDB(db *sql.DB) {
	d.db = db
}

var Db Database

func ConnectDb() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	var cfg = mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
		Net:                  "tcp",
		Addr:                 "sql11.freesqldatabase.com",
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
	}
	var err error
	Db.db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	pingErr := Db.db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}
