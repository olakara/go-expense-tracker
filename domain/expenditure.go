package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var ErrInvalidExpenditureAmount = errors.New("invalid expenditure amount")
var ErrExpenditureDescriptionEmpty = errors.New("expenditure description cannot be empty")
var ErrExpenditureFutureDate = errors.New("expenditure date cannot be in the future")

// Expenditure represents a money expenditure by a person
type Expenditure struct {
	ID          uuid.UUID `json:"id"`          // Unique identifier for the expenditure
	Description string    `json:"description"` // Description of what the money was spent on
	Amount      float64   `json:"amount"`      // Amount of money spent
	Date        time.Time `json:"date"`        // Date when the expenditure occurred
}

func NewExpenditure(description string, amount float64, date time.Time) (*Expenditure, error) {

	if description == "" {
		return nil, ErrExpenditureDescriptionEmpty
	}

	if amount <= 0 {
		return nil, ErrInvalidExpenditureAmount
	}

	// Check if the date is in the future
	if date.After(time.Now()) {
		return nil, ErrExpenditureFutureDate
	}

	return &Expenditure{
		ID:          uuid.New(),
		Description: description,
		Amount:      amount,
		Date:        date,
	}, nil
}
