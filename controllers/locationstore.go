package controllers

import (
	"github.com/arjunsaxaena/Library-Management/model"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

type DBLocationStore struct {
	db *sqlx.DB
}

func NewDBLocationStore(db *sqlx.DB) *DBLocationStore {
	return &DBLocationStore{db: db}
}

func (s *DBLocationStore) Location(id uuid.UUID) (model.Location, error) {
	var location model.Location
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("locations").Where(sb.Equal("id", id))

	query, args := sb.Build()
	err := s.db.Get(&location, query, args...)
	return location, err
}

func (s *DBLocationStore) Locations() ([]model.Location, error) {
	var locations []model.Location
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("locations")

	query, args := sb.Build()
	err := s.db.Select(&locations, query, args...)
	return locations, err
}

func (s *DBLocationStore) CreateLocation(l *model.Location) error {
	sb := sqlbuilder.NewInsertBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.InsertInto("locations").
		Cols("id", "name").
		Values(l.ID, l.Name)

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBLocationStore) UpdateLocation(l *model.Location) error {
	sb := sqlbuilder.NewUpdateBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Update("locations").
		Set(sb.Assign("name", l.Name)).
		Where(sb.Equal("id", l.ID))

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBLocationStore) DeleteLocation(id uuid.UUID) error {
	sb := sqlbuilder.NewDeleteBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.DeleteFrom("locations").Where(sb.Equal("id", id))

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}
