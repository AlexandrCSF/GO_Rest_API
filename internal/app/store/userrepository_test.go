package store_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"wb_cource/internal/app/model"
	"wb_cource/internal/app/store"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("users")

	u, err := s.User().Create(model.TestUser(t))
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
