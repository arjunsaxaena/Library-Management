package main

import (
	"log"
	"net/http"

	"github.com/arjunsaxaena/Library-Management/postgres"
	"github.com/arjunsaxaena/Library-Management/web"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Connect to the PostgreSQL database
	db, err := sqlx.Connect("postgres", "postgres://postgres:secret@postgres:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	// Initialize the data stores
	bookStore := postgres.NewDBBookStore(db)
	authorStore := postgres.NewDBAuthorStore(db)
	locationStore := postgres.NewDBLocationStore(db)
	userStore := postgres.NewDBUserStore(db)
	issuedBookStore := postgres.NewDBIssuedBookStore(db)

	// Initialize the handler with all stores
	handler := web.NewHandler(bookStore, authorStore, locationStore, userStore, issuedBookStore)

	// Set up the router
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// Routes for the main index and displaying available resources
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Routes for books
	router.GET("/books", handler.GetBooks)
	router.GET("/books/:id", handler.GetBook)
	router.GET("/create_book", handler.CreateBookForm)
	router.POST("/create_book", handler.CreateBook)
	router.DELETE("/books/:id", handler.DeleteBook)

	// Routes for users
	router.GET("/users", handler.GetUsers)
	router.GET("/create_user", handler.CreateUserForm)
	router.POST("/create_user", handler.CreateUser)
	router.DELETE("/users/:id", handler.DeleteUser)

	// Routes for locations
	router.GET("/locations", handler.GetLocations)
	router.GET("/create_location", handler.CreateLocationForm)
	router.POST("/create_location", handler.CreateLocation)
	router.DELETE("/locations/:id", handler.DeleteLocation)

	// Routes for issuing and returning books
	router.GET("/issue-book-form", handler.IssueBookForm)
	router.POST("/issue-book", handler.IssueBook)
	router.GET("/return-book-form", handler.ReturnBookForm)
	router.POST("/return-book", handler.ReturnBook)

	// Start the server
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
