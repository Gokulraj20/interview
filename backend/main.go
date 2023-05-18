package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
    "os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
  host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Connect to the PostgreSQL database
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the users table if it doesn't exist
	createTable(db)

	// Create a Gin router
	router := gin.Default()

	// Define the API routes
	router.POST("/users", createUser(db))
	router.GET("/users/:id", getUser(db))
	router.DELETE("/users/:id", deleteUser(db))

	// Start the server
	router.Run(":8080")
}

func createTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50),
		email VARCHAR(50)
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func createUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id"
		var userID int64
		err := db.QueryRow(query, user.Name, user.Email).Scan(&userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": userID})
	}
}

func getUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		query := "SELECT name, email FROM users WHERE id = $1"
		var name, email string
		err := db.QueryRow(query, id).Scan(&name, &email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"name": name, "email": email})
	}
}

func deleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		query := "DELETE FROM users WHERE id = $1"
		result, err := db.Exec(query, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
