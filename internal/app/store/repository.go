package store

import "wb_cource/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
}

type OrderRepository interface {
	Create(*model.Order) error
	FindByID(string) (*model.Order, error)
	GetAll() ([]*model.Order, error)
}
