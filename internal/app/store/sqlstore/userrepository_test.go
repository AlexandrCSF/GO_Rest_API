package sqlstore_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"wb_cource/internal/app/model"
	"wb_cource/internal/app/store/sqlstore"
)

func TestUserRepository_Create(t *testing.T) {
	store, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	u := model.TestUser(t)
	assert.NoError(t, store.User().Create(u))
	assert.NotNil(t, u)
}
