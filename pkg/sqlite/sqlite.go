package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/josephburnett/self/pkg/db"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteDb struct {
	db *sql.DB
}

func NewSqliteDb(path string) (db.Database, error) {
	d := &sqliteDb{}
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	d.db = database
	return d, nil
}

func (d *sqliteDb) ListNotes() ([]db.Id, error) {
	rows, err := d.db.Query("select id from notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := make([]db.Id, 0)
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, db.Id(id))
	}
	return ids, err
}

func (d *sqliteDb) GetNote(i db.Id) (*db.Note, error) {
	stmt, err := d.db.Prepare(`
select id, title, body, tags, created, updated
from notes
where id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var id, title, body, tags string
	var created, updated int64
	err = stmt.QueryRow(i).Scan(&id, &title, &body, &tags, &created, &updated)
	if err != nil {
		return nil, err
	}
	note := &db.Note{
		Id:      db.Id(id),
		Title:   title,
		Body:    body,
		Tags:    make([]db.Tag, 0),
		Created: time.Unix(created, 0),
		Updated: time.Unix(updated, 0),
	}
	ts := strings.Split(tags, ",")
	for _, t := range ts {
		note.Tags = append(note.Tags, db.Tag(t))
	}
	return note, nil
}

func (d *sqliteDb) PutNote(note *db.Note) error {
	stmt, err := d.db.Prepare(`
insert into notes(id, title, body, tags, created, updated) values(?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	tags := ""
	for _, t := range note.Tags {
		tags += string(t)
	}
	_, err = stmt.Exec(
		note.Id,
		note.Title,
		note.Body,
		tags,
		note.Created.Unix(),
		note.Updated.Unix(),
	)
	return err
}

func (d *sqliteDb) DeleteNote(id db.Id) error {
	return fmt.Errorf("unimplemented")
}

func (d *sqliteDb) ListTags() ([]db.Tag, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (d *sqliteDb) TagSearch(tags []db.Tag) ([]*db.Note, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (d *sqliteDb) TextSearch(s string) ([]*db.Note, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (d *sqliteDb) Reconcile() (bool, []error) {
	return false, nil
}

func (d *sqliteDb) Init() error {
	_, err := d.db.Exec(`
create table meta (key text not null primary key, value text not null);
create table notes (
  id text not null primary key,
  title text,
  body text,
  tags text,
  created integer,
  updated integer);
`)
	if err != nil {
		return err
	}
	_, err = d.db.Exec("insert into meta(key, value) values('version', '1')")
	return err
}
