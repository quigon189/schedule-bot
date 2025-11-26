package repository

import "database/sql"

type DB struct {
	*sql.DB
}

func NewPostgresDB(connString string) (*DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (d *DB) Close() error {
	return d.DB.Close()
}
