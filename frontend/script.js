// API Configuration
const API_BASE_URL = 'http://localhost:8080/api';

// Global Variables
let books = [];
let currentEditingBook = null;
let currentDeletingBook = null;
let authToken = localStorage.getItem('book_api_token') || '';

// DOM Elements
const booksGrid = document.getElementById('booksGrid');
const emptyState = document.getElementById('emptyState');
const loadingSpinner = document.getElementById('loadingSpinner');
const searchInput = document.getElementById('searchInput');
const bookModal = document.getElementById('bookModal');
const deleteModal = document.getElementById('deleteModal');
const bookForm = document.getElementById('bookForm');
const modalTitle = document.getElementById('modalTitle');
const submitBtn = document.getElementById('submitBtn');
const confirmDeleteBtn = document.getElementById('confirmDeleteBtn');
const deleteBookPreview = document.getElementById('deleteBookPreview');
// Auth UI Elements
const loginModal = document.getElementById('loginModal');
const loginForm = document.getElementById('loginForm');
const loginBtn = document.getElementById('loginBtn');
const logoutBtn = document.getElementById('logoutBtn');
const authStatus = document.getElementById('authStatus');
const addBookBtn = document.getElementById('addBookBtn');

// Stats Elements
const totalBooksEl = document.getElementById('totalBooks');
const totalAuthorsEl = document.getElementById('totalAuthors');
const avgYearEl = document.getElementById('avgYear');

// Initialize App
document.addEventListener('DOMContentLoaded', function() {
    loadBooks();
    setupEventListeners();
    updateAuthUI();
});

// Event Listeners
function setupEventListeners() {
    // Search functionality
    searchInput.addEventListener('input', handleSearch);
    
    // Form submission
    bookForm.addEventListener('submit', handleFormSubmit);
    
    // Modal close on outside click
    window.addEventListener('click', function(event) {
        if (event.target === bookModal) {
            closeModal();
        }
        if (event.target === deleteModal) {
            closeDeleteModal();
        }
    });
    
    // Confirm delete
    confirmDeleteBtn.addEventListener('click', handleDeleteConfirm);

    // Login form
    loginForm.addEventListener('submit', handleLoginSubmit);
}

// API Functions
async function fetchBooks() {
    try {
        const response = await fetch(`${API_BASE_URL}/books`, {
            headers: authHeaders()
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        return data.success ? data.data : [];
    } catch (error) {
        console.error('Error fetching books:', error);
        showToast('Failed to load books', 'error');
        return [];
    }
}

async function createBook(bookData) {
    try {
        const response = await fetch(`${API_BASE_URL}/books`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                ...authHeaders(),
            },
            body: JSON.stringify(bookData)
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.message || 'Failed to create book');
        }
        
        return data;
    } catch (error) {
        console.error('Error creating book:', error);
        throw error;
    }
}

async function updateBook(id, bookData) {
    try {
        const response = await fetch(`${API_BASE_URL}/books/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                ...authHeaders(),
            },
            body: JSON.stringify(bookData)
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.message || 'Failed to update book');
        }
        
        return data;
    } catch (error) {
        console.error('Error updating book:', error);
        throw error;
    }
}

async function deleteBook(id) {
    try {
        const response = await fetch(`${API_BASE_URL}/books/${id}`, {
            method: 'DELETE',
            headers: { ...authHeaders() }
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.message || 'Failed to delete book');
        }
        
        return data;
    } catch (error) {
        console.error('Error deleting book:', error);
        throw error;
    }
}

// Auth APIs
async function login(username, password) {
    const resp = await fetch(`${API_BASE_URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    });
    const data = await resp.json();
    if (!resp.ok || !data.success) {
        throw new Error(data.message || 'Login failed');
    }
    return data.token;
}

async function logout() {
    if (!authToken) return;
    await fetch(`${API_BASE_URL}/logout`, {
        method: 'POST',
        headers: { ...authHeaders() }
    });
}

// UI Functions
async function loadBooks() {
    showLoading(true);
    try {
        books = await fetchBooks();
        displayBooks(books);
        updateStats();
    } catch (error) {
        showToast('Failed to load books', 'error');
    } finally {
        showLoading(false);
    }
}

function displayBooks(booksToShow) {
    if (booksToShow.length === 0) {
        booksGrid.style.display = 'none';
        emptyState.style.display = 'block';
        return;
    }
    
    booksGrid.style.display = 'grid';
    emptyState.style.display = 'none';
    
    booksGrid.innerHTML = booksToShow.map(book => `
        <div class="book-card">
            <h3 class="book-title">${escapeHtml(book.judul)}</h3>
            <div class="book-author">
                <i class="fas fa-user"></i>
                ${escapeHtml(book.author)}
            </div>
            <div class="book-year">${book.tahun_terbit}</div>
            <div class="book-actions">
                <button class="btn btn-warning" onclick="openEditModal('${book.id}')">
                    <i class="fas fa-edit"></i> Edit
                </button>
                <button class="btn btn-danger" onclick="openDeleteModal('${book.id}')">
                    <i class="fas fa-trash"></i> Delete
                </button>
            </div>
        </div>
    `).join('');
}

function updateStats() {
    const totalBooks = books.length;
    const uniqueAuthors = new Set(books.map(book => book.author)).size;
    const avgYear = books.length > 0 
        ? Math.round(books.reduce((sum, book) => sum + book.tahun_terbit, 0) / books.length)
        : 0;
    
    totalBooksEl.textContent = totalBooks;
    totalAuthorsEl.textContent = uniqueAuthors;
    avgYearEl.textContent = avgYear;
}

function showLoading(show) {
    loadingSpinner.style.display = show ? 'block' : 'none';
    booksGrid.style.display = show ? 'none' : 'grid';
    emptyState.style.display = show ? 'none' : (books.length === 0 ? 'block' : 'none');
}

// Search Functionality
function handleSearch() {
    const searchTerm = searchInput.value.toLowerCase().trim();
    
    if (searchTerm === '') {
        displayBooks(books);
        return;
    }
    
    const filteredBooks = books.filter(book => 
        book.judul.toLowerCase().includes(searchTerm) ||
        book.author.toLowerCase().includes(searchTerm)
    );
    
    displayBooks(filteredBooks);
}

// Modal Functions
function openCreateModal() {
    currentEditingBook = null;
    modalTitle.textContent = 'Add New Book';
    submitBtn.innerHTML = '<i class="fas fa-save"></i> Save Book';
    bookForm.reset();
    bookModal.style.display = 'block';
    document.getElementById('bookTitle').focus();
}

function openEditModal(bookId) {
    const book = books.find(b => b.id === bookId);
    if (!book) return;
    
    currentEditingBook = book;
    modalTitle.textContent = 'Edit Book';
    submitBtn.innerHTML = '<i class="fas fa-save"></i> Update Book';
    
    // Populate form
    document.getElementById('bookTitle').value = book.judul;
    document.getElementById('bookAuthor').value = book.author;
    document.getElementById('bookYear').value = book.tahun_terbit;
    
    bookModal.style.display = 'block';
    document.getElementById('bookTitle').focus();
}

function closeModal() {
    bookModal.style.display = 'none';
    currentEditingBook = null;
    bookForm.reset();
}

function openDeleteModal(bookId) {
    const book = books.find(b => b.id === bookId);
    if (!book) return;
    
    currentDeletingBook = book;
    deleteBookPreview.innerHTML = `
        <strong>${escapeHtml(book.judul)}</strong><br>
        <small>by ${escapeHtml(book.author)} (${book.tahun_terbit})</small>
    `;
    deleteModal.style.display = 'block';
}

function closeDeleteModal() {
    deleteModal.style.display = 'none';
    currentDeletingBook = null;
}

// Login Modal functions
function openLoginModal() {
    loginModal.style.display = 'block';
    document.getElementById('username').focus();
}

function closeLoginModal() {
    loginModal.style.display = 'none';
    loginForm.reset();
}

async function handleLoginSubmit(e) {
    e.preventDefault();
    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value.trim();
    const btn = document.getElementById('loginSubmitBtn');
    try {
        btn.disabled = true;
        btn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Logging in...';
        const token = await login(username, password);
        authToken = token;
        localStorage.setItem('book_api_token', token);
        showToast('Login successful', 'success');
        closeLoginModal();
        updateAuthUI();
        await loadBooks();
    } catch (err) {
        showToast(err.message, 'error');
    } finally {
        btn.disabled = false;
        btn.innerHTML = '<i class="fas fa-sign-in-alt"></i> Login';
    }
}

async function handleLogout() {
    try {
        await logout();
    } catch (_) {}
    authToken = '';
    localStorage.removeItem('book_api_token');
    updateAuthUI();
    await loadBooks();
}

function updateAuthUI() {
    const loggedIn = !!authToken;
    authStatus.textContent = loggedIn ? 'Logged in' : 'Not logged in';
    loginBtn.style.display = loggedIn ? 'none' : 'inline-flex';
    logoutBtn.style.display = loggedIn ? 'inline-flex' : 'none';
    addBookBtn.disabled = !loggedIn;
    addBookBtn.title = loggedIn ? '' : 'Login to add a book';
}

function authHeaders() {
    return authToken ? { 'Authorization': `Bearer ${authToken}` } : {};
}

// Form Handling
async function handleFormSubmit(event) {
    event.preventDefault();
    
    const formData = new FormData(bookForm);
    const bookData = {
        judul: formData.get('judul').trim(),
        author: formData.get('author').trim(),
        tahun_terbit: parseInt(formData.get('tahun_terbit'))
    };
    
    // Validation
    if (!bookData.judul || !bookData.author || !bookData.tahun_terbit) {
        showToast('Please fill in all fields', 'warning');
        return;
    }
    
    if (bookData.tahun_terbit < 1000 || bookData.tahun_terbit > 2024) {
        showToast('Publication year must be between 1000 and 2024', 'warning');
        return;
    }
    
    try {
        submitBtn.disabled = true;
        submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Saving...';
        
        if (currentEditingBook) {
            await updateBook(currentEditingBook.id, bookData);
            showToast('Book updated successfully!', 'success');
        } else {
            await createBook(bookData);
            showToast('Book created successfully!', 'success');
        }
        
        closeModal();
        await loadBooks();
        
    } catch (error) {
        showToast(error.message || 'Failed to save book', 'error');
    } finally {
        submitBtn.disabled = false;
        submitBtn.innerHTML = currentEditingBook 
            ? '<i class="fas fa-save"></i> Update Book'
            : '<i class="fas fa-save"></i> Save Book';
    }
}

async function handleDeleteConfirm() {
    if (!currentDeletingBook) return;
    
    try {
        confirmDeleteBtn.disabled = true;
        confirmDeleteBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Deleting...';
        
        await deleteBook(currentDeletingBook.id);
        showToast('Book deleted successfully!', 'success');
        
        closeDeleteModal();
        await loadBooks();
        
    } catch (error) {
        showToast(error.message || 'Failed to delete book', 'error');
    } finally {
        confirmDeleteBtn.disabled = false;
        confirmDeleteBtn.innerHTML = '<i class="fas fa-trash"></i> Delete Book';
    }
}

// Utility Functions
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function showToast(message, type = 'success') {
    const toastContainer = document.getElementById('toastContainer');
    
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    
    const icon = type === 'success' ? 'fas fa-check-circle' :
                 type === 'error' ? 'fas fa-exclamation-circle' :
                 'fas fa-exclamation-triangle';
    
    toast.innerHTML = `
        <i class="${icon}"></i>
        <span>${message}</span>
    `;
    
    toastContainer.appendChild(toast);
    
    // Auto remove after 5 seconds
    setTimeout(() => {
        if (toast.parentNode) {
            toast.parentNode.removeChild(toast);
        }
    }, 5000);
}

// Keyboard shortcuts
document.addEventListener('keydown', function(event) {
    // Escape key to close modals
    if (event.key === 'Escape') {
        if (bookModal.style.display === 'block') {
            closeModal();
        }
        if (deleteModal.style.display === 'block') {
            closeDeleteModal();
        }
    }
    
    // Ctrl/Cmd + K to focus search
    if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
        event.preventDefault();
        searchInput.focus();
    }
    
    // Ctrl/Cmd + N to create new book
    if ((event.ctrlKey || event.metaKey) && event.key === 'n') {
        event.preventDefault();
        openCreateModal();
    }
});
