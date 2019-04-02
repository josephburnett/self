package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/josephburnett/self/pkg/db"
	"github.com/josephburnett/self/pkg/file"
)

var (
	filename    = flag.String("filename", "", "Path to single file.")
	filedb      = flag.String("filedb", "", "Path to file database.")
	command     = flag.String("command", "", "Command [search, read, list].")
	interactive = flag.Bool("interactive", false, "Interactive shell.")
	output      = flag.String("output", "", "Output format [body, json].")
	limit       = flag.Int("limit", 20, "Limit list output.")
	id          = flag.String("id", "", "Note Id.")
	tags        = flag.String("tags", "", "Comma separated list of tags.")
	sub         = flag.String("sub", "", "Sub-string to search for.")
)

func main() {
	flag.Parse()
	d, err := file.NewFileDb(*filedb)
	if err != nil {
		panic(err)
	}
	if *command != "" {
		switch *command {
		case "search":
			search(d)
		case "read":
			read(d)
		case "list":
			list(d)
		default:
			panic(fmt.Sprintf("Unknown command: %q.", *command))
		}
	}
}

func search(d db.Database) {
	ids, err := d.ListNotes()
	if err != nil {
		panic(err)
	}
	if *sub == "" {
		panic("Search requires sub.")
	}
	subString := strings.ToLower(*sub)
	filter := func(note *db.Note) bool {
		lowerTitle := strings.ToLower(note.Title)
		if strings.Contains(lowerTitle, subString) {
			return true
		}
		lowerBody := strings.ToLower(note.Body)
		if strings.Contains(lowerBody, subString) {
			return true
		}
		return false
	}
	notes := make([]*db.Note, 0)
	more := false
	for _, id := range ids {
		note, err := d.GetNote(id)
		if err != nil {
			panic(err)
		}
		if filter(note) {
			notes = append(notes, note)
		}
		if len(notes) == *limit {
			more = true
			break
		}
	}
	for _, note := range notes {
		fmt.Printf("=== %v ===\n\n%v %v\n", note.Title, note.Id, note.Tags)
		lines := strings.Split(note.Body, "\n")
		for _, line := range lines {
			lowerLine := strings.ToLower(line)
			if strings.Contains(lowerLine, subString) {
				fmt.Printf("%v\n", line)
			}
		}
		fmt.Printf("\n")
	}
	if more {
		fmt.Printf("...\n\n")
	}
}

func read(d db.Database) {
	if *id == "" {
		panic(fmt.Sprintf("Id is required for read command."))
	}
	note, err := d.GetNote(db.Id(*id))
	if err != nil {
		panic(err)
	}
	fmt.Printf("=== %v ===\n\n%v\n\n", note.Title, note.Body)
}

func list(d db.Database) {
	ids, err := d.ListNotes()
	if err != nil {
		panic(err)
	}
	more := false
	if *tags == "" {
		for _, id := range ids[:*limit] {
			note, err := d.GetNote(id)
			if err != nil {
				panic(err)
			}
			fmt.Printf("=== %v ===\n%v %v\n\n", note.Title, note.Id, note.Tags)
		}
		if len(ids) > *limit {
			more = true
		}
	} else {
		notes := make([]*db.Note, 0)
		ts := strings.Split(*tags, ",")
		filter := func(note *db.Note) bool {
			for _, wantTag := range ts {
				has := false
				for _, hasTag := range note.Tags {
					if string(hasTag) == wantTag {
						has = true
					}
				}
				if has == false {
					return false
				}
			}
			return true
		}
		for _, id := range ids {
			note, err := d.GetNote(id)
			if err != nil {
				panic(err)
			}
			if filter(note) {
				notes = append(notes, note)
			}
			if len(notes) == *limit {
				more = true
				break
			}
		}
		for _, note := range notes {
			fmt.Printf("=== %v ===\n%v %v\n\n", note.Title, note.Id, note.Tags)
		}
	}
	if more {
		fmt.Printf("...\n\n")
	}
}
