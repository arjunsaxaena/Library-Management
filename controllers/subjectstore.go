package controllers

import (
	"database/sql"
	"errors"

	"github.com/arjunsaxaena/Library-Management/model"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

type DBSubjectStore struct {
	db *sqlx.DB
}

func NewDBSubjectStore(db *sqlx.DB) *DBSubjectStore {
	return &DBSubjectStore{db: db}
}

func (s *DBSubjectStore) Subject(id uuid.UUID) (model.Subject, error) {
	var subject model.Subject
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("subjects").Where(sb.Equal("id", id))

	query, args := sb.Build()
	err := s.db.Get(&subject, query, args...)
	return subject, err
}

func (s *DBSubjectStore) Subjects() ([]model.Subject, error) {
	var subjects []model.Subject
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("subjects")

	query, args := sb.Build()
	err := s.db.Select(&subjects, query, args...)
	return subjects, err
}

func (s *DBSubjectStore) CreateSubject(subject *model.Subject) error {
	sb := sqlbuilder.NewInsertBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.InsertInto("subjects").Cols("id", "name", "language", "created_at").
		Values(subject.ID, subject.Name, subject.Language, subject.CreatedAt)

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBSubjectStore) UpdateSubject(subject *model.Subject) error {
	sb := sqlbuilder.NewUpdateBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Update("subjects").Set(
		sb.Assign("name", subject.Name),
		sb.Assign("language", subject.Language),
	).Where(sb.Equal("id", subject.ID))

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBSubjectStore) DeleteSubject(id uuid.UUID) error {
	sb := sqlbuilder.NewDeleteBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.DeleteFrom("subjects").Where(sb.Equal("id", id))

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBSubjectStore) SubjectByName(name string) (model.Subject, error) {
	var subject model.Subject
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("subjects").Where(sb.Equal("name", name))

	query, args := sb.Build()
	err := s.db.Get(&subject, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Subject{ID: uuid.Nil}, nil
		}
		return model.Subject{}, err
	}
	return subject, nil
}
