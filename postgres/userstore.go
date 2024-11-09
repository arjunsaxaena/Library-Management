package postgres

import (
	library "github.com/arjunsaxaena/Library-Management/library"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DBUserStore struct {
	db *sqlx.DB
}

func NewDBUserStore(db *sqlx.DB) *DBUserStore {
	return &DBUserStore{db: db}
}

func (s *DBUserStore) User(id uuid.UUID) (library.User, error) {
	var user library.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}

func (s *DBUserStore) Users() ([]library.User, error) {
	var users []library.User
	err := s.db.Select(&users, "SELECT * FROM users")
	return users, err
}

func (s *DBUserStore) CreateUser(u *library.User) error {
	_, err := s.db.NamedExec(`INSERT INTO users (id, name) VALUES (:id, :name)`, u)
	return err
}

func (s *DBUserStore) UpdateUser(u *library.User) error {
	_, err := s.db.NamedExec(`UPDATE users SET name=:name WHERE id=:id`, u)
	return err
}

func (s *DBUserStore) DeleteUser(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}
