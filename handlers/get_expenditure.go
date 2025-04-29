package handlers

import (
	"encoding/json"
	"go-expense-tracker/domain"
	"net/http"
	"strings"
)

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
