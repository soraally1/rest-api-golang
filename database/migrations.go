package database

import (
	"fmt"
	"log"
)

// CreateTables creates all necessary tables
func CreateTables() error {
	// Create books table
	booksTable := `
	CREATE TABLE IF NOT EXISTS books (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		judul VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		tahun_terbit INTEGER NOT NULL CHECK (tahun_terbit >= 1000 AND tahun_terbit <= 2024),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP WITH TIME ZONE NULL
	);`

	// Create users table for authentication
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		email VARCHAR(255),
		role VARCHAR(20) DEFAULT 'user',
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		last_login TIMESTAMP WITH TIME ZONE NULL
	);`

	// Create tokens table for session management
	tokensTable := `
	CREATE TABLE IF NOT EXISTS tokens (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		token VARCHAR(255) UNIQUE NOT NULL,
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		is_revoked BOOLEAN DEFAULT false
	);`

	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);",
		"CREATE INDEX IF NOT EXISTS idx_books_tahun_terbit ON books(tahun_terbit);",
		"CREATE INDEX IF NOT EXISTS idx_books_created_at ON books(created_at);",
		"CREATE INDEX IF NOT EXISTS idx_tokens_token ON tokens(token);",
		"CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_tokens_expires_at ON tokens(expires_at);",
	}

	// Create trigger to update updated_at column
	updateTrigger := `
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ language 'plpgsql';

	DROP TRIGGER IF EXISTS update_books_updated_at ON books;
	CREATE TRIGGER update_books_updated_at
		BEFORE UPDATE ON books
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();

	DROP TRIGGER IF EXISTS update_users_updated_at ON users;
	CREATE TRIGGER update_users_updated_at
		BEFORE UPDATE ON users
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();
	`

	tables := []string{booksTable, usersTable, tokensTable}
	
	// Execute table creation
	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	// Execute indexes
	for _, index := range indexes {
		if _, err := DB.Exec(index); err != nil {
			log.Printf("Warning: failed to create index: %v", err)
		}
	}

	// Execute trigger
	if _, err := DB.Exec(updateTrigger); err != nil {
		log.Printf("Warning: failed to create trigger: %v", err)
	}

	log.Println("✅ Database tables created successfully")
	return nil
}

// SeedData inserts initial data if tables are empty
func SeedData() error {
	// Check if users table is empty
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check users count: %w", err)
	}

	// Insert default users if table is empty
	if count == 0 {
		seedUsers := `
		INSERT INTO users (username, password, email, role) VALUES
		('admin', 'admin123', 'admin@example.com', 'admin'),
		('user', 'user123', 'user@example.com', 'user')
		ON CONFLICT (username) DO NOTHING;`

		if _, err := DB.Exec(seedUsers); err != nil {
			return fmt.Errorf("failed to seed users: %w", err)
		}
		log.Println("✅ Default users seeded successfully")
	}

	return nil
}
