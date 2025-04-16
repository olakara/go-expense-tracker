package domain

import (
	"errors"
	"github.com/google/uuid"
)

var ErrInvalidExpenditureAmount = errors.New("invalid expenditure amount")
var ErrExpenditureDescriptionEmpty = errors.New("expenditure description cannot be empty")

// Expenditure represents a money expenditure by a person
type Expenditure struct {
	ID          uuid.UUID `json:"id"`          // Unique identifier for the expenditure
	Description string    `json:"description"` // Description of what the money was spent on
	Amount      float64   `json:"amount"`      // Amount of money spent
}

func NewExpenditure(description string, amount float64) (*Expenditure, error) {

	if description == "" {
		return nil, ErrExpenditureDescriptionEmpty
	}

	if amount <= 0 {
		return nil, ErrInvalidExpenditureAmount
	}

	return &Expenditure{
		ID:          uuid.New(),
		Description: description,
		Amount:      amount,
	}, nil
}
