package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/LegationPro/corpus-reader/internal/service/counter"
)

type Handler struct {
	logger     *slog.Logger
	maxWorkers int
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{
		logger:     logger,
		maxWorkers: 10,
	}
}

// Handle counter request
func (h *Handler) handleCounter(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming request
	var request CounterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// Initialize a new counter instance
	// The root directory should ALWAYS be set to corpus
	counterInstance := counter.New(request.Word, "corpus", h.maxWorkers)

	// Start counting and process errors
	for err := range counterInstance.Count() {
		if err != nil {
			// If directory error occurs, attempt to look for it
			if strings.Contains(err.Error(), "lstat") {
				h.logger.Info("Directory not found, attempting to locate it")

				// Look for the directory
				foundPath, err := counterInstance.LookForDirectory(request.Directory)
				if err != nil {
					h.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to look for directory: %v", err))
					return
				}

				// If the directory is found, update root path
				counterInstance.UpdateRoot(foundPath)

				// Retry counting after updating root
				if err := h.retryCount(counterInstance, w, request.Directory); err != nil {
					return
				}

				break
			}

			// Handle other errors (non-directory issues)
			h.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to count word: %v", err))
			return
		}
	}

	// Successfully finished counting, respond with results
	h.respondWithSuccess(w, request.Directory, counterInstance)
}

// Retry the counting after updating the root directory
func (h *Handler) retryCount(counterInstance counter.ICounter, w http.ResponseWriter, directory string) error {
	// Start counting again
	for err := range counterInstance.Count() {
		if err != nil {
			h.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to count word after directory update: %v", err))
			return err
		}
	}

	// Return successful response after retry
	h.respondWithSuccess(w, directory, counterInstance)
	return nil
}

// Send error response with logging
func (h *Handler) respondWithError(w http.ResponseWriter, statusCode int, errorMessage string) {
	h.logger.Error(errorMessage)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errorMessage})
}

// Send success response with the count result
func (h *Handler) respondWithSuccess(w http.ResponseWriter, directory string, counterInstance counter.ICounter) {
	h.logger.Info(fmt.Sprintf("Successfully counted word occurrences in directory: %s", directory))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CounterResponse{
		Count:     int(counterInstance.GetCount()),
		Directory: directory,
	})
}
