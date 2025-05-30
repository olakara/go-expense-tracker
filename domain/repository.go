package domain

import "errors"

var ErrExpenditureAlreadyExists = errors.New("expenditure already exists")
var ErrExpenditureNotFound = errors.New("expenditure not found")

type ExpenditureRepository interface {
	AddExpenditure(expenditure *Expenditure) error
	GetExpenditureByID(id string) (*Expenditure, error)
	GetAllExpenditures() ([]*Expenditure, error)
	UpdateExpenditure(expenditure *Expenditure) error
	DeleteExpenditure(id string) error
}

var ErrCategoryNotFound = errors.New("category not found")

type CategoryRepository interface {
	GetCategoryByID(id string) (*Category, error)
	GetAllCategories() ([]*Category, error)
}
