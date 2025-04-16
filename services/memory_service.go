package services

import (
	domain "go-expense-tracker/domain"
	"log/slog"
	"sync"
)

type MemoryService struct {
	Expenditures map[string]*domain.Expenditure
	logger       *slog.Logger
	sync.RWMutex
}

func NewMemoryService(logger *slog.Logger) *MemoryService {
	return &MemoryService{
		Expenditures: make(map[string]*domain.Expenditure),
		logger:       logger,
	}
}

func (m *MemoryService) AddExpenditure(expenditure *domain.Expenditure) error {
	m.logger.Debug("Adding expenditure", "id", expenditure.ID, "description", expenditure.Description, "amount", expenditure.Amount)

	m.Lock()
	defer m.Unlock()

	if _, exists := m.Expenditures[expenditure.ID.String()]; exists {
		m.logger.Warn("Expenditure already exists", "id", expenditure.ID)
		return domain.ErrExpenditureAlreadyExists
	}

	m.Expenditures[expenditure.ID.String()] = expenditure
	m.logger.Info("Expenditure added successfully", "id", expenditure.ID, "total_count", len(m.Expenditures))
	return nil
}

func (m *MemoryService) GetExpenditureByID(id string) (*domain.Expenditure, error) {
	m.logger.Debug("Getting expenditure by ID", "id", id)

	m.RLock()
	defer m.RUnlock()

	expenditure, exists := m.Expenditures[id]
	if !exists {
		m.logger.Warn("Expenditure not found", "id", id)
		return nil, domain.ErrExpenditureNotFound
	}

	m.logger.Debug("Found expenditure", "id", id, "description", expenditure.Description, "amount", expenditure.Amount)
	return expenditure, nil
}

func (m *MemoryService) GetAllExpenditures() ([]*domain.Expenditure, error) {
	m.logger.Debug("Getting all expenditures")

	m.RLock()
	defer m.RUnlock()

	expenditures := make([]*domain.Expenditure, 0, len(m.Expenditures))
	for _, expenditure := range m.Expenditures {
		expenditures = append(expenditures, expenditure)
	}

	m.logger.Info("Retrieved all expenditures", "count", len(expenditures))
	return expenditures, nil
}

func (m *MemoryService) UpdateExpenditure(expenditure *domain.Expenditure) error {
	m.logger.Debug("Updating expenditure", "id", expenditure.ID, "description", expenditure.Description, "amount", expenditure.Amount)

	m.Lock()
	defer m.Unlock()

	id := expenditure.ID.String()
	if _, exists := m.Expenditures[id]; !exists {
		m.logger.Warn("Expenditure not found for update", "id", id)
		return domain.ErrExpenditureNotFound
	}

	m.Expenditures[id] = expenditure
	m.logger.Info("Expenditure updated successfully", "id", id)
	return nil
}

func (m *MemoryService) DeleteExpenditure(id string) error {
	m.logger.Debug("Deleting expenditure", "id", id)

	m.Lock()
	defer m.Unlock()

	if _, exists := m.Expenditures[id]; !exists {
		m.logger.Warn("Expenditure not found for deletion", "id", id)
		return domain.ErrExpenditureNotFound
	}

	delete(m.Expenditures, id)
	m.logger.Info("Expenditure deleted successfully", "id", id, "remaining_count", len(m.Expenditures))
	return nil
}
