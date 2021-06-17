package entrySrv

import (
	"errors"
	"github.com/Nikym/go-todo/mocks/pkg/core/ports"
	"github.com/Nikym/go-todo/pkg/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_Get(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedTitle       string
		expectedDescription string
		err                 bool
	}{
		{
			name:                "should return an entry when given valid id",
			input:               "17beccd2-c5e8-4744-9b5f-98163b4a479d",
			expectedTitle:       "Test Title",
			expectedDescription: "Test Description",
			err:                 false,
		},
		{
			name:  "should return error when given invalid id",
			input: "invalid",
			err:   true,
		},
	}

	mockEntryRepository := &mocks.EntryRepository{}
	mockEntryRepository.
		On("Get", "17beccd2-c5e8-4744-9b5f-98163b4a479d").
		Return(&domain.Entry{
			ID:          "17beccd2-c5e8-4744-9b5f-98163b4a479d",
			Title:       "Test Title",
			Description: "Test Description",
			Done:        false,
		}, nil)
	mockEntryRepository.
		On("Get", "invalid").
		Return(&domain.Entry{}, errors.New("error get"))

	service := New(mockEntryRepository)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := service.Get(test.input)
			if err != nil {
				assert.True(t, test.err)
			} else {
				assert.IsType(t, &domain.Entry{}, actual)
				assert.Equal(t, test.expectedTitle, actual.Title)
				assert.Equal(t, test.expectedDescription, actual.Description)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name                string
		inputTitle          string
		inputDescription    string
		expectedTitle       string
		expectedDescription string
		err                 bool
	}{
		{
			name:                "should create new entry and return domain.Entry object on valid input",
			inputTitle:          "Test Title",
			inputDescription:    "Test Description",
			expectedTitle:       "Test Title",
			expectedDescription: "Test Description",
			err:                 false,
		},
		{
			name:             "should return error when repository cannot be accessed",
			inputTitle:       "Error",
			inputDescription: "Error",
			err:              true,
		},
		{
			name:       "should return error when title is less than 3 chars long",
			inputTitle: "te",
			err:        true,
		},
	}

	mockEntryRepository := &mocks.EntryRepository{}
	mockEntryRepository.
		On("Save", mock.MatchedBy(
			func(e *domain.Entry) bool { return e.Title == "Test Title" && e.Description == "Test Description" }),
		).
		Return(nil)
	mockEntryRepository.
		On("Save", mock.MatchedBy(
			func(e *domain.Entry) bool { return e.Title == "Error" && e.Description == "Error" }),
		).
		Return(errors.New("error create"))

	service := New(mockEntryRepository)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := service.Create(test.inputTitle, test.inputDescription)
			if err != nil {
				assert.True(t, test.err)
			} else {
				assert.IsType(t, &domain.Entry{}, actual)
				assert.Equal(t, test.expectedTitle, actual.Title)
				assert.Equal(t, test.expectedDescription, actual.Description)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   bool
	}{
		{
			name:  "should return no error when valid id is given",
			input: "154b07a0-76bd-4f85-83a5-5090cbf46552",
			err:   false,
		},
		{
			name:  "should return error when invalid id is given",
			input: "invalid",
			err:   true,
		},
	}

	mockEntryRepository := &mocks.EntryRepository{}
	mockEntryRepository.
		On("Delete", "154b07a0-76bd-4f85-83a5-5090cbf46552").
		Return(nil)
	mockEntryRepository.
		On("Delete", "invalid").
		Return(errors.New("error"))

	service := New(mockEntryRepository)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.Delete(test.input)

			assert.Equal(t, test.err, err != nil)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name       string
		inputId    string
		inputEntry *domain.Entry
		err        bool
	}{
		{
			name:    "should return no error when id is valid",
			inputId: "154b07a0-76bd-4f85-83a5-5090cbf46552",
			inputEntry: &domain.Entry{
				ID:          "154b07a0-76bd-4f85-83a5-5090cbf46552",
				Title:       "Test Title",
				Description: "Test Description",
				Done:        false,
			},
			err: false,
		},
		{
			name:    "should return error when id is invalid",
			inputId: "invalid",
			inputEntry: &domain.Entry{
				ID:          "invalid",
				Title:       "Test Title",
				Description: "Test Description",
				Done:        false,
			},
			err: true,
		},
	}

	mockEntryRepository := &mocks.EntryRepository{}
	mockEntryRepository.
		On("Update", "154b07a0-76bd-4f85-83a5-5090cbf46552", mock.MatchedBy(
			func(e *domain.Entry) bool { return true },
		)).
		Return(nil)
	mockEntryRepository.
		On("Update", "invalid", mock.MatchedBy(
			func(e *domain.Entry) bool { return true },
		)).
		Return(errors.New("error"))

	service := New(mockEntryRepository)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.Update(test.inputId, test.inputEntry)
			assert.Equal(t, test.err, err != nil)
		})
	}
}
