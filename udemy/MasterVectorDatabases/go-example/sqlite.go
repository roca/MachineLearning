package main

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
	"log"
	"unsafe"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

var sqliteDatabase *sql.DB

func sqlite(vectors ...[]float64) {
	// os.Remove("vector.sqlite3") // I delete the file to avoid duplicated records.
	// // SQLite is a file based database.

	// log.Println("Creating vector.sqlite3...")
	// file, err := os.Create("vector.sqlite3") // Create SQLite file
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// file.Close()
	// log.Println("vector.sqlite3 created")

	sqliteDatabase, _ = sql.Open("sqlite3", "vector2.sqlite3") // Open the created SQLite File
	// Defer Closing the database

	// installVss()  // Install VSS
	createTable() // Create Database Tables
	clearTable()  // Clear Database Tables

	// INSERT RECORDS
	for _, vector := range vectors {
		insertVector(vector)
	}

	// DISPLAY INSERTED RECORDS
	displayVectors()
}

func installVss() {
	// stmt := `
	// .load ./vector0
	// .load ./vss0
	// create virtual table vss_vectors using vss0( vector_embedding(384) );
	// `

	// log.Println("Loading vss...")
	// statement, err := sqliteDatabase.Prepare(stmt) // Prepare SQL Statement
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// _, err = statement.Exec() // Execute SQL Statements
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// log.Println("vss loaded")

	var version, vector string
	err := sqliteDatabase.QueryRow("SELECT vss_version(), vector_to_json(?)", []byte{0x00, 0x00, 0x28, 0x42}).Scan(&version, &vector)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("version=%s vector=%s\n", version, vector)
}

func clearTable() {
	clearTableSQL := `DELETE FROM vectors;`
	log.Println("Clearing vector table...")
	statement, err := sqliteDatabase.Prepare(clearTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("vector table cleared")
}

func createTable() {
	createStudentTableSQL := `
	-- CREATE TABLE IF NOT EXISTS vectors (
	--	vector_id INTEGER PRIMARY KEY,		
	--	vector BLOB NOT NULL	
	--);
	DELETE FROM vectors;
	` // SQL Statement for Create Table

	log.Println("Create vector table...")
	statement, err := sqliteDatabase.Prepare(createStudentTableSQL) // Prepare SQL Statement
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

	l := len(src) / 8
	ptr := unsafe.Pointer(&src[0])
	// It is important to keep in mind that the Go garbage collector
	// will not interact with this data, and that if src if freed,
	// the behavior of any Go code using the slice is nondeterministic.
	// Reference: https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
	return (*[1 << 26]float64)((*[1 << 26]float64)(ptr))[:l:l]
}

// We are passing db reference connection from main to our method with other parameters
func insertVector(vector []float64) {
	log.Println("Inserting vector record ...")
	insertVectorSQL := `INSERT INTO vectors (vector) VALUES (?)`
	statement, err := sqliteDatabase.Prepare(insertVectorSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(encode(vector))
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayVectors() {
	row, err := sqliteDatabase.Query("SELECT * FROM vectors ORDER BY vector_id")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var vector []byte
		row.Scan(&id, &vector)
		log.Println("Vector: ", id, " length", len(decode(vector)))
	}
}

func findVectors(vector []float64) [][]float64 {
	var vectors [][]float64
	row, err := sqliteDatabase.Query("SELECT * FROM vectors WHERE vector == ?", encode(vector))
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var vector []byte
		row.Scan(&id, &vector)
		log.Println("Vector: ", id, " length", len(decode(vector)))
		vectors = append(vectors, decode(vector))
	}
	return vectors
}
