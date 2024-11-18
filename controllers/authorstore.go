package controllers

import (
	"github.com/arjunsaxaena/Library-Management/model"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

type DBAuthorStore struct {
	db *sqlx.DB
}

func NewDBAuthorStore(db *sqlx.DB) *DBAuthorStore {
	return &DBAuthorStore{db: db}
}

func (s *DBAuthorStore) Author(id uuid.UUID) (model.Author, error) {
	var author model.Author
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("authors").Where(sb.Equal("id", id))

	query, args := sb.Build()
	err := s.db.Get(&author, query, args...)
	return author, err
}

func (s *DBAuthorStore) Authors() ([]model.Author, error) {
	var authors []model.Author
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("authors")

	query, args := sb.Build()
	err := s.db.Select(&authors, query, args...)
	return authors, err
}

func (s *DBAuthorStore) CreateAuthor(a *model.Author) error {
	sb := sqlbuilder.NewInsertBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.InsertInto("authors").Cols("id", "name").Values(a.ID, a.Name)

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBAuthorStore) UpdateAuthor(a *model.Author) error {
	sb := sqlbuilder.NewUpdateBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Update("authors").Set(
		sb.Assign("name", a.Name),
	).Where(sb.Equal("id", a.ID))

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBAuthorStore) DeleteAuthor(id uuid.UUID) error {
	sb := sqlbuilder.NewDeleteBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.DeleteFrom("authors").Where(sb.Equal("id", id))

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}
