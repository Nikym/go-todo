package entrySrv

import (
	"errors"
	"github.com/Nikym/go-todo/internal/core/domain"
	"github.com/Nikym/go-todo/internal/core/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type TestEntry struct {
	Name     string
	Input    interface{}
	Expected interface{}
	Err      bool
}

func TestService_Get(t *testing.T) {
	tests := []TestEntry{
		{
			Name:  "should return an entry when given valid id",
			Input: "17beccd2-c5e8-4744-9b5f-98163b4a479d",
			Expected: &domain.Entry{
				ID:          "17beccd2-c5e8-4744-9b5f-98163b4a479d",
				Title:       "Test Title",
				Description: "Test Desc",
				Done:        false,
			},
			Err: false,
		},
		{
			Name:     "should return error when given invalid id",
			Input:    "invalid",
			Expected: &domain.Entry{},
			Err:      true,
		},
	}

	mockEntryRepository := &ports.MockEntryRepository{}
	mockEntryRepository.
		On("Get", "17beccd2-c5e8-4744-9b5f-98163b4a479d").
		Return(&domain.Entry{
			ID:          "17beccd2-c5e8-4744-9b5f-98163b4a479d",
			Title:       "Test Title",
			Description: "Test Desc",
			Done:        false,
		}, nil)
	mockEntryRepository.
		On("Get", "invalid").
		Return(&domain.Entry{}, errors.New("error get"))

	service := New(mockEntryRepository)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actual, err := service.Get(test.Input.(string))

			assert.EqualValues(t, test.Expected, actual)
			assert.Equal(t, test.Err, err != nil)
		})
	}
}

func TestService_Create(t *testing.T) {
	type input struct {
		title       string
		description string
	}
	tests := []TestEntry{
		{
			Name: "should create new entry and return domain.Entry object on valid input",
			Input: input{
				title:       "Test Title",
				description: "Test Description",
			},
			Expected: &domain.Entry{
				Title:       "Test Title",
				Description: "Test Description",
				Done:        false,
			},
			Err: false,
		},
		{
			Name: "should return error when repository cannot be accessed",
			Input: input{
				title:       "Error",
				description: "Error",
			},
			Expected: &domain.Entry{},
			Err:      true,
		},
	}

	mockEntryRepository := &ports.MockEntryRepository{}
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
		in := test.Input.(input)
		t.Run(test.Name, func(t *testing.T) {
			actual, err := service.Create(in.title, in.description)
			if err != nil {
				assert.True(t, test.Err)
			} else {
				expected := test.Expected.(*domain.Entry)

				assert.IsType(t, &domain.Entry{}, actual)
				assert.Equal(t, expected.Title, actual.Title)
				assert.Equal(t, expected.Description, actual.Description)
			}
		})
	}
}
