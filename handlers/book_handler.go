package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"rest-api-golang/database"
	"rest-api-golang/models"
	"rest-api-golang/repositories"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

// Repositories
var bookRepo *repositories.BookRepository
var userRepo *repositories.UserRepository
var tokenRepo *repositories.TokenRepository

// InitializeRepositories initializes all repositories
func InitializeRepositories() {
	bookRepo = repositories.NewBookRepository(database.DB)
	userRepo = repositories.NewUserRepository(database.DB)
	tokenRepo = repositories.NewTokenRepository(database.DB)
}

// AuthMiddleware protects endpoints with Bearer token except excluded paths
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow login, health, and docs without token
		if strings.HasPrefix(r.URL.Path, "/api/login") || 
		   strings.HasPrefix(r.URL.Path, "/health") ||
		   strings.HasPrefix(r.URL.Path, "/docs") ||
		   strings.HasPrefix(r.URL.Path, "/swagger") {
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

		// Check token in database
		_, err := tokenRepo.GetTokenByValue(token)
		if err != nil {
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

	// Get user from database
	user, err := userRepo.GetUserByUsername(req.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid username or password",
		})
		return
	}

	// Simple password check (in production, use bcrypt)
	if user.Password != req.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid username or password",
		})
		return
	}

	// Generate token
	tokenValue := randomToken()
	token := &models.Token{
		ID:        uuid.New().String(),
		Token:     tokenValue,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours
		CreatedAt: time.Now(),
		IsRevoked: false,
	}

	// Save token to database
	if err := tokenRepo.CreateToken(token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to create session",
		})
		return
	}

	// Update last login
	userRepo.UpdateLastLogin(user.ID)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token":   tokenValue,
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

	// Revoke token in database
	if err := tokenRepo.RevokeToken(token); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Token not found or already revoked",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}

// GetBooks handles GET /api/books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := bookRepo.GetAllBooks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to fetch books",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    books,
		"count":   len(books),
	})
}

// GetBook handles GET /api/books/{id}
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	book, err := bookRepo.GetBookByID(id)
	if err != nil {
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
	
	if err := bookRepo.CreateBook(book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to create book",
		})
		return
	}

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

	// Get existing book
	book, err := bookRepo.GetBookByID(id)
	if err != nil {
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
	
	if err := bookRepo.UpdateBook(book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to update book",
		})
		return
	}

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

	// Get book before deletion for response
	book, err := bookRepo.GetBookByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Book not found",
		})
		return
	}

	// Soft delete the book
	if err := bookRepo.DeleteBook(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to delete book",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Book deleted successfully",
		"data":    book,
	})
}
