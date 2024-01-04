package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func sqlite() {
	os.Remove("sqlite-database.sqlite3") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating sqlite-database.sqlite3...")
	file, err := os.Create("sqlite-database.sqlite3") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.sqlite3 created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.sqlite3") // Open the created SQLite File
	defer sqliteDatabase.Close()                                          // Defer Closing the database
	createTable(sqliteDatabase)                                           // Create Database Tables

	// INSERT RECORDS
	insertStock(sqliteDatabase, "TESLA")
	insertStock(sqliteDatabase, "Microsoft")

	// DISPLAY INSERTED RECORDS
	displayStocks(sqliteDatabase)
}

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE IF NOT EXISTS stocks (
		stock_code INTEGER PRIMARY KEY,		
		stock_name TEXT NOT NULL	
	  );` // SQL Statement for Create Table

	log.Println("Create stock table...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("stock table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertStock(db *sql.DB, name string) {
	log.Println("Inserting stock record ...")
	insertStockSQL := `INSERT INTO stocks (stock_name) VALUES (?)`
	statement, err := db.Prepare(insertStockSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(name)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayStocks(db *sql.DB) {
	row, err := db.Query("SELECT * FROM stocks ORDER BY stock_code")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var code int
		var name string
		row.Scan(&code, &name)
		log.Println("Stock: ", code, " ", name)
	}
}
