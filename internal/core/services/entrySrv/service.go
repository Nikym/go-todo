package entrySrv

import (
	"errors"
	"github.com/Nikym/go-todo/internal/core/domain"
	"github.com/Nikym/go-todo/internal/core/ports"
)

type service struct {
	entryRepository ports.EntryRepository
}

// New returns a pointer to a new service object.
func New(repository ports.EntryRepository) *service {
	return &service{
		entryRepository: repository,
	}
}

// Get returns the domain.Entry object with the given UUID.
func (srv *service) Get(id string) (*domain.Entry, error) {
	entry, err := srv.entryRepository.Get(id)
	if err != nil {
		return &domain.Entry{}, errors.New("retrieving entry from repository failed")
	}

	return entry, nil
}

// Create makes a new domain.Entry object and saves it to the repository.
func (srv *service) Create(title, description string) (*domain.Entry, error) {
	if len(title) < 3 {
		return &domain.Entry{}, errors.New("title must consist of 3 characters or more")
	}

	entry := domain.NewEntry(title, description)
	if err := srv.entryRepository.Save(entry); err != nil {
		return &domain.Entry{}, errors.New("saving entry to repository failed")
	}

	return entry, nil
}
