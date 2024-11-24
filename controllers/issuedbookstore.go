package controllers

import (
	"time"

	"github.com/arjunsaxaena/Library-Management/model"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

type DBIssuedBookStore struct {
	db *sqlx.DB
}

func NewDBIssuedBookStore(db *sqlx.DB) *DBIssuedBookStore {
	return &DBIssuedBookStore{db: db}
}

func (s *DBIssuedBookStore) CreateIssuedBook(issuedBook *model.IssuedBook) error {
	tx := s.db.MustBegin()
	defer tx.Rollback()

	sbInsert := sqlbuilder.NewInsertBuilder()
	sbInsert.SetFlavor(sqlbuilder.PostgreSQL)
	sbInsert.InsertInto("issued_books").
		Cols("id", "book_id", "user_id", "issue_date").
		Values(issuedBook.ID, issuedBook.BookID, issuedBook.UserID, issuedBook.IssueDate)

	queryInsert, argsInsert := sbInsert.Build()
	_, err := tx.Exec(queryInsert, argsInsert...)
	if err != nil {
		return err
	}

	sbUpdate := sqlbuilder.NewUpdateBuilder()
	sbUpdate.SetFlavor(sqlbuilder.PostgreSQL)
	sbUpdate.Update("books").
		Set(sbUpdate.Assign("is_checked_out", true)).
		Where(sbUpdate.Equal("id", issuedBook.BookID))

	queryUpdate, argsUpdate := sbUpdate.Build()
	_, err = tx.Exec(queryUpdate, argsUpdate...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *DBIssuedBookStore) ReturnBook(bookID uuid.UUID) (float64, error) {
	tx := s.db.MustBegin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var issuedBook model.IssuedBook
	sbSelect := sqlbuilder.NewSelectBuilder()
	sbSelect.SetFlavor(sqlbuilder.PostgreSQL)
	sbSelect.Select("*").
		From("issued_books").
		Where(sbSelect.Equal("book_id", bookID)).
		Where(sbSelect.IsNull("return_date"))

	querySelect, argsSelect := sbSelect.Build()
	err := s.db.Get(&issuedBook, querySelect, argsSelect...)
	if err != nil {
		return 0, err
	}

	var lateFees float64
	daysLate := time.Since(issuedBook.IssueDate).Hours() / 24
	if daysLate > 15 {
		lateFees = (daysLate - 15) * 2
	}

	currentTime := time.Now()

	sbUpdateReturn := sqlbuilder.NewUpdateBuilder()
	sbUpdateReturn.SetFlavor(sqlbuilder.PostgreSQL)
	sbUpdateReturn.Update("issued_books").Set(
		sbUpdateReturn.Assign("return_date", currentTime),
		sbUpdateReturn.Assign("late_fees", lateFees),
	).Where(sbUpdateReturn.Equal("book_id", bookID))

	queryUpdateReturn, argsUpdateReturn := sbUpdateReturn.Build()
	_, err = tx.Exec(queryUpdateReturn, argsUpdateReturn...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	sbUpdateBook := sqlbuilder.NewUpdateBuilder()
	sbUpdateBook.SetFlavor(sqlbuilder.PostgreSQL)
	sbUpdateBook.Update("books").
		Set(sbUpdateBook.Assign("is_checked_out", false)).
		Where(sbUpdateBook.Equal("id", bookID))

	queryUpdateBook, argsUpdateBook := sbUpdateBook.Build()
	_, err = tx.Exec(queryUpdateBook, argsUpdateBook...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return lateFees, nil
}

func (s *DBIssuedBookStore) GetIssuedBookByBookID(bookID uuid.UUID) (model.IssuedBook, error) {
	var issuedBook model.IssuedBook
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select(
		"id",
		"book_id",
		"user_id",
		"issue_date",
		"return_date",
		"CASE "+
			"WHEN return_date IS NULL THEN GREATEST(0, (EXTRACT(DAY FROM (CURRENT_DATE - issue_date)) - 15) * 2) "+
			"ELSE late_fees "+
			"END AS late_fees",
	).
		From("issued_books").
		Where(sb.Equal("book_id", bookID)).
		Where(sb.IsNull("return_date"))

	query, args := sb.Build()
	err := s.db.Get(&issuedBook, query, args...)
	return issuedBook, err
}

func (s *DBIssuedBookStore) IssuedBooks() ([]model.IssuedBook, error) {
	var issuedBooks []model.IssuedBook
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select(
		"id",
		"book_id",
		"user_id",
		"issue_date",
		"return_date",
		"CASE "+
			"WHEN return_date IS NULL THEN GREATEST(0, (EXTRACT(DAY FROM (CURRENT_DATE - issue_date)) - 15) * 2) "+
			"ELSE late_fees "+
			"END AS late_fees",
	).
		From("issued_books").
		Where(sb.IsNull("return_date"))

	query, args := sb.Build()
	err := s.db.Select(&issuedBooks, query, args...)
	return issuedBooks, err
}

func (s *DBIssuedBookStore) DeleteIssuedBook(id uuid.UUID) error {
	sb := sqlbuilder.NewDeleteBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.DeleteFrom("issued_books").
		Where(sb.Equal("id", id))

	query, args := sb.Build()
	_, err := s.db.Exec(query, args...)
	return err
}
