package handlers

import (
	"github.com/google/uuid"
	"time"
)

type ExpenditureRequest struct {
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	CategoryId  uuid.UUID `json:"categoryId"`
}
