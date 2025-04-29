package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"go-expense-tracker/domain"
	"net/http"
	"strings"
	"time"
)

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
