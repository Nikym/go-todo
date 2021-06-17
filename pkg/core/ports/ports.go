package ports

import "github.com/Nikym/go-todo/pkg/core/domain"

// EntryRepository is the interface for the repository port handling the
// retrieval and storage of to-do entries.
type EntryRepository interface {
	Get(id string) (*domain.Entry, error)
	Save(entry *domain.Entry) error
	Delete(id string) error
	Update(id string, entry *domain.Entry) error
}

// EntryService is the interface for the driver port handling the
// interactions with entries (domain.Entry)
type EntryService interface {
	Get(id string) (*domain.Entry, error)
	Create(title, description string) (*domain.Entry, error)
	Update(id string, entry *domain.Entry) error
	Delete(id string) error
}
