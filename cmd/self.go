package main

import (
	"flag"
	"fmt"

	"github.com/josephburnett/self/pkg/file"
)

var (
	filename = flag.String("filename", "", "Path to single file.")
	filedb   = flag.String("filedb", "", "Path to file database.")
)

func main() {
	flag.Parse()
	db, err := file.NewFileDb(*filedb)
	if err != nil {
		panic(err)
	}
	ids, err := db.ListNotes()
	if err != nil {
		panic(err)
	}
	for _, id := range ids[:20] {
		note, err := db.GetNote(id)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v\n", note.Title)
	}
}
