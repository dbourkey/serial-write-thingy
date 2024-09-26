package status

import (
	"encoding/json"
	"io"
	"net/http"
)

type responseBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   error  `json:"error,omitempty"`
}

type StatusHandler struct {
	Serialiser *Serialiser
}

func NewStatusHandler(serialiser *Serialiser) *StatusHandler {
	return &StatusHandler{Serialiser: serialiser}
}

func respondError(w http.ResponseWriter, statusCode int, msg string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&responseBody{Status: statusCode, Message: msg, Error: err})
}

func respondOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&responseBody{Status: http.StatusOK, Message: "ok"})
}

func (uh StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to read request body", err)
		return
	}

	var containerStatus ContainerReport
	if err := json.Unmarshal(b, &containerStatus); err != nil {
		respondError(w, http.StatusBadRequest, "Failed to decode request body", err)
		return
	}

	if err := uh.Serialiser.Write(r.Context(), containerStatus); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to queue update", err)
		return
	}

	respondOK(w)
}
