package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"rest-api-golang/models"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

// GetAllBooks retrieves all non-deleted books
func (r *BookRepository) GetAllBooks() ([]*models.Book, error) {
	query := `
		SELECT id, judul, author, tahun_terbit, created_at, updated_at, deleted_at
		FROM books 
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query books: %w", err)
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		book := &models.Book{}
		err := rows.Scan(
			&book.ID,
			&book.Judul,
			&book.Author,
			&book.TahunTerbit,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, book)
	}

	return books, nil
}

// GetBookByID retrieves a book by ID
func (r *BookRepository) GetBookByID(id string) (*models.Book, error) {
	query := `
		SELECT id, judul, author, tahun_terbit, created_at, updated_at, deleted_at
		FROM books 
		WHERE id = $1 AND deleted_at IS NULL`

	book := &models.Book{}
	err := r.db.QueryRow(query, id).Scan(
		&book.ID,
		&book.Judul,
		&book.Author,
		&book.TahunTerbit,
		&book.CreatedAt,
		&book.UpdatedAt,
		&book.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("book not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return book, nil
}

// CreateBook creates a new book
func (r *BookRepository) CreateBook(book *models.Book) error {
	query := `
		INSERT INTO books (id, judul, author, tahun_terbit, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query,
		book.ID,
		book.Judul,
		book.Author,
		book.TahunTerbit,
		book.CreatedAt,
		book.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	return nil
}

// UpdateBook updates an existing book
func (r *BookRepository) UpdateBook(book *models.Book) error {
	query := `
		UPDATE books 
		SET judul = $2, author = $3, tahun_terbit = $4, updated_at = $5
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query,
		book.ID,
		book.Judul,
		book.Author,
		book.TahunTerbit,
		book.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

// DeleteBook soft deletes a book
func (r *BookRepository) DeleteBook(id string) error {
	query := `
		UPDATE books 
		SET deleted_at = $2
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

// HardDeleteBook permanently deletes a book
func (r *BookRepository) HardDeleteBook(id string) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}
