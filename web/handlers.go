package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/arjunsaxaena/Library-Management/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	BookStore       model.BookStore
	AuthorStore     model.AuthorStore
	LocationStore   model.LocationStore
	UserStore       model.UserStore
	IssuedBookStore model.IssuedBookStore
}

func NewHandler(bs model.BookStore, as model.AuthorStore, ls model.LocationStore, us model.UserStore, ibs model.IssuedBookStore) *Handler {
	return &Handler{
		BookStore:       bs,
		AuthorStore:     as,
		LocationStore:   ls,
		UserStore:       us,
		IssuedBookStore: ibs,
	}
}

// ISSUE AND RETURN HANDLERS

func (h *Handler) IssueBook(c *gin.Context) {
	fmt.Println("Received issue book request")

	type IssueBookRequest struct {
		BookID uuid.UUID `json:"book_id" binding:"required"`
		UserID uuid.UUID `json:"user_id" binding:"required"`
	}

	var request IssueBookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	existingIssuedBook, err := h.IssuedBookStore.GetIssuedBookByBookID(request.BookID)
	if err == nil && existingIssuedBook.UserID == request.UserID && existingIssuedBook.ReturnDate == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "This book is already issued to the user and cannot be reissued."})
		return
	}

	issuedBook := model.IssuedBook{
		ID:        uuid.New(),
		BookID:    request.BookID,
		UserID:    request.UserID,
		IssueDate: time.Now(),
	}

	err = h.IssuedBookStore.CreateIssuedBook(&issuedBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to issue book: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Book issued successfully",
		"issued_book": issuedBook,
	})
}

func (h *Handler) ReturnBook(c *gin.Context) {
	fmt.Println("Received return book request")

	type ReturnBookRequest struct {
		BookID uuid.UUID `json:"book_id" binding:"required"`
		UserID uuid.UUID `json:"user_id" binding:"required"`
	}

	var request ReturnBookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	issuedBook, err := h.IssuedBookStore.GetIssuedBookByBookID(request.BookID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "This book is not currently issued."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving issued book record"})
		return
	}

	if issuedBook.UserID != request.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "This book was not issued to the user and cannot be returned."})
		return
	}

	if issuedBook.ReturnDate != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This book has already been returned and cannot be returned again."})
		return
	}

	err = h.IssuedBookStore.ReturnBook(request.BookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to return book: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}

// HELPER FUNCTIONS

func (h *Handler) getOrCreateAuthor(name string) (*model.Author, error) {
	authors, err := h.AuthorStore.Authors()
	if err != nil {
		return nil, err
	}
	for _, a := range authors {
		if a.Name == name {
			return &a, nil
		}
	}
	newAuthor := model.Author{
		ID:   uuid.New(),
		Name: name,
	}
	if err := h.AuthorStore.CreateAuthor(&newAuthor); err != nil {
		return nil, err
	}
	return &newAuthor, nil
}

func (h *Handler) getOrCreateLocation(name string) (*model.Location, error) {
	locations, err := h.LocationStore.Locations()
	if err != nil {
		return nil, err
	}
	for _, l := range locations {
		if l.Name == name {
			return &l, nil
		}
	}
	newLocation := model.Location{
		ID:   uuid.New(),
		Name: name,
	}
	if err := h.LocationStore.CreateLocation(&newLocation); err != nil {
		return nil, err
	}
	return &newLocation, nil
}

// GET HANDLERS (BY ID OR ALL)

func (h *Handler) GetBooks(c *gin.Context) {
	books, err := h.BookStore.Books()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *Handler) GetBook(c *gin.Context) {
	idParam := c.Param("id")
	bookID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	book, err := h.BookStore.Book(bookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *Handler) GetAuthors(c *gin.Context) {
	authors, err := h.AuthorStore.Authors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve authors"})
		return
	}

	c.JSON(http.StatusOK, authors)
}

func (h *Handler) GetAuthor(c *gin.Context) {
	idParam := c.Param("id")
	authorID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	author, err := h.AuthorStore.Author(authorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	c.JSON(http.StatusOK, author)
}

func (h *Handler) GetLocations(c *gin.Context) {
	locations, err := h.LocationStore.Locations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve locations"})
		return
	}

	c.JSON(http.StatusOK, locations)
}

func (h *Handler) GetLocation(c *gin.Context) {
	fmt.Println("Received get location request")

	type GetLocationResponse struct {
		Location model.Location `json:"location"`
	}

	idParam := c.Param("id")
	locationID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	location, err := h.LocationStore.Location(locationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}

	c.JSON(http.StatusOK, GetLocationResponse{Location: location})
}

func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.UserStore.Users()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUser(c *gin.Context) {
	fmt.Println("Received get user request")

	type GetUserResponse struct {
		User model.User `json:"user"`
	}

	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.UserStore.User(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, GetUserResponse{User: user})
}

// CREATE HANDLERS

func (h *Handler) CreateBook(c *gin.Context) {
	type CreateBookRequest struct {
		Title        string `json:"title"`
		AuthorName   string `json:"author_name"`
		LocationName string `json:"location_name"`
	}

	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	existingBooks, err := h.BookStore.Books()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing books"})
		return
	}

	for _, book := range existingBooks {
		if book.Title == req.Title {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "A book with the same title already exists",
				"book_id": book.ID,
			})
			return
		}
	}

	newUUID := uuid.New()

	for _, book := range existingBooks {
		if book.ID == newUUID {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "A book with the same UUID already exists",
				"book_id": book.ID,
			})
			return
		}
	}

	author, err := h.getOrCreateAuthor(req.AuthorName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or find author"})
		return
	}

	location, err := h.getOrCreateLocation(req.LocationName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or find location"})
		return
	}

	newBook := model.Book{
		ID:           newUUID,
		Title:        req.Title,
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
	fmt.Println("Received create user request")

	type CreateUserRequest struct {
		Name string `json:"name"`
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	//existingUsers, err := h.UserStore.Users()
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing users"})
	//	return
	//}

	//for _, user := range existingUsers {
	//	if user.Name == req.Name {
	//		c.JSON(http.StatusConflict, gin.H{
	//			"error":   "A user with the same name already exists",
	//			"user_id": user.ID,
	//		})
	//		return
	//	}
	//}

	newUUID := uuid.New()

	newUser := model.User{
		ID:   newUUID,
		Name: req.Name,
	}

	if err := h.UserStore.CreateUser(&newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": newUser})
}

func (h *Handler) CreateLocation(c *gin.Context) {
	fmt.Println("Received create location request")

	type CreateLocationRequest struct {
		Name string `json:"name"`
	}

	var req CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	existingLocations, err := h.LocationStore.Locations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing locations"})
		return
	}

	for _, location := range existingLocations {
		if location.Name == req.Name {
			c.JSON(http.StatusConflict, gin.H{
				"error":       "A location with the same name already exists",
				"location_id": location.ID,
			})
			return
		}
	}

	newLocation := model.Location{
		ID:   uuid.New(),
		Name: req.Name,
	}

	if err := h.LocationStore.CreateLocation(&newLocation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location created successfully", "location": newLocation})
}

func (h *Handler) CreateAuthor(c *gin.Context) {
	fmt.Println("Received create author request")

	type CreateAuthorRequest struct {
		Name string `json:"name"`
	}

	var req CreateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	//existingAuthors, err := h.AuthorStore.Authors()
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing authors"})
	//	return
	//}

	//for _, author := range existingAuthors {
	//	if author.Name == req.Name {
	//		c.JSON(http.StatusConflict, gin.H{
	//			"error":     "An author with the same name already exists",
	//			"author_id": author.ID,
	//		})
	//		return
	//	}
	//}

	newAuthor := model.Author{
		ID:   uuid.New(),
		Name: req.Name,
	}

	if err := h.AuthorStore.CreateAuthor(&newAuthor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Author created successfully", "author": newAuthor})
}

// UPDATE HANDLERS

func (h *Handler) UpdateBook(c *gin.Context) {
	idParam := c.Param("id")
	bookID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	book.ID = bookID

	if err := h.BookStore.UpdateBook(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully", "book": book})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user.ID = userID

	if err := h.UserStore.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

func (h *Handler) UpdateLocation(c *gin.Context) {
	idParam := c.Param("id")
	locationID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	var location model.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	location.ID = locationID

	if err := h.LocationStore.UpdateLocation(&location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location updated successfully", "location": location})
}

func (h *Handler) UpdateAuthor(c *gin.Context) {
	idParam := c.Param("id")
	authorID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	var author model.Author
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	author.ID = authorID

	if err := h.AuthorStore.UpdateAuthor(&author); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Author updated successfully", "author": author})
}

// DELETE HANDLERS

func (h *Handler) DeleteBook(c *gin.Context) {
	idParam := c.Param("id")
	bookID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := h.BookStore.DeleteBook(bookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.UserStore.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *Handler) DeleteLocation(c *gin.Context) {
	idParam := c.Param("id")
	locationID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	if err := h.LocationStore.DeleteLocation(locationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location deleted successfully"})
}

func (h *Handler) DeleteAuthor(c *gin.Context) {
	idParam := c.Param("id")
	authorID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	if err := h.AuthorStore.DeleteAuthor(authorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete author"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Author deleted successfully"})
}
