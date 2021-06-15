package entrySrv

import "github.com/Nikym/go-todo/internal/core/domain"

type service struct{}

// New returns a pointer to a new service object.
func New() *service {
	return &service{}
}

// Get returns the domain.Entry object with the given UUID.
func (srv *service) Get(id string) (domain.Entry, error) {
	return domain.Entry{}, nil
}
