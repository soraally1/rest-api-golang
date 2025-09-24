package models

import (
	"time"

	"github.com/google/uuid"
)

// Book represents a book entity
type Book struct {
	ID          string    `json:"id"`
	Judul       string    `json:"judul"`
	Author      string    `json:"author"`
	TahunTerbit int       `json:"tahun_terbit"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateBookRequest represents the request payload for creating a book
type CreateBookRequest struct {
	Judul       string `json:"judul" validate:"required"`
	Author      string `json:"author" validate:"required"`
	TahunTerbit int    `json:"tahun_terbit" validate:"required,min=1000,max=2024"`
}

// UpdateBookRequest represents the request payload for updating a book
type UpdateBookRequest struct {
	Judul       string `json:"judul,omitempty"`
	Author      string `json:"author,omitempty"`
	TahunTerbit int    `json:"tahun_terbit,omitempty"`
}

// NewBook creates a new Book instance
func NewBook(req CreateBookRequest) *Book {
	now := time.Now()
	return &Book{
		ID:          uuid.New().String(),
		Judul:       req.Judul,
		Author:      req.Author,
		TahunTerbit: req.TahunTerbit,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
