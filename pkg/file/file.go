package file

import (
	"encoding/json"
	"io/ioutil"
)

type File struct {
	Content       string   `json:"content"`
	Created       int64    `json:"created"`
	Files         []string `json:"files"`
	Id            string   `json:"id"`
	IsFavorite    int32    `json:"isFavorite"`
	NotebookId    string   `json:"notebookId"`
	Tags          []string `json:"tags"`
	TaskAll       int32    `json:"taskAll"`
	TaskCompleted int32    `json:"taskCompleted"`
	Title         string   `json:"title"`
	Trash         int32    `json:"trash"`
	Type          string   `json:"type"`
	Updated       int64    `json:"updated"`
}

func Load(filename string) (*File, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	f := &File{}
	err = json.Unmarshal(b, f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func Store(filename string, file *File) error {
	b, err := json.Marshal(file)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}
