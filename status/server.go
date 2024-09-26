package status

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Initialise a basic (not very secure) http server with hardcoded params.
// This is not meant to represent good code.
func NewServer(statusHandler *StatusHandler) *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	router.PathPrefix("/api/update").Methods("POST").Handler(statusHandler)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv
}
