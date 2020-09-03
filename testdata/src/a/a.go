package a

import (
	"fmt"
	"os"
	"io/ioutil"
	"log"
)

func readFile() {
	file := "/tmp/a"
	content, err := ioutil.ReadFile(file) // OK
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))
}

func readTaintedFile() {
	file := os.Getenv("TAINTED_FILE_PATH")
	fmt.Println(file)
	content, err := ioutil.ReadFile(file) // want "NG"
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))
}
