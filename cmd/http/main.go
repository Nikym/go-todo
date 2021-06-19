package main

import (
	"github.com/Nikym/go-todo/internal/core/services/entrySrv"
	"github.com/Nikym/go-todo/internal/handlers/entryHandler"
	"github.com/Nikym/go-todo/internal/repositories/entryRepo"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func SetupRoutes(router *mux.Router, httpHandler *entryHandler.HTTPEntryHandler) {
	log.Println("Setting up routes...")
	router.HandleFunc("/api/entry/{id}", httpHandler.Get).Methods("GET")
	router.HandleFunc("/api/entry/{id}", httpHandler.Delete).Methods("DELETE")
	router.HandleFunc("/api/entry/{id}", httpHandler.Update).Methods("PATCH")
	router.HandleFunc("/api/entry", httpHandler.Create).Methods("POST")
}

func main() {
	log.Println("Started HTTP server")

	entryRepository := entryRepo.NewMemKVS()
	entryService := entrySrv.New(entryRepository)
	httpHandler := entryHandler.NewHTTPEntryHandler(entryService)

	router := mux.NewRouter()
	SetupRoutes(router, httpHandler)

	log.Println("Finished setup")
	log.Fatal(http.ListenAndServe(":8080", router))
}
