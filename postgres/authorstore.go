package postgres

import (
	library "github.com/arjunsaxaena/Library-Management/library"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DBAuthorStore struct {
	db *sqlx.DB
}

func NewDBAuthorStore(db *sqlx.DB) *DBAuthorStore {
	return &DBAuthorStore{db: db}
}

func (s *DBAuthorStore) Author(id uuid.UUID) (library.Author, error) {
	var author library.Author
	err := s.db.Get(&author, "SELECT * FROM authors WHERE id=$1", id)
	return author, err
}

func (s *DBAuthorStore) Authors() ([]library.Author, error) {
	var authors []library.Author
	err := s.db.Select(&authors, "SELECT * FROM authors")
	return authors, err
}

func (s *DBAuthorStore) CreateAuthor(a *library.Author) error {
	_, err := s.db.NamedExec(`INSERT INTO authors (id, name) VALUES (:id, :name)`, a)
	return err
}

func (s *DBAuthorStore) UpdateAuthor(a *library.Author) error {
	_, err := s.db.NamedExec(`UPDATE authors SET name=:name WHERE id=:id`, a)
	return err
}

func (s *DBAuthorStore) DeleteAuthor(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM authors WHERE id=$1", id)
	return err
}
