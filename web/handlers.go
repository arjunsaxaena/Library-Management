package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arjunsaxaena/Library-Management/library"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	BookStore       library.BookStore
	AuthorStore     library.AuthorStore
	LocationStore   library.LocationStore
	UserStore       library.UserStore
	IssuedBookStore library.IssuedBookStore
}

func NewHandler(bs library.BookStore, as library.AuthorStore, ls library.LocationStore, us library.UserStore, ibs library.IssuedBookStore) *Handler {
	return &Handler{
		BookStore:       bs,
		AuthorStore:     as,
		LocationStore:   ls,
		UserStore:       us,
		IssuedBookStore: ibs,
	}
}

// Render form pages
func (h *Handler) CreateBookForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create_book.html", nil)
}

func (h *Handler) CreateUserForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create_user.html", nil)
}

func (h *Handler) CreateLocationForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create_location.html", nil)
}

func (h *Handler) IssueBookForm(c *gin.Context) {
	c.HTML(http.StatusOK, "issue_book.html", nil)
}

func (h *Handler) ReturnBookForm(c *gin.Context) {
	c.HTML(http.StatusOK, "return_book.html", nil)
}

func (h *Handler) IssueBook(c *gin.Context) {
	fmt.Println("Received issue book request")
	bookID, err := uuid.Parse(c.PostForm("book_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	userID, err := uuid.Parse(c.PostForm("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Create the issued book entry
	issuedBook := library.IssuedBook{
		ID:        uuid.New(),
		BookID:    bookID,
		UserID:    userID,
		IssueDate: time.Now(),
	}

	err = h.IssuedBookStore.CreateIssuedBook(&issuedBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to issue book: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book issued successfully", "issued_book": issuedBook})
}

func (h *Handler) ReturnBook(c *gin.Context) {
	// Parse the book ID from the form data
	bookID, err := uuid.Parse(c.PostForm("book_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Check if the book is currently issued
	issuedBook, err := h.IssuedBookStore.GetIssuedBookByBookID(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving issued book record"})
		return
	}
	if issuedBook.ReturnDate != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This book has already been returned"})
		return
	}

	// Update the return date and mark the book as available
	err = h.IssuedBookStore.ReturnBook(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to return book: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}

// Helper function to find or create an author
func (h *Handler) getOrCreateAuthor(name string) (*library.Author, error) {
	authors, err := h.AuthorStore.Authors()
	if err != nil {
		return nil, err
	}
	for _, a := range authors {
		if a.Name == name {
			return &a, nil
		}
	}
	newAuthor := library.Author{
		ID:   uuid.New(),
		Name: name,
	}
	if err := h.AuthorStore.CreateAuthor(&newAuthor); err != nil {
		return nil, err
	}
	return &newAuthor, nil
}

// Helper function to find or create a location
func (h *Handler) getOrCreateLocation(name string) (*library.Location, error) {
	locations, err := h.LocationStore.Locations()
	if err != nil {
		return nil, err
	}
	for _, l := range locations {
		if l.Name == name {
			return &l, nil
		}
	}
	newLocation := library.Location{
		ID:   uuid.New(),
		Name: name,
	}
	if err := h.LocationStore.CreateLocation(&newLocation); err != nil {
		return nil, err
	}
	return &newLocation, nil
}

func (h *Handler) CreateBook(c *gin.Context) {
	title := c.PostForm("title")
	authorName := c.PostForm("author_name")
	locationName := c.PostForm("location_name")

	author, err := h.getOrCreateAuthor(authorName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or find author"})
		return
	}

	location, err := h.getOrCreateLocation(locationName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or find location"})
		return
	}

	newBook := library.Book{
		ID:           uuid.New(),
		Title:        title,
		AuthorID:     author.ID,
		LocationID:   location.ID,
		IsCheckedOut: false,
	}

	if err := h.BookStore.CreateBook(&newBook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book created successfully", "book": newBook})
}

func (h *Handler) CreateUser(c *gin.Context) {
	name := c.PostForm("name")
	newUser := library.User{
		ID:   uuid.New(),
		Name: name,
	}
	if err := h.UserStore.CreateUser(&newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": newUser})
}

func (h *Handler) CreateLocation(c *gin.Context) {
	name := c.PostForm("name")
	newLocation := library.Location{
		ID:   uuid.New(),
		Name: name,
	}
	if err := h.LocationStore.CreateLocation(&newLocation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create location"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Location created successfully", "location": newLocation})
}

func (h *Handler) DeleteBook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := h.BookStore.DeleteBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.UserStore.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *Handler) DeleteLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	if err := h.LocationStore.DeleteLocation(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location deleted successfully"})
}

func (h *Handler) GetBooks(c *gin.Context) {
	books, err := h.BookStore.Books()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

func (h *Handler) GetBook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	book, err := h.BookStore.Book(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *Handler) GetLocations(c *gin.Context) {
	locations, err := h.LocationStore.Locations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve locations"})
		return
	}
	c.JSON(http.StatusOK, locations)
}

func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.UserStore.Users()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}
