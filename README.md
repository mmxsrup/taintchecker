# Taint Checker
The tool statically analyzes the source code written in Golang and
perform [taint checking](https://en.wikipedia.org/wiki/Taint_checking).
The purpose of this tool is to find security issues.

## Install
```
go get -u git@github.com:mmxsrup/taintchecker.git
```

## Check item
* Detect file access using tainted file path. 
