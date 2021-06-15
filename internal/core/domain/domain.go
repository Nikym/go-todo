package domain

import uuid2 "github.com/google/uuid"

// Entry object describes a to-do instance.
type Entry struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// NewEntry returns a pointer to a new Entry object.
func NewEntry(title, description string) *Entry {
	id := uuid2.NewString()
	return &Entry{
		ID:          id,
		Title:       title,
		Description: description,
		Done:        false,
	}
}
