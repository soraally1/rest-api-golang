# 📚 Book Management Frontend

Modern, responsive web application untuk mengelola koleksi buku dengan interface yang user-friendly.

## 🚀 Fitur

- **📖 CRUD Operations** - Create, Read, Update, Delete buku
- **🔍 Real-time Search** - Pencarian berdasarkan judul dan author
- **📊 Statistics Dashboard** - Statistik total buku, author unik, dan rata-rata tahun
- **🎨 Modern UI/UX** - Design yang menarik dengan animasi smooth
- **📱 Responsive Design** - Optimal di desktop, tablet, dan mobile
- **⚡ Real-time Updates** - Update data secara real-time
- **🔔 Toast Notifications** - Notifikasi sukses/error yang informatif
- **⌨️ Keyboard Shortcuts** - Shortcut untuk efisiensi

## 🎯 Fitur UI/UX

### Dashboard
- **Stats Cards** - Menampilkan total buku, author unik, dan rata-rata tahun terbit
- **Search Bar** - Pencarian real-time dengan highlight
- **Book Grid** - Tampilan kartu buku yang menarik
- **Empty State** - State ketika belum ada buku

### Book Cards
- **Book Information** - Judul, author, tahun terbit
- **Action Buttons** - Edit dan Delete dengan konfirmasi
- **Hover Effects** - Animasi hover yang smooth
- **Responsive Layout** - Adaptif untuk berbagai ukuran layar

### Modals
- **Create/Edit Modal** - Form untuk menambah/edit buku
- **Delete Confirmation** - Konfirmasi sebelum menghapus
- **Form Validation** - Validasi input yang real-time
- **Loading States** - Indikator loading saat proses

## 🛠️ Teknologi

- **HTML5** - Semantic markup
- **CSS3** - Modern styling dengan Flexbox/Grid
- **Vanilla JavaScript** - No framework dependencies
- **Font Awesome** - Icons yang konsisten
- **Responsive Design** - Mobile-first approach

## 📁 Struktur File

```
frontend/
├── index.html          # Main HTML file
├── styles.css          # CSS styling
├── script.js           # JavaScript functionality
└── README.md           # Documentation
```

## 🚀 Cara Menjalankan

1. **Pastikan backend Go sudah berjalan:**
   ```bash
   cd ../
   go run main.go
   ```

2. **Buka frontend:**
   - Buka `frontend/index.html` di browser
   - Atau gunakan live server untuk development

3. **Akses aplikasi:**
   - URL: `http://localhost:8080` (jika menggunakan live server)
   - Atau buka file `index.html` langsung di browser

## 🎨 Design Features

### Color Scheme
- **Primary**: Gradient blue-purple (#667eea → #764ba2)
- **Success**: Green (#51cf66)
- **Warning**: Yellow (#ffd43b)
- **Danger**: Red (#ff6b6b)
- **Background**: Gradient background

### Typography
- **Font**: Segoe UI, system fonts
- **Hierarchy**: Clear heading structure
- **Readability**: Optimal contrast ratios

### Animations
- **Hover Effects**: Smooth transitions
- **Modal Animations**: Slide-in effects
- **Loading Spinners**: Rotating animations
- **Toast Notifications**: Slide-in from right

## ⌨️ Keyboard Shortcuts

- **Ctrl/Cmd + K** - Focus search bar
- **Ctrl/Cmd + N** - Open create book modal
- **Escape** - Close any open modal

## 📱 Responsive Breakpoints

- **Desktop**: > 768px (Grid layout)
- **Tablet**: 768px - 480px (Single column)
- **Mobile**: < 480px (Optimized for touch)

## 🔧 Customization

### Mengubah API URL
Edit di `script.js`:
```javascript
const API_BASE_URL = 'http://localhost:8080/api';
```

### Mengubah Theme Colors
Edit di `styles.css`:
```css
:root {
    --primary-color: #667eea;
    --secondary-color: #764ba2;
    --success-color: #51cf66;
    --warning-color: #ffd43b;
    --danger-color: #ff6b6b;
}
```

## 🐛 Troubleshooting

### CORS Issues
Jika ada masalah CORS, pastikan backend Go sudah mengaktifkan CORS middleware.

### API Connection
- Pastikan backend berjalan di `http://localhost:8080`
- Check browser console untuk error messages
- Pastikan tidak ada firewall yang memblokir

### Performance
- Gunakan browser modern untuk performa optimal
- Clear browser cache jika ada masalah styling

## 📈 Future Enhancements

- [ ] Pagination untuk buku banyak
- [ ] Filter berdasarkan tahun/author
- [ ] Export data ke CSV/PDF
- [ ] Dark mode toggle
- [ ] Book cover upload
- [ ] Advanced search dengan multiple criteria
- [ ] Book categories/tags
- [ ] Reading progress tracking
