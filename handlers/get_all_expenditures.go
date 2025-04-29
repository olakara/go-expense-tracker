package handlers

import (
	"encoding/json"
	"net/http"
)

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
