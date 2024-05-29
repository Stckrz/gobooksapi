package main

import (
	"time"
)

type Book struct {
	ID               int       `json:"id"`
	CreatedAt        time.Time `json:"timestamp"`
	Title            string    `json:"title"`
	Author           string    `json:"author"`
	OriginalLanguage string    `json:"original_language"`
	AuthorId         int       `json:"author_id"`
}
