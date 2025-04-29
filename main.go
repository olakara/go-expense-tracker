package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"go-expense-tracker/domain"
	"go-expense-tracker/handlers"
	"go-expense-tracker/services"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

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

	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Parse command line flags
	useDB := flag.Bool("db", false, "Use PostgreSQL database instead of in-memory storage")
	flag.Parse()

	// Configure structured logger
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	logger.Info("Starting expense tracker application")

	// Initialize the appropriate service
	var service domain.ExpenditureRepository
	var err error

	if *useDB {
		// Get database parameters from environment variables
		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			dbHost = "localhost" // Default value
		}

		dbPortStr := os.Getenv("DB_PORT")
		dbPort := 5432 // Default value
		if dbPortStr != "" {
			var err error
			dbPort, err = strconv.Atoi(dbPortStr)
			if err != nil {
				logger.Error("Invalid DB_PORT value", "error", err, "value", dbPortStr)
				os.Exit(1)
			}
		}

		dbUser := os.Getenv("DB_USER")
		if dbUser == "" {
			dbUser = "postgres" // Default value
		}

		dbPassword := os.Getenv("DB_PASSWORD")
		if dbPassword == "" {
			dbPassword = "postgres" // Default value
		}

		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "expense_tracker" // Default value
		}

		logger.Info("Using PostgreSQL database for storage",
			"host", dbHost,
			"port", dbPort,
			"user", dbUser,
			"database", dbName)

		dbService, err := services.NewDBService(dbHost, dbPort, dbUser, dbPassword, dbName, logger)
		if err != nil {
			logger.Error("Failed to initialize database service", "error", err)
			os.Exit(1)
		}
		defer dbService.Close()

		service = dbService
	} else {
		logger.Info("Using in-memory storage")
		service = services.NewMemoryService(logger)
	}

	handler := handlers.NewExpenditureHandler(service, logger)

	// Set up the routes
	router := handlers.ExpenditureRouter(handler)

	// Apply logging middleware
	loggedRouter := LoggingMiddleware(logger, router)

	http.Handle("/expenditures", loggedRouter)
	http.Handle("/expenditures/", loggedRouter)

	// Start the server
	serverAddr := fmt.Sprintf(":%d", port)
	logger.Info("Starting HTTP server", "address", serverAddr)
	err = http.ListenAndServe(serverAddr, nil)
	if err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
