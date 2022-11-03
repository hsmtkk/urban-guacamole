package entryrepo

import "github.com/hsmtkk/urban-guacamole/entry"

type EntryRepo interface {
	Scan() ([]entry.Entry, error)
	Save(e entry.Entry) error
	Get(id int64) (entry.Entry, error)
	Delete(id int64) error
}
