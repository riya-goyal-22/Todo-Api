package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"todo/logger"
	"todo/repository"
	"todo/service"
)

var client *sql.DB

func main() {
	logger.Info("Starting server...")
	initDB()
	todoHandler := &TodoHandler{
		service.NewTodoService(repository.NewTodoRepository(client)),
	}
	authHandler := &AuthHandler{
		service: service.NewUserService(repository.NewUserRepository(client)),
	}
	router := mux.NewRouter()

	// Route handlers
	router.HandleFunc("/auth", authHandler.ValidateUser).Methods("POST")
	router.HandleFunc("/todos", todoHandler.createTodoHandler).Methods("POST")
	router.HandleFunc("/todos", todoHandler.getTodosHandler).Methods("GET")
	router.HandleFunc("/todos/{id}", todoHandler.updateTodoHandler).Methods("PUT")
	router.HandleFunc("/todos/{id}", todoHandler.updateTodoStatus).Methods("PATCH")
	router.HandleFunc("/todos/{id}", todoHandler.deleteTodoHandler).Methods("DELETE")

	// Start the server
	logger.Info("Server running on port 8080")
	http.ListenAndServe(":8080", router)
}

func initDB() {
	var err error
	dsn := "root:mySql@tcp(localhost:3306)/banking"
	client, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Check if the connection is successful
	err = client.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create the todos table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		status VARCHAR(50) NOT NULL
	);`
	_, err = client.Exec(query)
	if err != nil {
		log.Fatalf("Failed to execute table creation query: %v", err)
	}
}
