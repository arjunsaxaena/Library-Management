package main

import (
	"log"
	"net/http"

	"github.com/arjunsaxaena/Library-Management/controllers"
	"github.com/arjunsaxaena/Library-Management/web"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "postgres://postgres:secret@postgres:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	bookStore := controllers.NewDBBookStore(db)
	authorStore := controllers.NewDBAuthorStore(db)
	locationStore := controllers.NewDBLocationStore(db)
	userStore := controllers.NewDBUserStore(db)
	issuedBookStore := controllers.NewDBIssuedBookStore(db)

	handler := web.NewHandler(bookStore, authorStore, locationStore, userStore, issuedBookStore)

	router := gin.Default()

	// Book routes
	router.GET("/books", handler.GetBooks)
	router.GET("/books/:id", handler.GetBook)
	router.POST("/books", handler.CreateBook)
	router.PUT("/books/:id", handler.UpdateBook)
	router.DELETE("/books/:id", handler.DeleteBook)

	// User routes
	router.GET("/users", handler.GetUsers)
	router.GET("/users/:id", handler.GetUser)
	router.POST("/users", handler.CreateUser)
	router.PUT("/users/:id", handler.UpdateUser)
	router.DELETE("/users/:id", handler.DeleteUser)

	// Location routes
	router.GET("/locations", handler.GetLocations)
	router.GET("/locations/:id", handler.GetLocation)
	router.POST("/locations", handler.CreateLocation)
	router.PUT("/locations/:id", handler.UpdateLocation)
	router.DELETE("/locations/:id", handler.DeleteLocation)

	// Author routes
	router.GET("/authors", handler.GetAuthors)
	router.GET("/authors/:id", handler.GetAuthor)
	router.POST("/authors", handler.CreateAuthor)
	router.PUT("/authors/:id", handler.UpdateAuthor)
	router.DELETE("/authors/:id", handler.DeleteAuthor)

	// Issued Book routes
	router.POST("books/issue", handler.IssueBook)
	router.POST("books/return", handler.ReturnBook)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Start the server
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
