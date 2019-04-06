package main

import (
	"flag"
	"fmt"

	"github.com/josephburnett/self/pkg/command"
	"github.com/josephburnett/self/pkg/db"
	"github.com/josephburnett/self/pkg/file"
	"github.com/josephburnett/self/pkg/sqlite"
)

var (
	filename    = flag.String("filename", "", "Path to single file.")
	filedb      = flag.String("filedb", "", "Path to file database.")
	sqlitedb    = flag.String("sqlitedb", "", "Path to sqlite database.")
	cmd         = flag.String("command", "", "Command [search, read, list].")
	interactive = flag.Bool("interactive", false, "Interactive shell.")
	output      = flag.String("output", "", "Output format [body, json].")
	limit       = flag.Int("limit", 20, "Limit list output.")
	id          = flag.String("id", "", "Note Id.")
	tags        = flag.String("tags", "", "Comma separated list of tags.")
	sub         = flag.String("sub", "", "Sub-string to search for.")
)

func main() {
	flag.Parse()
	var d db.Database
	var err error
	if *filedb != "" {
		d, err = file.NewFileDb(*filedb)
		if err != nil {
			panic(err)
		}
	}
	if *sqlitedb != "" {
		d, err = sqlite.NewSqliteDb(*sqlitedb)
		if err != nil {
			panic(err)
		}
	}
	if d == nil {
		panic("no database given")
	}
	if *cmd != "" {
		var err error
		switch *cmd {
		case "search":
			err = command.Search(d, *sub, *limit)
		case "read":
			err = command.Read(d, *id)
		case "list":
			err = command.List(d, *tags, *limit)
		default:
			panic(fmt.Sprintf("Unknown command: %q.", *cmd))
		}
		if err != nil {
			panic(err)
		}
	}
}
