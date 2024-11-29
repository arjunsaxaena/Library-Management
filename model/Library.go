package model

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID           uuid.UUID `db:"id"`
	Title        string    `db:"title"`
	AuthorID     uuid.UUID `db:"author_id"`
	LocationID   uuid.UUID `db:"location_id"`
	IsCheckedOut bool      `db:"is_checked_out"`
	BookType     string    `db:"book_type"`
	CreatedAt    time.Time `db:"created_at"`
}

type Author struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type Location struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type User struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Class string    `db:"class"`
}

type IssuedBook struct {
	ID         uuid.UUID  `db:"id"`
	BookID     uuid.UUID  `db:"book_id"`
	UserID     uuid.UUID  `db:"user_id"`
	IssueDate  time.Time  `db:"issue_date"`
	ReturnDate *time.Time `db:"return_date"`
	LateFees   float64    `db:"late_fees"`
}

type Subject struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Language  string    `db:"language"`
	CreatedAt time.Time `db:"created_at"`
}

type Material struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Notes       string    `db:"notes"`
	Type        string    `db:"type"`
	Link        string    `db:"link"`
	Language    string    `db:"language"`
	SubjectName string    `db:"subject_name"`
	CreatedAt   time.Time `db:"created_at"`
}

// Interfaces for CRUD operations
type BookStore interface {
	Book(id uuid.UUID) (Book, error)
	Books() ([]Book, error)
	CreateBook(b *Book) error
	UpdateBook(b *Book) error
	DeleteBook(id uuid.UUID) error
}

type AuthorStore interface {
	Author(id uuid.UUID) (Author, error)
	Authors() ([]Author, error)
	CreateAuthor(a *Author) error
	UpdateAuthor(a *Author) error
	DeleteAuthor(id uuid.UUID) error
}

type LocationStore interface {
	Location(id uuid.UUID) (Location, error)
	Locations() ([]Location, error)
	CreateLocation(l *Location) error
	UpdateLocation(l *Location) error
	DeleteLocation(id uuid.UUID) error
}

type UserStore interface {
	User(id uuid.UUID) (User, error)
	Users() ([]User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(id uuid.UUID) error
}

type IssuedBookStore interface {
	CreateIssuedBook(issuedBook *IssuedBook) error
	ReturnBook(bookID uuid.UUID) (float64, error)
	GetIssuedBookByBookID(bookID uuid.UUID) (IssuedBook, error)
	IssuedBooks() ([]IssuedBook, error)
}

type SubjectStore interface {
	Subject(id uuid.UUID) (Subject, error)
	Subjects() ([]Subject, error)
	CreateSubject(s *Subject) error
	UpdateSubject(s *Subject) error
	DeleteSubject(id uuid.UUID) error
	SubjectByName(name string) (Subject, error)
}

type MaterialStore interface {
	Material(id uuid.UUID) (Material, error)
	Materials() ([]Material, error)
	CreateMaterial(material *Material) error
	UpdateMaterial(material *Material) error
	DeleteMaterial(id uuid.UUID) error
	GetMaterialsBySubject(subjectName string) ([]Material, error)
	GetMaterialsByLanguage(language string) ([]Material, error)
}
