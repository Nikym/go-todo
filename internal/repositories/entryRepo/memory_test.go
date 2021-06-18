package entryRepo

import (
	"encoding/json"
	"github.com/Nikym/go-todo/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

var repo = NewMemKVS()

func setUp() {
	exampleEntry := domain.Entry{
		ID:          "5b2c9d9f-7bb2-401d-b0e1-1e5d8ea955ca",
		Title:       "Test Title",
		Description: "Test Description",
		Done:        false,
	}
	bytes, err := json.Marshal(exampleEntry)
	if err != nil {
		panic(err)
	}

	repo.kvs[exampleEntry.ID] = bytes
}

func tearDown() {
	repo.kvs = map[string][]byte{}
}

func TestMemKVS_Get(t *testing.T) {
	setUp()
	defer tearDown()

	tests := []struct {
		name     string
		key      string
		expected *domain.Entry
		err      bool
	}{
		{
			name: "should return entry when id of stored entry given",
			key:  "5b2c9d9f-7bb2-401d-b0e1-1e5d8ea955ca",
			expected: &domain.Entry{
				ID:          "5b2c9d9f-7bb2-401d-b0e1-1e5d8ea955ca",
				Title:       "Test Title",
				Description: "Test Description",
				Done:        false,
			},
			err: false,
		},
		{
			name:     "should return error when id not found in repository",
			key:      "invalid",
			expected: &domain.Entry{},
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := repo.Get(test.key)

			assert.EqualValues(t, test.expected, actual)
			assert.EqualValues(t, test.err, err != nil)
		})
	}
}

func TestMemKVS_Save(t *testing.T) {
	setUp()
	defer tearDown()

	tests := []struct {
		name    string
		input   *domain.Entry
		present bool
		err     bool
	}{
		{
			name: "should save an entry when a valid entry given",
			input: &domain.Entry{
				ID:          "valid",
				Title:       "Test Title 2",
				Description: "Test Description 2",
				Done:        false,
			},
			present: true,
			err:     false,
		},
		{
			name: "should return error when empty string given for id",
			input: &domain.Entry{
				ID:          "",
				Title:       "Test Title 2",
				Description: "Test Description 2",
				Done:        false,
			},
			present: false,
			err:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.Save(test.input)

			_, present := repo.kvs[test.input.ID]
			assert.EqualValues(t, test.present, present)
			assert.EqualValues(t, test.err, err != nil)
		})
	}
}

func TestMemKVS_Delete(t *testing.T) {
	setUp()
	defer tearDown()

	tests := []struct {
		name  string
		input string
		err   bool
	}{
		{
			name:  "should remove entry when id of a stored entry given",
			input: "5b2c9d9f-7bb2-401d-b0e1-1e5d8ea955ca",
			err:   false,
		},
		{
			name:  "should return no error when id of not stored entry given",
			input: "test",
			err:   false,
		},
		{
			name:  "should return error when empty string given for id",
			input: "",
			err:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.Delete(test.input)

			_, present := repo.kvs[test.input]
			assert.False(t, present)
			assert.EqualValues(t, test.err, err != nil)
		})
	}
}

func TestMemKVS_Update(t *testing.T) {
	setUp()
	defer tearDown()

	tests := []struct {
		name       string
		inputId    string
		inputEntry *domain.Entry
		present    bool
		err        bool
	}{
		{
			name:    "should change the stored entry to new details when given entry id present in repository",
			inputId: "5b2c9d9f-7bb2-401d-b0e1-1e5d8ea955ca",
			inputEntry: &domain.Entry{
				ID:          "5b2c9d9f-7bb2-401d-b0e1-1e5d8ea955ca",
				Title:       "Updated Title",
				Description: "Updated Description",
				Done:        false,
			},
			present: true,
			err:     false,
		},
		{
			name:    "should return error when no entry present with given id",
			inputId: "invalid",
			inputEntry: &domain.Entry{
				ID:          "invalid",
				Title:       "Updated Title",
				Description: "Updated Description",
				Done:        false,
			},
			present: false,
			err:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.Update(test.inputId, test.inputEntry)

			stored, present := repo.kvs[test.inputId]

			if present {
				bytes, err := json.Marshal(test.inputEntry)
				if err != nil {
					panic(err)
				}
				assert.EqualValues(t, bytes, stored)
			}
			assert.EqualValues(t, test.present, present)
			assert.EqualValues(t, test.err, err != nil)
		})
	}
}
