package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/josephburnett/self/pkg/db"
)

func Insert(d db.Database, title, body string) error {
	if title == "" {
		return fmt.Errorf("Title is required")
	}
	if body == "" {
		return fmt.Errorf("Body is required")
	}
	now := time.Now()
	n := &db.Note{
		Id:      db.NewId(),
		Title:   title,
		Body:    body,
		Tags:    []db.Tag{},
		Created: now,
		Updated: now,
	}
	return d.PutNote(n)
}

func Search(d db.Database, sub string, limit int) error {
	ids, err := d.ListNotes()
	if err != nil {
		return err
	}
	if sub == "" {
		return fmt.Errorf("Search requires sub.")
	}
	subString := strings.ToLower(sub)
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
			return err
		}
		if filter(note) {
			notes = append(notes, note)
		}
		if len(notes) == limit {
			more = true
			break
		}
	}
	for _, note := range notes {
		fmt.Printf("=== %v ===\n%v %v\n\n", note.Title, note.Id, note.Tags)
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
	return nil
}

func Read(d db.Database, id string) error {
	if id == "" {
		fmt.Errorf("Id is required for read command.")
	}
	note, err := d.GetNote(db.Id(id))
	if err != nil {
		return err
	}
	fmt.Printf("=== %v ===\n\n%v\n\n", note.Title, note.Body)
	return nil
}

func List(d db.Database, tags string, limit int) error {
	ids, err := d.ListNotes()
	if err != nil {
		return err
	}
	if len(ids) < limit {
		limit = len(ids)
	}
	more := false
	if tags == "" {
		for _, id := range ids[:limit] {
			note, err := d.GetNote(id)
			if err != nil {
				return err
			}
			fmt.Printf("=== %v ===\n%v %v\n\n", note.Title, note.Id, note.Tags)
		}
		if len(ids) > limit {
			more = true
		}
	} else {
		notes := make([]*db.Note, 0)
		ts := strings.Split(tags, ",")
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
				return err
			}
			if filter(note) {
				notes = append(notes, note)
			}
			if len(notes) == limit {
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
	return nil
}
