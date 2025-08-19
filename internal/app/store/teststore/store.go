package teststore

import (
	"wb_cource/internal/app/model"
	"wb_cource/internal/app/store"
)

type Store struct {
	userRepository store.UserRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}

	return s.userRepository
}
