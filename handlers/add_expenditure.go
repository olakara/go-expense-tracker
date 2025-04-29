package handlers

import (
	"encoding/json"
	"go-expense-tracker/domain"
	"net/http"
)

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

	//TODO: Need to check if category exists

	expenditure, err := domain.NewExpenditure(req.Description, req.Amount, req.Date, req.CategoryId)

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
