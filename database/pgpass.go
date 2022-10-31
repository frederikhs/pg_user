package database

import (
	"github.com/tg/pgpass"
)

func ListHosts() ([]pgpass.Entry, error) {
	f, err := pgpass.OpenDefault()
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var entries []pgpass.Entry
	er := pgpass.NewEntryReader(f)
	for er.Next() {
		e := er.Entry()
		entries = append(entries, e)
	}

	return entries, er.Err()
}
