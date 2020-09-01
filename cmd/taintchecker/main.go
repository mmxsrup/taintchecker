package main

import (
	"taintchecker"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(taintchecker.Analyzer) }

