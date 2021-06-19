package entryHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nikym/go-todo/internal/core/domain"
	mocks "github.com/Nikym/go-todo/mocks/core/ports"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setUp() (*mocks.EntryService, *HTTPEntryHandler) {
	mockService := &mocks.EntryService{}
	httpEntryHandler := NewHTTPEntryHandler(mockService)
	return mockService, httpEntryHandler
}

func TestHTTPEntryHandler_Get(t *testing.T) {
	mockService, httpEntryHandler := setUp()
	mockService.
		On("Get", "1d126f09-4daf-447e-aaab-74765d8aefa2").
		Return(&domain.Entry{
			ID:          "1d126f09-4daf-447e-aaab-74765d8aefa2",
			Title:       "Test Title",
			Description: "Test Description",
			Done:        false,
		}, nil)
	mockService.
		On("Get", "invalid").
		Return(&domain.Entry{}, errors.New("invalid"))

	tests := []struct {
		name   string
		id     string
		status int
	}{
		{
			name:   "should return OK when given an entry ID that is present",
			id:     "1d126f09-4daf-447e-aaab-74765d8aefa2",
			status: http.StatusOK,
		},
		{
			name:   "should return Internal Server Error when given entry ID not present",
			id:     "invalid",
			status: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := fmt.Sprintf("/api/entry/%s", test.id)
			req := httptest.NewRequest("GET", path, nil)
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/api/entry/{id}", httpEntryHandler.Get)
			router.ServeHTTP(rr, req)

			assert.EqualValues(t, test.status, rr.Code)
		})
	}
}

func TestHTTPEntryHandler_Delete(t *testing.T) {
	mockService, httpEntryHandler := setUp()
	mockService.
		On("Delete", "1d126f09-4daf-447e-aaab-74765d8aefa2").
		Return(nil)
	mockService.
		On("Delete", "invalid").
		Return(errors.New("invalid"))

	tests := []struct {
		name    string
		id      string
		success bool
	}{
		{
			name:    "should return OK when given valid entry id",
			id:      "1d126f09-4daf-447e-aaab-74765d8aefa2",
			success: true,
		},
		{
			name:    "should not be successful when given invalid entry id",
			id:      "invalid",
			success: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := fmt.Sprintf("/api/entry/%s", test.id)
			req := httptest.NewRequest("DELETE", path, nil)
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/api/entry/{id}", httpEntryHandler.Delete)
			router.ServeHTTP(rr, req)

			assert.Equal(t, test.success, rr.Code == http.StatusOK)
		})
	}
}

func TestHTTPEntryHandler_Create(t *testing.T) {
	mockService, httpEntryHandler := setUp()
	mockService.
		On("Create", "Test Title", "Test Description").
		Return(&domain.Entry{
			ID:          "1d126f09-4daf-447e-aaab-74765d8aefa2",
			Title:       "Test Title",
			Description: "Test Description",
			Done:        false,
		}, nil)
	mockService.
		On("Create", "Invalid", "Invalid").
		Return(&domain.Entry{}, errors.New("invalid"))

	tests := []struct {
		name             string
		inputTitle       string
		inputDescription string
		success          bool
	}{
		{
			name:             "should return entry JSON when valid title and description given",
			inputTitle:       "Test Title",
			inputDescription: "Test Description",
			success:          true,
		},
		{
			name:             "should not be successful when invalid json body given",
			inputTitle:       "Invalid",
			inputDescription: "Invalid",
			success:          false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payload := fmt.Sprintf(`
				{
					"title": "%s",
					"description": "%s"
				}
			`, test.inputTitle, test.inputDescription)

			req := httptest.NewRequest(
				"POST",
				"/api/entry",
				strings.NewReader(payload),
			)
			rr := httptest.NewRecorder()
			httpEntryHandler.Create(rr, req)

			assert.Equal(t, test.success, rr.Code == http.StatusOK)
			if test.success {
				var returnedEntry domain.Entry
				if err := json.Unmarshal(rr.Body.Bytes(), &returnedEntry); err != nil {
					panic(err)
				}
				assert.EqualValues(t, test.inputTitle, returnedEntry.Title)
				assert.EqualValues(t, test.inputDescription, returnedEntry.Description)
			}
		})
	}
}

func TestHTTPEntryHandler_Update(t *testing.T) {
	mockService, httpEntryHandler := setUp()
	testEntry := &domain.Entry{
		ID:          "1d126f09-4daf-447e-aaab-74765d8aefa2",
		Title:       "Test Title",
		Description: "Test Description",
		Done:        false,
	}
	testUpdateEntry := &domain.Entry{
		ID:          "1d126f09-4daf-447e-aaab-74765d8aefa2",
		Title:       "Test Title 2",
		Description: "Test Description 2",
		Done:        false,
	}
	mockService.
		On("Get", "1d126f09-4daf-447e-aaab-74765d8aefa2").
		Return(testEntry, nil)
	mockService.
		On("Get", "invalid").
		Return(&domain.Entry{}, errors.New("invalid"))
	mockService.
		On("Update", "1d126f09-4daf-447e-aaab-74765d8aefa2", testUpdateEntry).
		Return(nil)

	tests := []struct {
		name       string
		inputId    string
		inputEntry string
		success    bool
	}{
		{
			name:    "should return OK when given valid entry ID and JSON",
			inputId: "1d126f09-4daf-447e-aaab-74765d8aefa2",
			inputEntry: `{
				"title": "Test Title 2",
				"description": "Test Description 2"
			}`,
			success: true,
		},
		{
			name:       "should return error when given invalid entry ID",
			inputId:    "invalid",
			inputEntry: `{"title": "Test Title 2"}`,
			success:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := fmt.Sprintf("/api/entry/%s", test.inputId)
			payload := strings.NewReader(test.inputEntry)
			req := httptest.NewRequest("PATCH", path, payload)
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/api/entry/{id}", httpEntryHandler.Update)
			router.ServeHTTP(rr, req)

			assert.EqualValues(t, test.success, rr.Code == http.StatusOK)
		})
	}
}
