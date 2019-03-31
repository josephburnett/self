package file

import "github.com/josephburnett/self/pkg/db"

type fileDb struct {
	path string
}

func NewFileDb(path string) (Database, error) {
	return &fileDb{
		path: path,
	}
}

func (db *fileDb) ListNotes() ([]db.Id, error) {

}

func (db *fileDb) GetNote(id db.Id) (*db.Note, error) {

}

func (db *fileDb) PutNote(note *db.Note) error {

}

func (db *fileDb) DeleteNote(id db.Id) error {

}

func (db *fileDb) ListTags() ([]db.Tag, error) {

}

func (db *fileDb) GetTaggedNotes(tag db.Tag) ([]*Note, error) {

}

func (db *fileDb) Repair() (bool, []error) {

}
