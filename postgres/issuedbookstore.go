package postgres

import (
	"time"

	library "github.com/arjunsaxaena/Library-Management/library"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DBIssuedBookStore struct {
	db *sqlx.DB
}

func NewDBIssuedBookStore(db *sqlx.DB) *DBIssuedBookStore {
	return &DBIssuedBookStore{db: db}
}

func (s *DBIssuedBookStore) CreateIssuedBook(issuedBook *library.IssuedBook) error {
	// Start a transaction because two tables involved
	tx := s.db.MustBegin()
	defer tx.Rollback()

	_, err := tx.NamedExec(`INSERT INTO issued_books (id, book_id, user_id, issue_date) 
		VALUES (:id, :book_id, :user_id, :issue_date)`, issuedBook)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE books SET is_checked_out=true WHERE id=$1`, issuedBook.BookID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *DBIssuedBookStore) ReturnBook(bookID uuid.UUID) error {
	// Start a transaction
	tx := s.db.MustBegin()
	defer tx.Rollback()

	_, err := tx.Exec(`UPDATE issued_books SET return_date=$1 WHERE book_id=$2 AND return_date IS NULL`,
		time.Now(), bookID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE books SET is_checked_out=false WHERE id=$1`, bookID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *DBIssuedBookStore) GetIssuedBookByBookID(bookID uuid.UUID) (library.IssuedBook, error) {
	var issuedBook library.IssuedBook
	err := s.db.Get(&issuedBook, "SELECT * FROM issued_books WHERE book_id=$1 AND return_date IS NULL", bookID)
	return issuedBook, err
}

func (s *DBIssuedBookStore) IssuedBooks() ([]library.IssuedBook, error) {
	var issuedBooks []library.IssuedBook
	err := s.db.Select(&issuedBooks, "SELECT * FROM issued_books")
	return issuedBooks, err
}

func (s *DBIssuedBookStore) DeleteIssuedBook(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM issued_books WHERE id=$1", id)
	return err
}
