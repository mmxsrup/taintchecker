package a

import (
	"os"
	"io/ioutil"
	"log"
	"database/sql"
)

func openFile() {
	filepath := "test.txt"
	os.Open(filepath) // OK

	homepath := "/home/"
	os.Open(homepath + filepath) // OK

	dirpath := "test/"
	os.Open(homepath + dirpath + filepath) // OK
}

func openTaintedFile() {
	filepath := os.Getenv("TAINTED_FILE_PATH")
	os.Open(filepath) // want "NG"

	homepath := "/home/"
	os.Open(homepath + filepath) // want "NG"

	dirpath := "test/"
	os.Open(homepath + dirpath + filepath) // want "NG"
}

// TODO: Should check if the argument is tainted.
func openFileFromArg(filepath string) {
	os.Open(filepath) // want "NG"
}

func readFile() {
	filepath := "test.txt"
	ioutil.ReadFile(filepath) // OK

	homepath := "/home/"
	ioutil.ReadFile(homepath + filepath) // OK

	dirpath := "test/"
	ioutil.ReadFile(homepath + dirpath + filepath) // OK
}

func readTaintedFile() {
	filepath := os.Getenv("TAINTED_FILE_PATH")
	ioutil.ReadFile(filepath) // want "NG"

	homepath := "/home/"
	ioutil.ReadFile(homepath + filepath) // want "NG"

	dirpath := "test/"
	ioutil.ReadFile(homepath + dirpath + filepath) // want "NG"
}

func sqlQuery() {
	query := "SELECT id, name FROM person"
	db, err := sql.Open("mysql", "user:password@tcp(host:port)/dbname")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Query(query) // OK
}

func sqlTaintedQuery() {
	query := os.Getenv("TAINTED_SQL_Query")
	db, err := sql.Open("mysql", "user:password@tcp(host:port)/dbname")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Query(query) // want "NG"
}

func sqlQueryRow() {
	query := "SELECT id, name FROM person"
	db, err := sql.Open("mysql", "user:password@tcp(host:port)/dbname")
	if err != nil {
		log.Fatal(err)
	}
	db.QueryRow(query) // OK
}

func sqlTaintedQueryRow() {
	query := os.Getenv("TAINTED_SQL_Query")
	db, err := sql.Open("mysql", "user:password@tcp(host:port)/dbname")
	if err != nil {
		log.Fatal(err)
	}
	db.QueryRow(query) // want "NG"
}
