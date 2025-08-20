package sqlstore

import (
	"database/sql"
	_ "github.com/lib/pq"
	"wb_cource/internal/app/store"
)

type Store struct {
	db              *sql.DB
	databaseURL     string
	userRepository  store.UserRepository
	orderRepository store.OrderRepository
}

func New(databaseURL string) *Store {
	return &Store{
		databaseURL: databaseURL,
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

func (s *Store) Order() store.OrderRepository {
	if s.orderRepository != nil {
		return s.orderRepository
	}

	s.orderRepository = &OrderRepository{
		store: s,
	}

	return s.orderRepository
}

func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.databaseURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
