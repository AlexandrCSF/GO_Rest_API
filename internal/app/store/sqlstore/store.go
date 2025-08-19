package sqlstore

import (
	"database/sql"
	_ "github.com/lib/pq"
	"wb_cource/internal/app/store"
)

type Store struct {
	db             *sql.DB
	userRepository store.UserRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
