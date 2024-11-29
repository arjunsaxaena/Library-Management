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
	db, err := sqlx.Connect("postgres", "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable") // go run cmd/main.go
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	bookStore := controllers.NewDBBookStore(db)
	authorStore := controllers.NewDBAuthorStore(db)
	locationStore := controllers.NewDBLocationStore(db)
	userStore := controllers.NewDBUserStore(db)
	issuedBookStore := controllers.NewDBIssuedBookStore(db)
	subjectStore := controllers.NewDBSubjectStore(db)
	materialStore := controllers.NewDBMaterialStore(db)

	handler := web.NewHandler(bookStore, authorStore, locationStore, userStore, issuedBookStore, subjectStore, materialStore)

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
	router.GET("books/issue/:id", handler.GetIssuedBook)
	router.GET("books/issue", handler.GetIssuedBooks)
	router.POST("books/issue", handler.IssueBook)
	router.POST("books/return", handler.ReturnBook)

	// Material routes
	router.GET("/materials", handler.GetMaterials)
	router.GET("/materials/:id", handler.GetMaterial)
	router.POST("/materials", handler.CreateMaterial)
	router.GET("/materials/subject/:subject_name", handler.GetMaterialsBySubject)
	router.GET("/materials/language/:language", handler.GetMaterialsByLanguage)
	router.PUT("/materials/:id", handler.UpdateMaterial)
	router.DELETE("/materials/:id", handler.DeleteMaterial)

	// Subject routes
	router.GET("/subjects", handler.GetSubjects)
	router.GET("/subjects/:id", handler.GetSubject)
	router.POST("/subjects", handler.CreateSubject)
	router.GET("/subjects/name/:name", handler.GetSubjectByName)
	router.PUT("/subjects/:id", handler.UpdateSubject)
	router.DELETE("/subjects/:id", handler.DeleteSubject)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Start the server
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
