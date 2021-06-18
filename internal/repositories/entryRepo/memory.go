package entryRepo

import (
	"encoding/json"
	"errors"
	"github.com/Nikym/go-todo/internal/core/domain"
)

type memKVS struct {
	kvs map[string][]byte
}

// NewMemKVS returns a pointer to an in-memory entry repository.
func NewMemKVS() *memKVS {
	return &memKVS{
		kvs: map[string][]byte{},
	}
}

// Get retrieves an entry with a specified ID from the in-memory KVS repository.
func (r *memKVS) Get(id string) (*domain.Entry, error) {
	if val, ok := r.kvs[id]; ok {
		entry := domain.Entry{}
		if err := json.Unmarshal(val, &entry); err != nil {
			return &domain.Entry{}, err
		}
		return &entry, nil
	}

	return &domain.Entry{}, errors.New("entry not found in repository")
}

// Save stores a given domain.Entry object in the in-memory KVS repository.
func (r *memKVS) Save(entry *domain.Entry) error {
	if entry.ID != "" {
		bytes, err := json.Marshal(*entry)
		if err != nil {
			return err
		}
		r.kvs[entry.ID] = bytes
		return nil
	}

	return errors.New("id cannot be an empty string")
}

// Delete removes a domain.Entry object with a given ID from the in-memory KVS repository.
func (r *memKVS) Delete(id string) error {
	if id != "" {
		delete(r.kvs, id)
		return nil
	}
	return errors.New("id cannot be an empty string")
}

// Update sets the entry stored in KVS repository with given ID to the domain.Entry specified.
func (r *memKVS) Update(id string, entry *domain.Entry) error {
	if _, ok := r.kvs[id]; ok {
		bytes, err := json.Marshal(*entry)
		if err != nil {
			return err
		}

		r.kvs[id] = bytes
		return nil
	}

	return errors.New("no entry with given id found in repository")
}
