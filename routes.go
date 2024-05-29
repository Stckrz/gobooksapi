package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type BookHandler struct {
}

func (b BookHandler) listBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var u Book
		if err := rows.Scan(&u.ID, &u.CreatedAt, &u.Title, &u.Author, &u.OriginalLanguage, &u.AuthorId); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		books = append(books, u)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(books)
}

func (b BookHandler) getBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	rows, err := db.Query("SELECT * FROM books WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()
	var books []Book
	for rows.Next() {
		var u Book
		if err := rows.Scan(&u.ID, &u.CreatedAt, &u.Title, &u.Author, &u.OriginalLanguage, &u.AuthorId); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		books = append(books, u)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)

}

func (b BookHandler) postBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var newId int
	err := db.QueryRow("INSERT INTO books (title, author, original_language, author_id) VALUES ($1, $2, $3, $4) RETURNING id", book.Title, book.Author, book.OriginalLanguage, book.AuthorId).Scan(&newId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	book.ID = newId
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)

}

func (b BookHandler) editBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "no object in request body", 500)
		return
	}
	rows, err := db.Exec("UPDATE books SET title = $1, author = $2, original_language = $3, author_id = $4 WHERE id = $5", book.Title, book.Author, book.OriginalLanguage, book.AuthorId, id)
	if err != nil {
		// http.Error(w, "invalid object", 500)
		http.Error(w, err.Error(), 500)
		return
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Item not found", 500)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item updated successfully"))
}
