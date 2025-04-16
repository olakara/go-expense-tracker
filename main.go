package main

import (
	"encoding/json"
	"fmt"
	"go-expense-tracker/domain"
	"go-expense-tracker/services"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ExpenditureHandler struct {
	service *services.MemoryService
	logger  *slog.Logger
}

type ExpenditureRequest struct {
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
}

func NewExpenditureHandler(service *services.MemoryService, logger *slog.Logger) *ExpenditureHandler {
	return &ExpenditureHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ExpenditureHandler) AddExpenditure(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling add expenditure request", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)

	if r.Method != http.MethodPost {
		h.logger.Warn("Method not allowed", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ExpenditureRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Decoded expenditure request", "description", req.Description, "amount", req.Amount, "date", req.Date)

	expenditure, err := domain.NewExpenditure(req.Description, req.Amount, req.Date)
	if err != nil {
		h.logger.Error("Failed to create expenditure", "error", err, "description", req.Description, "amount", req.Amount, "date", req.Date)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.AddExpenditure(expenditure)
	if err != nil {
		h.logger.Error("Failed to add expenditure", "error", err, "id", expenditure.ID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully added expenditure", "id", expenditure.ID, "description", expenditure.Description, "date", expenditure.Date)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expenditure)
}

func (h *ExpenditureHandler) GetAllExpenditures(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling get all expenditures request", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)

	if r.Method != http.MethodGet {
		h.logger.Warn("Method not allowed", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	expenditures, err := h.service.GetAllExpenditures()
	if err != nil {
		h.logger.Error("Failed to get all expenditures", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully retrieved all expenditures", "count", len(expenditures))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenditures)
}

func (h *ExpenditureHandler) GetExpenditureByID(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling get expenditure by ID request", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)

	if r.Method != http.MethodGet {
		h.logger.Warn("Method not allowed", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/expenditures/")
	h.logger.Debug("Getting expenditure by ID", "id", id)

	expenditure, err := h.service.GetExpenditureByID(id)
	if err != nil {
		if err == domain.ErrExpenditureNotFound {
			h.logger.Warn("Expenditure not found", "id", id)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to get expenditure by ID", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully retrieved expenditure", "id", id, "description", expenditure.Description, "date", expenditure.Date)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenditure)
}

func (h *ExpenditureHandler) UpdateExpenditure(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling update expenditure request", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)

	if r.Method != http.MethodPut {
		h.logger.Warn("Method not allowed", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/expenditures/")
	h.logger.Debug("Updating expenditure", "id", id)

	_, err := h.service.GetExpenditureByID(id)
	if err != nil {
		if err == domain.ErrExpenditureNotFound {
			h.logger.Warn("Expenditure not found for update", "id", id)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to check expenditure existence", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req ExpenditureRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("Failed to decode update request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Decoded update request", "id", id, "description", req.Description, "amount", req.Amount, "date", req.Date)

	if req.Description == "" {
		h.logger.Warn("Empty description in update request", "id", id)
		http.Error(w, domain.ErrExpenditureDescriptionEmpty.Error(), http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		h.logger.Warn("Invalid amount in update request", "id", id, "amount", req.Amount)
		http.Error(w, domain.ErrInvalidExpenditureAmount.Error(), http.StatusBadRequest)
		return
	}

	// Check if the date is in the future
	if req.Date.After(time.Now()) {
		h.logger.Warn("Future date in update request", "id", id, "date", req.Date)
		http.Error(w, domain.ErrExpenditureFutureDate.Error(), http.StatusBadRequest)
		return
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Error("Failed to parse UUID", "id", id, "error", err)
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	expenditure := &domain.Expenditure{
		ID:          parsedUUID,
		Description: req.Description,
		Amount:      req.Amount,
		Date:        req.Date,
	}

	err = h.service.UpdateExpenditure(expenditure)
	if err != nil {
		h.logger.Error("Failed to update expenditure", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully updated expenditure", "id", id, "description", expenditure.Description, "date", expenditure.Date)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenditure)
}

func (h *ExpenditureHandler) DeleteExpenditure(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling delete expenditure request", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)

	if r.Method != http.MethodDelete {
		h.logger.Warn("Method not allowed", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/expenditures/")
	h.logger.Debug("Deleting expenditure", "id", id)

	err := h.service.DeleteExpenditure(id)
	if err != nil {
		if err == domain.ErrExpenditureNotFound {
			h.logger.Warn("Expenditure not found for deletion", "id", id)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to delete expenditure", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully deleted expenditure", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

func expenditureRouter(handler *ExpenditureHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/expenditures" {
			switch r.Method {
			case http.MethodGet:
				handler.GetAllExpenditures(w, r)
			case http.MethodPost:
				handler.AddExpenditure(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if strings.HasPrefix(path, "/expenditures/") {
			switch r.Method {
			case http.MethodGet:
				handler.GetExpenditureByID(w, r)
			case http.MethodPut:
				handler.UpdateExpenditure(w, r)
			case http.MethodDelete:
				handler.DeleteExpenditure(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		http.NotFound(w, r)
	})
}

// LoggingMiddleware adds request logging to all HTTP requests
func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response wrapper to capture the status code
		wrapped := NewResponseWriter(w)

		// Process the request
		next.ServeHTTP(wrapped, r)

		// Log the request details
		duration := time.Since(start)
		logger.Info("HTTP request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
	})
}

// ResponseWriter wraps http.ResponseWriter to capture the status code
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

// WriteHeader captures the status code and passes it to the wrapped ResponseWriter
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func main() {
	port := 8080

	// Configure structured logger
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	logger.Info("Starting expense tracker application")

	// Initialize the service with logger
	service := services.NewMemoryService(logger)
	handler := NewExpenditureHandler(service, logger)

	// Set up the routes
	router := expenditureRouter(handler)

	// Apply logging middleware
	loggedRouter := LoggingMiddleware(logger, router)

	http.Handle("/expenditures", loggedRouter)
	http.Handle("/expenditures/", loggedRouter)

	// Start the server
	serverAddr := fmt.Sprintf(":%d", port)
	logger.Info("Starting HTTP server", "address", serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
