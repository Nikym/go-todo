package entrySrv

import (
	"errors"
	"github.com/Nikym/go-todo/internal/core/domain"
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

type MockEntryRepository struct {
	mock.Mock
}

func (m *MockEntryRepository) Get(id string) (*domain.Entry, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Entry), args.Error(1)
}

func (m *MockEntryRepository) Save(entry *domain.Entry) error {
	args := m.Called(entry)
	return args.Error(0)
}

func (m *MockEntryRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
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

	mockEntryRepository := &MockEntryRepository{}
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
		Return(&domain.Entry{}, errors.New("error"))

	service := New(mockEntryRepository)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actual, err := service.Get(test.Input.(string))

			assert.EqualValues(t, test.Expected, actual)
			assert.Equal(t, test.Err, err != nil)
		})
	}
}
