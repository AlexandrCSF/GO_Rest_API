package teststore

import "wb_cource/internal/app/model"

type UserRepository struct {
	store *Store
	users map[int]*model.User
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	if err := u.BeforeCreate(); err != nil {
		return err
	}
	u.ID = len(r.users) + 1
	r.users[u.ID] = u

	return nil
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	panic("Not Implemented")
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	panic("Not Implemented")
}
