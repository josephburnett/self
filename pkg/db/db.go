package db

import "time"

type Id string
type Tag string

type Note struct {
	Id      Id
	Title   string
	Body    string
	Tags    []Tag
	Created time.Time
	Updated time.Time
}

type Database interface {
	ListNotes() ([]Id, error)
	GetNote(id Id) (*Note, error)
	PutNote(note *Note) error
	DeleteNote(id Id) error

	ListTags() ([]Tag, error)
	GetTaggedNotes(tag Tag) ([]*Note, error)

	Repair() (bool, []error)
}
