package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"rest-api-golang/handlers"
	"rest-api-golang/models"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

func main() {
	// Create router
	r := mux.NewRouter()

	// Load users from config.yaml
	users := loadUsersFromConfig()
	handlers.LoadUsers(users)

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
	fmt.Println("  GET    /api/books       - Get all books")
	fmt.Println("  POST   /api/books       - Create a new book")
	fmt.Println("  GET    /api/books/{id}  - Get book by ID")
	fmt.Println("  PUT    /api/books/{id}  - Update book by ID")
	fmt.Println("  DELETE /api/books/{id}  - Delete book by ID")
	fmt.Println("  GET    /health          - Health check")
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
