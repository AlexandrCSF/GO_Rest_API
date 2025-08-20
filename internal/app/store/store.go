package store

type Store interface {
	User() UserRepository
	Order() OrderRepository
}

type Config interface {
	DatabaseURL() string
}
