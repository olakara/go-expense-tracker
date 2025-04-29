package handlers

import (
	"go-expense-tracker/domain"
	"net/http"
	"strings"
)

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
