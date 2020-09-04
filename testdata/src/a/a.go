package a

import (
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"database/sql"
)

func openFile() {
	filepath := "test.txt"
	file, err := os.Open(filepath) // OK
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func openTaintedFile() {
	filepath := os.Getenv("TAINTED_FILE_PATH")
	file, err := os.Open(filepath) // want "NG"
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func readFile() {
	filepath := "test.txt"
	content, err := ioutil.ReadFile(filepath) // OK
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))
}

func readTaintedFile() {
	filepath := os.Getenv("TAINTED_FILE_PATH")
	content, err := ioutil.ReadFile(filepath) // want "NG"
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))
}

func sqlQuery() {
	query := "SELECT id, name FROM person"
	db, err := sql.Open("mysql", "user:password@tcp(host:port)/dbname")
	rows, err := db.Query(query) // OK
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}

func sqlTaintedQuery() {
	query := os.Getenv("TAINTED_SQL_Query")
	db, err := sql.Open("mysql", "user:password@tcp(host:port)/dbname")
	rows, err := db.Query(query) // want "NG"
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}
