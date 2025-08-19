package sqlstore_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"wb_cource/internal/app/model"
	"wb_cource/internal/app/store/sqlstore"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := sqlstore.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}
