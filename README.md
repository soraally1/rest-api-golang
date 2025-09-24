# 📚 Complete Book Management System

Full-stack application untuk mengelola koleksi buku dengan backend Go dan frontend modern.

## 🏗️ Architecture

```
rest-api-golang/
├── main.go                    # Go backend server
├── go.mod                     # Go dependencies
├── models/
│   └── book.go               # Book data models
├── handlers/
│   └── book_handler.go       # API handlers
├── frontend/                 # Web frontend
│   ├── index.html            # Main HTML
│   ├── styles.css            # CSS styling
│   ├── script.js             # JavaScript
│   └── README.md             # Frontend docs
└── README.md                 # This file
```

## 🚀 Quick Start

### 1. Backend (Go API)
```bash
# Install dependencies
go mod tidy

# Run backend server
go run main.go
```
Backend akan berjalan di: `http://localhost:8080`

### 2. Frontend (Web App)
```bash
# Buka frontend/index.html di browser
# Atau gunakan live server untuk development
```
Frontend akan berjalan di: `http://localhost:3000` (jika menggunakan live server)

## 📋 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/books` | Get all books |
| POST | `/api/books` | Create new book |
| GET | `/api/books/{id}` | Get book by ID |
| PUT | `/api/books/{id}` | Update book |
| DELETE | `/api/books/{id}` | Delete book |
| GET | `/health` | Health check |

## 🎯 Features

### Backend Features
- ✅ RESTful API dengan 4 endpoint utama
- ✅ CRUD operations lengkap
- ✅ Data validation
- ✅ Error handling
- ✅ CORS support
- ✅ UUID generation
- ✅ In-memory storage

### Frontend Features
- ✅ Modern responsive UI
- ✅ Real-time search
- ✅ Statistics dashboard
- ✅ Modal forms
- ✅ Toast notifications
- ✅ Loading states
- ✅ Keyboard shortcuts
- ✅ Mobile-friendly design

## 🛠️ Technology Stack

### Backend
- **Go 1.21+** - Programming language
- **Gorilla Mux** - HTTP router
- **Google UUID** - Unique ID generation
- **JSON** - Data format

### Frontend
- **HTML5** - Markup
- **CSS3** - Styling dengan modern features
- **Vanilla JavaScript** - No framework dependencies
- **Font Awesome** - Icons
- **Responsive Design** - Mobile-first

## 📊 Data Model

```json
{
  "id": "uuid-string",
  "judul": "Book Title",
  "author": "Author Name", 
  "tahun_terbit": 2024,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

## 🧪 Testing

### Backend Testing
```bash
# Health check
curl http://localhost:8080/health

# Get all books
curl http://localhost:8080/api/books

# Create book
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -d '{"judul":"1984","author":"George Orwell","tahun_terbit":1949}'
```

### Frontend Testing
1. Buka `frontend/index.html` di browser
2. Test semua CRUD operations
3. Test search functionality
4. Test responsive design

## 🔧 Configuration

### Backend Configuration
- **Port**: 8080 (default)
- **CORS**: Enabled for all origins
- **Storage**: In-memory (data hilang saat restart)

### Frontend Configuration
- **API URL**: `http://localhost:8080/api`
- **Theme**: Blue-purple gradient
- **Responsive**: Mobile-first design

## 📱 Screenshots

### Desktop View
- Modern dashboard dengan stats cards
- Book grid layout
- Modal forms untuk create/edit
- Toast notifications

### Mobile View
- Responsive single-column layout
- Touch-friendly buttons
- Optimized modal sizing
- Swipe-friendly interface

## 🚀 Deployment

### Backend Deployment
```bash
# Build binary
go build -o book-api main.go

# Run binary
./book-api
```

### Frontend Deployment
- Upload `frontend/` folder ke web server
- Atau deploy ke static hosting (Netlify, Vercel, etc.)
- Update API URL jika backend di server berbeda

## 🔒 Security Notes

- **CORS**: Currently open for development
- **Validation**: Basic input validation
- **Storage**: In-memory (not persistent)
- **Authentication**: None (for demo purposes)

## 📈 Performance

### Backend
- **Memory Usage**: Minimal (in-memory storage)
- **Response Time**: < 100ms for most operations
- **Concurrency**: Go's goroutines handle multiple requests

### Frontend
- **Load Time**: < 2 seconds
- **Bundle Size**: Minimal (no frameworks)
- **Rendering**: Optimized DOM updates
- **Caching**: Browser caching for static assets

## 🐛 Troubleshooting

### Common Issues

1. **CORS Errors**
   - Pastikan backend CORS middleware aktif
   - Check browser console untuk error details

2. **API Connection Failed**
   - Pastikan backend berjalan di port 8080
   - Check firewall settings
   - Verify API URL di frontend

3. **Styling Issues**
   - Clear browser cache
   - Check CSS file path
   - Verify Font Awesome CDN

4. **JavaScript Errors**
   - Check browser console
   - Verify API responses
   - Check network tab untuk failed requests

## 📝 Development Notes

### Backend Development
- Add database integration untuk persistence
- Implement proper logging
- Add authentication/authorization
- Add API documentation (Swagger)

### Frontend Development
- Add unit tests
- Implement state management
- Add PWA features
- Optimize bundle size

## 🤝 Contributing

1. Fork repository
2. Create feature branch
3. Make changes
4. Test thoroughly
5. Submit pull request

## 📄 License

MIT License - feel free to use for learning and development.

## 🙏 Acknowledgments

- Go community untuk excellent documentation
- Font Awesome untuk beautiful icons
- Modern CSS techniques untuk responsive design
- Vanilla JavaScript untuk lightweight solution