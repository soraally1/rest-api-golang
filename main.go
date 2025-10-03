package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"rest-api-golang/database"
	"rest-api-golang/handlers"
	"rest-api-golang/models"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

// @title Book Management API
// @version 1.0
// @description A REST API for managing books with authentication
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Connect to database
	if err := database.ConnectDatabase(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDatabase()

	// Create tables and seed data
	if err := database.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
	if err := database.SeedData(); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	// Initialize repositories
	handlers.InitializeRepositories()

	// Create router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/login", handlers.Login).Methods("POST")
	api.HandleFunc("/logout", handlers.Logout).Methods("POST")

	// Book routes
	api.HandleFunc("/books", handlers.GetBooks).Methods("GET")
	api.HandleFunc("/books", handlers.CreateBook).Methods("POST")
	api.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")
	api.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
	api.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "healthy", "message": "Server is running"}`)
	}).Methods("GET")

	// API Documentation endpoint
	r.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>Book Management API Documentation</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { background: #f5f5f5; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .method { font-weight: bold; color: #007bff; }
        .auth { color: #dc3545; font-weight: bold; }
        .no-auth { color: #28a745; font-weight: bold; }
    </style>
</head>
<body>
    <h1>📚 Book Management API</h1>
    <p>REST API for managing books with authentication</p>
    
    <h2>🔐 Authentication</h2>
    <p>Most endpoints require a Bearer token. Login first to get a token.</p>
    
    <h2>📋 Endpoints</h2>
    
    <div class="endpoint">
        <span class="method">POST</span> /api/login <span class="no-auth">(No Auth)</span><br>
        <strong>Description:</strong> Login to get authentication token<br>
        <strong>Body:</strong> {"username": "admin", "password": "admin123"}<br>
        <strong>Response:</strong> {"success": true, "token": "..."}
    </div>
    
    <div class="endpoint">
        <span class="method">POST</span> /api/logout <span class="auth">(Auth Required)</span><br>
        <strong>Description:</strong> Logout and revoke token<br>
        <strong>Headers:</strong> Authorization: Bearer YOUR_TOKEN
    </div>
    
    <div class="endpoint">
        <span class="method">GET</span> /api/books <span class="auth">(Auth Required)</span><br>
        <strong>Description:</strong> Get all books<br>
        <strong>Headers:</strong> Authorization: Bearer YOUR_TOKEN
    </div>
    
    <div class="endpoint">
        <span class="method">POST</span> /api/books <span class="auth">(Auth Required)</span><br>
        <strong>Description:</strong> Create a new book<br>
        <strong>Headers:</strong> Authorization: Bearer YOUR_TOKEN<br>
        <strong>Body:</strong> {"judul": "Book Title", "author": "Author Name", "tahun_terbit": 2024}
    </div>
    
    <div class="endpoint">
        <span class="method">GET</span> /api/books/{id} <span class="auth">(Auth Required)</span><br>
        <strong>Description:</strong> Get book by ID<br>
        <strong>Headers:</strong> Authorization: Bearer YOUR_TOKEN
    </div>
    
    <div class="endpoint">
        <span class="method">PUT</span> /api/books/{id} <span class="auth">(Auth Required)</span><br>
        <strong>Description:</strong> Update book by ID<br>
        <strong>Headers:</strong> Authorization: Bearer YOUR_TOKEN<br>
        <strong>Body:</strong> {"judul": "New Title", "author": "New Author", "tahun_terbit": 2024}
    </div>
    
    <div class="endpoint">
        <span class="method">DELETE</span> /api/books/{id} <span class="auth">(Auth Required)</span><br>
        <strong>Description:</strong> Delete book by ID<br>
        <strong>Headers:</strong> Authorization: Bearer YOUR_TOKEN
    </div>
    
    <div class="endpoint">
        <span class="method">GET</span> /health <span class="no-auth">(No Auth)</span><br>
        <strong>Description:</strong> Health check endpoint
    </div>
    
    <h2>👥 Default Users</h2>
    <ul>
        <li><strong>Username:</strong> admin, <strong>Password:</strong> admin123</li>
        <li><strong>Username:</strong> user, <strong>Password:</strong> user123</li>
    </ul>
    
    <h2>🧪 Testing</h2>
    <p>Use tools like Postman, Insomnia, or curl to test the API endpoints.</p>
</body>
</html>`
		fmt.Fprint(w, html)
	}).Methods("GET")

	// CORS middleware
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Apply CORS and Auth middleware
	secured := handlers.AuthMiddleware(r)
	handler := corsHandler(secured)

	// Start server
	port := ":8080"
	fmt.Printf("🚀 Server starting on port %s\n", port)
	fmt.Println("📚 Book API Endpoints:")
	fmt.Println("  POST   /api/login       - Login to get token")
	fmt.Println("  POST   /api/logout      - Logout (requires token)")
	fmt.Println("  GET    /api/books       - Get all books (requires token)")
	fmt.Println("  POST   /api/books       - Create a new book (requires token)")
	fmt.Println("  GET    /api/books/{id}  - Get book by ID (requires token)")
	fmt.Println("  PUT    /api/books/{id}  - Update book by ID (requires token)")
	fmt.Println("  DELETE /api/books/{id}  - Delete book by ID (requires token)")
	fmt.Println("  GET    /health          - Health check")
	fmt.Println("  GET    /docs            - API documentation")
	fmt.Println()

	log.Fatal(http.ListenAndServe(port, handler))
}

// loadUsersFromConfig reads users from config.yaml; returns empty if not found
func loadUsersFromConfig() []models.User {
	path := "config.yaml"
	f, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("⚠️  Could not read %s: %v\n", path, err)
		return nil
	}
	var raw struct {
		Users []models.User `yaml:"users"`
	}
	if err := yaml.Unmarshal(f, &raw); err != nil {
		fmt.Printf("⚠️  Could not parse %s: %v\n", path, err)
		return nil
	}
	return raw.Users
}
