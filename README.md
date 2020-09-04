# Taint Checker
The tool statically analyzes the source code written in Golang and
perform [taint checking](https://en.wikipedia.org/wiki/Taint_checking).
The purpose of this tool is to find security issues.

## Install
```
go get -u git@github.com:mmxsrup/taintchecker.git
```

## Taint definition
We consider only constants that are embedded in the program to be safe.
In other words, we assume that all values other than the constant are taint values.

## Check item
* Detect file access using tainted file path. 
  * os.Open()
  * io/ioutil.ReadFile()
* Detect sql operations with tainted queries.
  * database/sql.Query()
  * database/sql.QueryRow()
