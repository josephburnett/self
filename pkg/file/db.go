package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/josephburnett/self/pkg/db"
)

const notesDbDir = "notes-db"
const notesDir = "notes"
const notesExt = ".json"

type fileDb struct {
	notes string
}

func NewFileDb(path string) (db.Database, error) {
	_, file := filepath.Split(path)
	if file != notesDbDir {
		return nil, fmt.Errorf("Invalid filedb %q. Must be something %q.", path, notesDbDir)
	}
	fileInfo, err := os.Stat(filepath.Join(path, notesDir))
	if err != nil || !fileInfo.IsDir() {
		return nil, fmt.Errorf("Invalid filedb %q. Must contain %q directory.", path, notesDir)
	}
	return &fileDb{
		notes: filepath.Join(path, notesDir),
	}, nil
}

func (d *fileDb) ListNotes() ([]db.Id, error) {
	files, err := ioutil.ReadDir(d.notes)
	if err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool {
		return !files[i].ModTime().Before(files[j].ModTime())
	})
	ids := make([]db.Id, 0)
	for _, f := range files {
		name := f.Name()
		if filepath.Ext(name) == notesExt {
			ids = append(ids, db.Id(name[:len(name)-len(notesExt)]))
		}
	}
	return ids, nil
}

func (d *fileDb) GetNote(id db.Id) (*db.Note, error) {
	file, err := Load(filepath.Join(d.notes, string(id)+notesExt))
	if err != nil {
		return nil, err
	}
	tags := make([]db.Tag, len(file.Tags))
	for i, t := range file.Tags {
		tags[i] = db.Tag(t)
	}
	return &db.Note{
		Id:      id,
		Title:   file.Title,
		Body:    file.Content,
		Tags:    tags,
		Created: time.Unix(file.Created, 0),
		Updated: time.Unix(file.Updated, 0),
	}, nil
}

func (d *fileDb) PutNote(note *db.Note) error {
	return nil
}

func (d *fileDb) DeleteNote(id db.Id) error {
	return nil
}

func (d *fileDb) ListTags() ([]db.Tag, error) {
	return nil, nil
}

func (d *fileDb) TagSearch(tags []db.Tag) ([]*db.Note, error) {
	return nil, nil
}

func (d *fileDb) TextSearch(s string) ([]*db.Note, error) {
	return nil, nil
}

func (d *fileDb) Reconcile() (bool, []error) {
	return false, nil
}
