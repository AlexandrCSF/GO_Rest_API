package sqlstore

import (
	"fmt"
	"strings"
	"testing"
)

func TestDB(t *testing.T, databaseURL string) (*Store, func(...string)) {
	t.Helper()
	store := New(databaseURL)
	if err := store.Open(); err != nil {
		t.Fatal(err)
	}

	return store, func(tables ...string) {
		if len(tables) > 0 {
			store.db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
		}
		err := store.Close()
		if err != nil {
			return
		}
	}
}
