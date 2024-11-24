package controllers

import (
	"github.com/arjunsaxaena/Library-Management/model"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

type DBBookStore struct {
	db *sqlx.DB
}

func NewDBBookStore(db *sqlx.DB) *DBBookStore {
	return &DBBookStore{db: db}
}

func (s *DBBookStore) Book(id uuid.UUID) (model.Book, error) {
	var book model.Book
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("books").Where(sb.Equal("id", id))

	query, args := sb.Build()
	err := s.db.Get(&book, query, args...)
	return book, err
}

func (s *DBBookStore) Books() ([]model.Book, error) {
	var books []model.Book
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").
		From("books").
		Where(sb.Equal("is_checked_out", false))

	query, args := sb.Build()
	err := s.db.Select(&books, query, args...)
	return books, err
}

func (s *DBBookStore) CreateBook(b *model.Book) error {
	sb := sqlbuilder.NewInsertBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.InsertInto("books").
		Cols("id", "title", "author_id", "location_id", "is_checked_out", "book_type", "created_at").
		Values(b.ID, b.Title, b.AuthorID, b.LocationID, b.IsCheckedOut, b.BookType, b.CreatedAt)

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBBookStore) UpdateBook(b *model.Book) error {
	sb := sqlbuilder.NewUpdateBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Update("books").
		Set(
			sb.Assign("title", b.Title),
			sb.Assign("author_id", b.AuthorID),
			sb.Assign("location_id", b.LocationID),
			sb.Assign("is_checked_out", b.IsCheckedOut),
			sb.Assign("book_type", b.BookType),
		).
		Where(sb.Equal("id", b.ID))

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBBookStore) DeleteBook(id uuid.UUID) error {
	sb := sqlbuilder.NewDeleteBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.DeleteFrom("books").Where(sb.Equal("id", id))

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}
