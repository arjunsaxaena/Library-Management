package postgres

import (
	library "github.com/arjunsaxaena/Library-Management/library"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DBBookStore struct {
	db *sqlx.DB
}

func NewDBBookStore(db *sqlx.DB) *DBBookStore {
	return &DBBookStore{db: db}
}

func (s *DBBookStore) Book(id uuid.UUID) (library.Book, error) {
	var book library.Book
	err := s.db.Get(&book, "SELECT * FROM books WHERE id=$1", id)
	return book, err
}

func (s *DBBookStore) Books() ([]library.Book, error) {
	var books []library.Book
	err := s.db.Select(&books, "SELECT * FROM books")
	return books, err
}

func (s *DBBookStore) CreateBook(b *library.Book) error {
	_, err := s.db.NamedExec(`INSERT INTO books (id, title, author_id, location_id, is_checked_out)   
		VALUES (:id, :title, :author_id, :location_id, :is_checked_out)`, b) // '_' when we want to ignore the result of the query
	return err
}

func (s *DBBookStore) UpdateBook(b *library.Book) error {
	_, err := s.db.NamedExec(`UPDATE books SET title=:title, author_id=:author_id, 
		location_id=:location_id, is_checked_out=:is_checked_out WHERE id=:id`, b)
	return err
}

func (s *DBBookStore) DeleteBook(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM books WHERE id=$1", id)
	return err
}
