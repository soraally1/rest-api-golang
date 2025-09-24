package main

import (
	"fmt"
	"log"
	"net/http"

	"rest-api-golang/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Create router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

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

	// Apply CORS middleware
	handler := corsHandler(r)

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
