package controllers

import (
	"github.com/arjunsaxaena/Library-Management/model"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

type DBMaterialStore struct {
	db *sqlx.DB
}

func NewDBMaterialStore(db *sqlx.DB) *DBMaterialStore {
	return &DBMaterialStore{db: db}
}

func (s *DBMaterialStore) Material(id uuid.UUID) (model.Material, error) {
	var material model.Material
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("materials").Where(sb.Equal("id", id))

	query, args := sb.Build()
	err := s.db.Get(&material, query, args...)
	return material, err
}

func (s *DBMaterialStore) Materials() ([]model.Material, error) {
	var materials []model.Material
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("materials")

	query, args := sb.Build()
	err := s.db.Select(&materials, query, args...)
	return materials, err
}

func (s *DBMaterialStore) CreateMaterial(material *model.Material) error {
	sb := sqlbuilder.NewInsertBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.InsertInto("materials").Cols(
		"id", "title", "description", "notes", "type", "link", "language", "subject_name", "created_at",
	).Values(
		material.ID, material.Title, material.Description, material.Notes, material.Type, material.Link, material.Language, material.SubjectName, material.CreatedAt,
	)

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBMaterialStore) UpdateMaterial(material *model.Material) error {
	sb := sqlbuilder.NewUpdateBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Update("materials").Set(
		sb.Assign("title", material.Title),
		sb.Assign("description", material.Description),
		sb.Assign("notes", material.Notes),
		sb.Assign("type", material.Type),
		sb.Assign("link", material.Link),
		sb.Assign("language", material.Language),
		sb.Assign("subject_name", material.SubjectName),
	).Where(sb.Equal("id", material.ID))

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBMaterialStore) DeleteMaterial(id uuid.UUID) error {
	sb := sqlbuilder.NewDeleteBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.DeleteFrom("materials").Where(sb.Equal("id", id))

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *DBMaterialStore) GetMaterialsBySubject(subjectName string) ([]model.Material, error) {
	var materials []model.Material
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("materials").Where(sb.Equal("subject_name", subjectName))

	query, args := sb.Build()
	err := s.db.Select(&materials, query, args...)
	return materials, err
}

func (s *DBMaterialStore) GetMaterialsByLanguage(language string) ([]model.Material, error) {
	var materials []model.Material
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("*").From("materials").Where(sb.Equal("language", language))

	query, args := sb.Build()
	err := s.db.Select(&materials, query, args...)
	return materials, err
}
