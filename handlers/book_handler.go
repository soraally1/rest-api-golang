package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"rest-api-golang/models"

	"github.com/gorilla/mux"
)

// In-memory storage for books
var books = make(map[string]*models.Book)

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
