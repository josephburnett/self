package main

import (
	"flag"
)

var (
	filename = flag.String("filename", "", "Path to single file.")
	filedb   = flag.String("filedb", "", "Path to file database.")
)

func main() {
	flag.Parse()
}
