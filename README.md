📚 Book Management REST API
REST API modern untuk mengelola koleksi buku dengan authentication menggunakan Go, PostgreSQL, dan Gorilla Mux.

🎯 Features
✅ CRUD Operations - Create, Read, Update, Delete books
✅ Authentication - Bearer token-based authentication
✅ PostgreSQL Database - Persistent data storage
✅ Auto Migrations - Automatic table creation and seeding
✅ Soft Delete - Safe deletion with recovery option
✅ CORS Support - Cross-origin resource sharing
✅ Input Validation - Data validation and error handling
✅ API Documentation - Built-in HTML documentation
✅ Health Check - Endpoint untuk monitoring
🛠️ Technology Stack
Go 1.21+ - Programming language
Gorilla Mux - HTTP router and URL matcher
PostgreSQL - Relational database
lib/pq - PostgreSQL driver
Google UUID - Unique identifier generation
YAML - Configuration file format
📋 Prerequisites
Go 1.21 atau lebih baru
PostgreSQL 12+ terinstall dan berjalan
Git (optional)
🚀 Installation & Setup
1. Clone Repository
bash
git clone <repository-url>
cd rest-api-golang
2. Install Dependencies
bash
go mod tidy
3. Setup PostgreSQL Database
Opsi A: Manual Setup
sql
-- Login ke PostgreSQL
psql -U postgres

-- Create database
CREATE DATABASE bookdb;

-- Keluar dari psql
\q
Opsi B: Menggunakan DBeaver
Buka DBeaver
Klik "New Database Connection" (icon plug)
Pilih PostgreSQL → Next
Masukkan connection details:
Host: localhost
Port: 5432
Database: postgres
Username: postgres
Password:
Test Connection → Finish
Right-click connection → SQL Editor → New SQL Script
Jalankan: CREATE DATABASE bookdb;
4. Configure Environment Variables
Copy file 
env.example
 ke .env atau set environment variables:

bash
# Windows CMD
set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=postgres
set DB_PASSWORD=password
set DB_NAME=bookdb
set DB_SSLMODE=disable

# Windows PowerShell
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="password"
$env:DB_NAME="bookdb"
$env:DB_SSLMODE="disable"

# Linux/Mac
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=bookdb
export DB_SSLMODE=disable
Default Values (jika tidak di-set):

DB_HOST: localhost
DB_PORT: 5432
DB_USER: postgres
DB_PASSWORD: password
DB_NAME: bookdb
DB_SSLMODE: disable
5. Run Application
bash
go run main.go
Server akan berjalan di http://localhost:8080

🗄️ Database Schema
Tables
1. books
sql
CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    judul VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    tahun_terbit INTEGER NOT NULL CHECK (tahun_terbit >= 1000 AND tahun_terbit <= 2024),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);
2. users
sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    role VARCHAR(20) DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE NULL
);
3. tokens
sql
CREATE TABLE tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_revoked BOOLEAN DEFAULT false
);
Indexes
idx_books_author - Index pada kolom author
idx_books_tahun_terbit - Index pada tahun terbit
idx_books_created_at - Index pada created_at
idx_tokens_token - Index pada token
idx_tokens_user_id - Index pada user_id
idx_tokens_expires_at - Index pada expires_at
Database Triggers
Auto-update updated_at column untuk tabel books dan users
📡 API Endpoints
Base URL: http://localhost:8080

Authentication
1. Login
POST /api/login
Request Body:

json
{
  "username": "admin",
  "password": "admin123"
}
Response:

json
{
  "success": true,
  "token": "Xyz123AbC456..."
}
2. Logout
POST /api/logout
Headers: Authorization: Bearer <token>
Response:

json
{
  "success": true,
  "message": "Logged out successfully"
}
Books Management (Requires Authentication)
3. Get All Books
GET /api/books
Headers: Authorization: Bearer <token>
Response:

json
{
  "success": true,
  "data": [
    {
      "id": "uuid-string",
      "judul": "To Kill a Mockingbird",
      "author": "Harper Lee",
      "tahun_terbit": 1960,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ],
  "count": 1
}
4. Get Book by ID
GET /api/books/{id}
Headers: Authorization: Bearer <token>
5. Create Book
POST /api/books
Headers: Authorization: Bearer <token>
Request Body:

json
{
  "judul": "1984",
  "author": "George Orwell",
  "tahun_terbit": 1949
}
6. Update Book
PUT /api/books/{id}
Headers: Authorization: Bearer <token>
Request Body:

json
{
  "judul": "Nineteen Eighty-Four",
  "author": "George Orwell",
  "tahun_terbit": 1949
}
7. Delete Book
DELETE /api/books/{id}
Headers: Authorization: Bearer <token>
Utility Endpoints
8. Health Check
GET /health
Response:

json
{
  "status": "healthy",
  "message": "Server is running"
}
9. API Documentation
GET /docs
Menampilkan dokumentasi HTML interaktif di browser.

👥 Default Users
Aplikasi akan otomatis membuat 2 user saat pertama kali dijalankan:

Username	Password	Role	Email
admin	admin123	admin	admin@example.com
user	user123	user	user@example.com
🧪 Testing API
Menggunakan cURL
1. Login untuk mendapatkan token
bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"admin\",\"password\":\"admin123\"}"
2. Get all books
bash
curl -X GET http://localhost:8080/api/books \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
3. Create new book
bash
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d "{\"judul\":\"The Great Gatsby\",\"author\":\"F. Scott Fitzgerald\",\"tahun_terbit\":1925}"
4. Update book
bash
curl -X PUT http://localhost:8080/api/books/{book-id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d "{\"judul\":\"Updated Title\"}"
5. Delete book
bash
curl -X DELETE http://localhost:8080/api/books/{book-id} \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
Menggunakan Postman
Import Collection: Buat collection baru
Set Base URL: http://localhost:8080
Login: POST ke /api/login dengan credentials
Save Token: Copy token dari response
Set Authorization: Pada requests lain, tambahkan header:
Key: Authorization
Value: Bearer <your-token>
💾 Melihat Database di DBeaver
Setup Connection
Buka DBeaver → Click icon "New Database Connection" atau Ctrl+Shift+N
Pilih PostgreSQL dari list database → Click Next
Masukkan Connection Details:
Host: localhost
Port: 5432
Database: bookdb
Username: postgres
Password: password
Test Connection:
Click tombol "Test Connection"
Jika prompt download driver, click "Download"
Tunggu sampai muncul "Connected"
Click Finish untuk save connection
Viewing Data
Expand Connection: bookdb → Schemas → public → Tables
View Books Table:
Right-click table books
Select "View Data" atau double-click
Semua data buku akan tampil
View Users Table:
Right-click table users
Select "View Data"
View Tokens Table:
Right-click table tokens
Select "View Data"
Lihat active sessions dan expired tokens
Useful Queries di DBeaver
sql
-- Lihat semua buku yang tidak dihapus
SELECT * FROM books WHERE deleted_at IS NULL;

-- Lihat semua users dengan login terakhir
SELECT username, email, last_login FROM users ORDER BY last_login DESC;

-- Lihat active tokens
SELECT t.token, u.username, t.expires_at 
FROM tokens t 
JOIN users u ON t.user_id = u.id 
WHERE t.is_revoked = false AND t.expires_at > NOW();

-- Count books per author
SELECT author, COUNT(*) as total 
FROM books 
WHERE deleted_at IS NULL 
GROUP BY author 
ORDER BY total DESC;

-- Books by year range
SELECT * FROM books 
WHERE tahun_terbit BETWEEN 1900 AND 2000 
AND deleted_at IS NULL
ORDER BY tahun_terbit DESC;
🏗️ Project Structure
rest-api-golang/
├── main.go                     # Entry point & router setup
├── go.mod                      # Go dependencies
├── go.sum                      # Dependency checksums
├── config.yaml                 # User configuration
├── env.example                 # Environment variables example
├── Dockerfile                  # Docker configuration
├── docker-compose.yml          # Docker compose setup
│
├── models/                     # Data models
│   └── book.go                 # Book & User models
│
├── handlers/                   # HTTP handlers
│   └── book_handler.go         # CRUD & Auth handlers
│
├── database/                   # Database layer
│   ├── database.go             # DB connection
│   └── migrations.go           # Schema & seeding
│
└── repositories/               # Data access layer
    ├── book_repository.go      # Book CRUD operations
    ├── user_repository.go      # User operations
    └── token_repository.go     # Token management
🔒 Security Notes
⚠️ Development Only - Aplikasi ini untuk development/learning:

Password Storage: Plain text (gunakan bcrypt di production)
Token: Simple random string (gunakan JWT di production)
CORS: Open untuk semua origins
SSL: Disabled (enable di production)
Input Validation: Basic validation only
Untuk Production:

Hash passwords dengan bcrypt
Gunakan JWT tokens
Enable SSL/TLS
Restrict CORS origins
Add rate limiting
Implement proper logging
Add request validation middleware
🐛 Troubleshooting
Database Connection Failed
Error: failed to connect to database
Solution:

Pastikan PostgreSQL berjalan
Check credentials di environment variables
Verify database bookdb sudah dibuat
Test connection: psql -U postgres -d bookdb
Port Already in Use
Error: listen tcp :8080: bind: address already in use
Solution:

Kill process di port 8080: netstat -ano | findstr :8080
Atau ubah port di code
CORS Errors
Access to fetch blocked by CORS policy
Solution:

CORS middleware sudah enabled
Check frontend menggunakan http://localhost:8080
Clear browser cache
Authentication Failed
{"success": false, "message": "Token expired or invalid"}
Solution:

Login ulang untuk mendapatkan token baru
Token expired setelah 24 jam
Pastikan header: Authorization: Bearer <token>
📈 Performance Tips
Database Indexes: Sudah dibuat otomatis untuk query performance
Connection Pooling: PostgreSQL driver menangani otomatis
Soft Delete: Gunakan WHERE deleted_at IS NULL untuk filtering
Token Cleanup: Periodic cleanup untuk expired tokens (implement cron job)
🚢 Deployment
Build Binary
bash
# Build untuk Windows
go build -o book-api.exe main.go

# Build untuk Linux
GOOS=linux GOARCH=amd64 go build -o book-api main.go

# Run
./book-api
Docker (Coming Soon)
bash
docker-compose up -d
📝 Development Roadmap
 JWT token implementation
 Password hashing dengan bcrypt
 Pagination untuk list endpoints
 Search & filter functionality
 Rate limiting
 API versioning
 Swagger/OpenAPI documentation
 Unit tests
 Integration tests
 Docker deployment
 CI/CD pipeline
📄 License
MIT License - Free to use for learning and development.

🙏 Credits
Go Team - Excellent programming language
Gorilla Toolkit - Powerful HTTP libraries
PostgreSQL - Robust database system
DBeaver - Amazing database management tool
Happy Coding! 🚀

Untuk pertanyaan atau issue, silakan buat issue di repository atau hubungi maintainer.

Summary
Saya telah membuat README yang komprehensif dengan:

✅ Database Setup - Langkah-langkah setup PostgreSQL dan DBeaver
✅ Schema Documentation - Detail struktur tabel books, users, tokens
✅ API Endpoints - Semua endpoint dengan contoh request/response
✅ DBeaver Guide - Cara connect dan query database
✅ Testing Examples - Curl commands dan Postman setup
✅ Troubleshooting - Common issues dan solusinya
✅ Environment Variables - Konfigurasi untuk Windows/Linux/Mac
