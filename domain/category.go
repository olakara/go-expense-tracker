package domain

import (
	"errors"
	"github.com/google/uuid"
)

var ErrCategoryColorEmpty = errors.New("category color cannot be empty")
var ErrCategoryNameEmpty = errors.New("category name cannot be empty")

type Category struct {
	ID    uuid.UUID `json:"id"`   // Unique identifier for the category
	Name  string    `json:"name"` // Name of the category
	Color string    `json:"color"`
}

func NewCategory(name string, color string) (*Category, error) {
	if name == "" {
		return nil, ErrCategoryNameEmpty
	}

	if color == "" {
		return nil, ErrCategoryColorEmpty
	}

	return &Category{
		ID:    uuid.New(),
		Name:  name,
		Color: color,
	}, nil
}

func (c *Category) Update(name string, color string) error {
	if name == "" {
		return ErrCategoryNameEmpty
	}

	if color == "" {
		return ErrCategoryColorEmpty
	}

	c.Name = name
	c.Color = color

	return nil
}
