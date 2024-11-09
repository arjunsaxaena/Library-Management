package postgres

import (
	library "github.com/arjunsaxaena/Library-Management/library"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DBLocationStore struct {
	db *sqlx.DB
}

func NewDBLocationStore(db *sqlx.DB) *DBLocationStore {
	return &DBLocationStore{db: db}
}

func (s *DBLocationStore) Location(id uuid.UUID) (library.Location, error) {
	var location library.Location
	err := s.db.Get(&location, "SELECT * FROM locations WHERE id=$1", id)
	return location, err
}

func (s *DBLocationStore) Locations() ([]library.Location, error) {
	var locations []library.Location
	err := s.db.Select(&locations, "SELECT * FROM locations")
	return locations, err
}

func (s *DBLocationStore) CreateLocation(l *library.Location) error {
	_, err := s.db.NamedExec(`INSERT INTO locations (id, name) VALUES (:id, :name)`, l)
	return err
}

func (s *DBLocationStore) UpdateLocation(l *library.Location) error {
	_, err := s.db.NamedExec(`UPDATE locations SET name=:name WHERE id=:id`, l)
	return err
}

func (s *DBLocationStore) DeleteLocation(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM locations WHERE id=$1", id)
	return err
}
