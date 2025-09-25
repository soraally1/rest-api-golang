package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"rest-api-golang/models"

	"github.com/gorilla/mux"
)

// In-memory storage for books
var books = make(map[string]*models.Book)

// In-memory token store (token -> username)
var activeTokens = make(map[string]string)

// In-memory users loaded from config
var usersFromConfig []models.User

// LoadUsers sets configured users at startup
func LoadUsers(users []models.User) {
	usersFromConfig = users
}

// AuthMiddleware protects endpoints with Bearer token except excluded paths
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow login and health without token
		if strings.HasPrefix(r.URL.Path, "/api/login") || strings.HasPrefix(r.URL.Path, "/health") {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Missing or invalid Authorization header",
			})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Invalid token",
			})
			return
		}
		if _, ok := activeTokens[token]; !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Token expired or invalid",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LoginRequest represents login payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// randomToken creates a random token string
func randomToken() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Login handles POST /api/login
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid JSON format",
		})
		return
	}

	// Simple auth check against config users
	valid := false
	for _, u := range usersFromConfig {
		if u.Username == req.Username && u.Password == req.Password {
			valid = true
			break
		}
	}
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid username or password",
		})
		return
	}

	token := randomToken()
	activeTokens[token] = req.Username

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token":   token,
	})
}

// Logout handles POST /api/logout (requires Bearer token)
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Missing or invalid Authorization header",
		})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid token",
		})
		return
	}

	if _, ok := activeTokens[token]; ok {
		delete(activeTokens, token)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}

// GetBooks handles GET /api/books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Convert map to slice
	bookList := make([]*models.Book, 0, len(books))
	for _, book := range books {
		bookList = append(bookList, book)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    bookList,
		"count":   len(bookList),
	})
}

// GetBook handles GET /api/books/{id}
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	book, exists := books[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Book not found",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    book,
	})
}

// CreateBook handles POST /api/books
func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid JSON format",
		})
		return
	}

	// Basic validation
	if req.Judul == "" || req.Author == "" || req.TahunTerbit < 1000 || req.TahunTerbit > 2024 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid book data. Judul and Author are required, TahunTerbit must be between 1000-2024",
		})
		return
	}

	book := models.NewBook(req)
	books[book.ID] = book

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Book created successfully",
		"data":    book,
	})
}

// UpdateBook handles PUT /api/books/{id}
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	book, exists := books[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Book not found",
		})
		return
	}

	var req models.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid JSON format",
		})
		return
	}

	// Update fields if provided
	if req.Judul != "" {
		book.Judul = req.Judul
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.TahunTerbit > 0 {
		if req.TahunTerbit < 1000 || req.TahunTerbit > 2024 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "TahunTerbit must be between 1000-2024",
			})
			return
		}
		book.TahunTerbit = req.TahunTerbit
	}

	book.UpdatedAt = time.Now()
	books[id] = book

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Book updated successfully",
		"data":    book,
	})
}

// DeleteBook handles DELETE /api/books/{id}
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	book, exists := books[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Book not found",
		})
		return
	}

	delete(books, id)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Book deleted successfully",
		"data":    book,
	})
}
