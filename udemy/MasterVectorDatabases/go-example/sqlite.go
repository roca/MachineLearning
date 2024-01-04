package main

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"log"
	"os"
	"unsafe"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func sqlite() {
	os.Remove("vector.sqlite3") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating vector.sqlite3...")
	file, err := os.Create("vector.sqlite3") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("vector.sqlite3 created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./vector.sqlite3") // Open the created SQLite File
	defer sqliteDatabase.Close()                                 // Defer Closing the database
	createTable(sqliteDatabase)                                  // Create Database Tables

	// INSERT RECORDS
	insertVector(sqliteDatabase, []float64{1.3, 3.5, 2.2, 0.9, 0.6, 1.5})
	insertVector(sqliteDatabase, []float64{2.8, 1.6, 3.8, 2.2, 3.2, 17.0})

	// DISPLAY INSERTED RECORDS
	displayVectors(sqliteDatabase)
}

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE IF NOT EXISTS vectors (
		vector_id INTEGER PRIMARY KEY,		
		vector BLOB NOT NULL	
	  );` // SQL Statement for Create Table

	log.Println("Create vector table...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("vector table created")
}

func encode(fs []float64) []byte {
	//buf := make([]byte, len(fs)*8)
	buf := new(bytes.Buffer)
	for _, f := range fs {
		binary.Write(buf, binary.LittleEndian, f)
	}
	return buf.Bytes()
}

func decode(src []byte) []float64 {
	if len(src) == 0 {
		return nil
	}

	l := len(src) / 4
	ptr := unsafe.Pointer(&src[0])
	// It is important to keep in mind that the Go garbage collector
	// will not interact with this data, and that if src if freed,
	// the behavior of any Go code using the slice is nondeterministic.
	// Reference: https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
	return (*[1 << 26]float64)((*[1 << 26]float64)(ptr))[:l:l]
}

// We are passing db reference connection from main to our method with other parameters
func insertVector(db *sql.DB, vector []float64) {
	log.Println("Inserting vector record ...")
	insertVectorSQL := `INSERT INTO vectors (vector) VALUES (?)`
	statement, err := db.Prepare(insertVectorSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(encode(vector))
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayVectors(db *sql.DB) {
	row, err := db.Query("SELECT * FROM vectors ORDER BY vector_id")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var vector []byte
		row.Scan(&id, &vector)
		log.Println("Vector: ", id, " ", decode(vector))
	}
}
