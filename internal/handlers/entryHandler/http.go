package entryHandler

import (
	"encoding/json"
	"github.com/Nikym/go-todo/internal/core/ports"
	"github.com/gorilla/mux"
	"net/http"
)

type response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type createJSON struct {
	Title       string
	Description string
}

type HTTPEntryHandler struct {
	EntryService ports.EntryService
}

// NewHTTPEntryHandler returns a pointer to the HTTP adapter for the ports.EntryService interface.
func NewHTTPEntryHandler(entryService ports.EntryService) *HTTPEntryHandler {
	return &HTTPEntryHandler{
		EntryService: entryService,
	}
}

// Get handles retrieval of a to-do entry through HTTP with a specified UUID within the URL.
func (h *HTTPEntryHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	id := vars["id"]

	entry, err := h.EntryService.Get(id)
	if err != nil {
		sendErrorResponse(w, "failed to retrieve entry with given ID", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		panic(err)
	}
}

// Create handles the creation of a new to-do entry through HTTP with given Title and Description within body.
func (h *HTTPEntryHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var details *createJSON
	if err := json.NewDecoder(r.Body).Decode(&details); err != nil {
		sendErrorResponse(w, "failed to decode json body", err)
		return
	}

	newEntry, err := h.EntryService.Create(details.Title, details.Description)
	if err != nil {
		sendErrorResponse(w, "failed to create to-do entry", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(*newEntry); err != nil {
		panic(err)
	}
}

// Delete removes an entry with a given ID through HTTP.
func (h *HTTPEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	id := vars["id"]

	err := h.EntryService.Delete(id)
	if err != nil {
		sendErrorResponse(w, "failed to delete entry with given id", err)
	}

	w.WriteHeader(http.StatusOK)
}

// Update updates the entry specified by the ID with the new values given in the body.
func (h *HTTPEntryHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	id := vars["id"]

	entry, err := h.EntryService.Get(id)
	if err != nil {
		sendErrorResponse(w, "failed to find entry with given id", err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(entry); err != nil {
		sendErrorResponse(w, "failed to decode json body", err)
		return
	}

	entry.ID = id
	if err := h.EntryService.Update(id, entry); err != nil {
		sendErrorResponse(w, "failed to update entry", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(
		response{Message: message, Error: err.Error()},
	); err != nil {
		panic(err)
	}
}
