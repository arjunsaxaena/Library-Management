package controllers

import (
	"github.com/arjunsaxaena/Library-Management/model"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

type DBUserStore struct {
	db *sqlx.DB
}

func NewDBUserStore(db *sqlx.DB) *DBUserStore {
	return &DBUserStore{db: db}
}

func (s *DBUserStore) User(id uuid.UUID) (model.User, error) {
	var user model.User
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("users").Where(sb.Equal("id", id))

	query, args := sb.Build()
	err := s.db.Get(&user, query, args...)
	return user, err
}

func (s *DBUserStore) Users() ([]model.User, error) {
	var users []model.User
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("users")

	query, args := sb.Build()
	err := s.db.Select(&users, query, args...)
	return users, err
}

func (s *DBUserStore) CreateUser(u *model.User) error {
	sb := sqlbuilder.NewInsertBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.InsertInto("users").
		Cols("id", "name", "standard").
		Values(u.ID, u.Name, u.Standard)

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBUserStore) UpdateUser(u *model.User) error {
	sb := sqlbuilder.NewUpdateBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Update("users").
		Set(
			sb.Assign("name", u.Name),
			sb.Assign("standard", u.Standard),
		).
		Where(sb.Equal("id", u.ID))

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBUserStore) DeleteUser(id uuid.UUID) error {
	sb := sqlbuilder.NewDeleteBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.DeleteFrom("users").Where(sb.Equal("id", id))

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}
