package services

import (
	"database/sql"
	"errors"
	"fmt"
	"go-expense-tracker/domain"
	"log/slog"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// DBService implements the ExpenditureRepository interface using PostgreSQL
type DBService struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewDBService creates a new DBService with the given connection parameters
func NewDBService(host string, port int, user, password, dbname string, logger *slog.Logger) (*DBService, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create the expenditures table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS expenditures (
			id UUID PRIMARY KEY,
			description TEXT NOT NULL,
			amount DECIMAL(10, 2) NOT NULL,
			date TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create expenditures table: %w", err)
	}

	return &DBService{
		db:     db,
		logger: logger,
	}, nil
}

// Close closes the database connection
func (s *DBService) Close() error {
	return s.db.Close()
}

// AddExpenditure adds a new expenditure to the database
func (s *DBService) AddExpenditure(expenditure *domain.Expenditure) error {
	s.logger.Debug("Adding expenditure to database", 
		"id", expenditure.ID, 
		"description", expenditure.Description, 
		"amount", expenditure.Amount, 
		"date", expenditure.Date)

	// Check if expenditure with this ID already exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM expenditures WHERE id = $1)", expenditure.ID).Scan(&exists)
	if err != nil {
		s.logger.Error("Error checking if expenditure exists", "error", err, "id", expenditure.ID)
		return fmt.Errorf("error checking if expenditure exists: %w", err)
	}

	if exists {
		s.logger.Warn("Expenditure already exists", "id", expenditure.ID)
		return domain.ErrExpenditureAlreadyExists
	}

	// Insert the expenditure
	_, err = s.db.Exec(
		"INSERT INTO expenditures (id, description, amount, date) VALUES ($1, $2, $3, $4)",
		expenditure.ID, expenditure.Description, expenditure.Amount, expenditure.Date,
	)
	if err != nil {
		s.logger.Error("Error inserting expenditure", "error", err, "id", expenditure.ID)
		return fmt.Errorf("error inserting expenditure: %w", err)
	}

	s.logger.Info("Expenditure added successfully", "id", expenditure.ID)
	return nil
}

// GetExpenditureByID retrieves an expenditure by its ID
func (s *DBService) GetExpenditureByID(id string) (*domain.Expenditure, error) {
	s.logger.Debug("Getting expenditure by ID", "id", id)

	// Parse the ID string to UUID
	expenditureID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error("Invalid UUID format", "error", err, "id", id)
		return nil, fmt.Errorf("invalid UUID format: %w", err)
	}

	// Query the expenditure
	var expenditure domain.Expenditure
	err = s.db.QueryRow(
		"SELECT id, description, amount, date FROM expenditures WHERE id = $1",
		expenditureID,
	).Scan(&expenditure.ID, &expenditure.Description, &expenditure.Amount, &expenditure.Date)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("Expenditure not found", "id", id)
			return nil, domain.ErrExpenditureNotFound
		}
		s.logger.Error("Error querying expenditure", "error", err, "id", id)
		return nil, fmt.Errorf("error querying expenditure: %w", err)
	}

	s.logger.Debug("Found expenditure", 
		"id", id, 
		"description", expenditure.Description, 
		"amount", expenditure.Amount, 
		"date", expenditure.Date)
	return &expenditure, nil
}

// GetAllExpenditures retrieves all expenditures from the database
func (s *DBService) GetAllExpenditures() ([]*domain.Expenditure, error) {
	s.logger.Debug("Getting all expenditures")

	// Query all expenditures
	rows, err := s.db.Query("SELECT id, description, amount, date FROM expenditures")
	if err != nil {
		s.logger.Error("Error querying all expenditures", "error", err)
		return nil, fmt.Errorf("error querying all expenditures: %w", err)
	}
	defer rows.Close()

	// Collect all expenditures
	var expenditures []*domain.Expenditure
	for rows.Next() {
		var expenditure domain.Expenditure
		err := rows.Scan(&expenditure.ID, &expenditure.Description, &expenditure.Amount, &expenditure.Date)
		if err != nil {
			s.logger.Error("Error scanning expenditure row", "error", err)
			return nil, fmt.Errorf("error scanning expenditure row: %w", err)
		}
		expenditures = append(expenditures, &expenditure)
	}

	if err = rows.Err(); err != nil {
		s.logger.Error("Error iterating expenditure rows", "error", err)
		return nil, fmt.Errorf("error iterating expenditure rows: %w", err)
	}

	s.logger.Info("Retrieved all expenditures", "count", len(expenditures))
	return expenditures, nil
}

// UpdateExpenditure updates an existing expenditure
func (s *DBService) UpdateExpenditure(expenditure *domain.Expenditure) error {
	s.logger.Debug("Updating expenditure", 
		"id", expenditure.ID, 
		"description", expenditure.Description, 
		"amount", expenditure.Amount, 
		"date", expenditure.Date)

	// Check if expenditure exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM expenditures WHERE id = $1)", expenditure.ID).Scan(&exists)
	if err != nil {
		s.logger.Error("Error checking if expenditure exists", "error", err, "id", expenditure.ID)
		return fmt.Errorf("error checking if expenditure exists: %w", err)
	}

	if !exists {
		s.logger.Warn("Expenditure not found for update", "id", expenditure.ID)
		return domain.ErrExpenditureNotFound
	}

	// Update the expenditure
	_, err = s.db.Exec(
		"UPDATE expenditures SET description = $1, amount = $2, date = $3 WHERE id = $4",
		expenditure.Description, expenditure.Amount, expenditure.Date, expenditure.ID,
	)
	if err != nil {
		s.logger.Error("Error updating expenditure", "error", err, "id", expenditure.ID)
		return fmt.Errorf("error updating expenditure: %w", err)
	}

	s.logger.Info("Expenditure updated successfully", "id", expenditure.ID)
	return nil
}

// DeleteExpenditure deletes an expenditure by its ID
func (s *DBService) DeleteExpenditure(id string) error {
	s.logger.Debug("Deleting expenditure", "id", id)

	// Parse the ID string to UUID
	expenditureID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error("Invalid UUID format", "error", err, "id", id)
		return fmt.Errorf("invalid UUID format: %w", err)
	}

	// Check if expenditure exists
	var exists bool
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM expenditures WHERE id = $1)", expenditureID).Scan(&exists)
	if err != nil {
		s.logger.Error("Error checking if expenditure exists", "error", err, "id", id)
		return fmt.Errorf("error checking if expenditure exists: %w", err)
	}

	if !exists {
		s.logger.Warn("Expenditure not found for deletion", "id", id)
		return domain.ErrExpenditureNotFound
	}

	// Delete the expenditure
	_, err = s.db.Exec("DELETE FROM expenditures WHERE id = $1", expenditureID)
	if err != nil {
		s.logger.Error("Error deleting expenditure", "error", err, "id", id)
		return fmt.Errorf("error deleting expenditure: %w", err)
	}

	s.logger.Info("Expenditure deleted successfully", "id", id)
	return nil
}
