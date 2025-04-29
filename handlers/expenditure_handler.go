package handlers

import (
	"go-expense-tracker/domain"
	"log/slog"
	"net/http"
	"strings"
)

type ExpenditureHandler struct {
	service domain.ExpenditureRepository
	logger  *slog.Logger
}

func NewExpenditureHandler(service domain.ExpenditureRepository, logger *slog.Logger) *ExpenditureHandler {
	return &ExpenditureHandler{
		service: service,
		logger:  logger,
	}
}

func ExpenditureRouter(handler *ExpenditureHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/expenditures" {
			switch r.Method {
			case http.MethodGet:
				handler.GetAllExpenditures(w, r)
			case http.MethodPost:
				handler.AddExpenditure(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if strings.HasPrefix(path, "/expenditures/") {
			switch r.Method {
			case http.MethodGet:
				handler.GetExpenditureByID(w, r)
			case http.MethodPut:
				handler.UpdateExpenditure(w, r)
			case http.MethodDelete:
				handler.DeleteExpenditure(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		http.NotFound(w, r)
	})
}
