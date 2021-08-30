package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"io/ioutil"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("postgres",
		"postgres://my_test_user:myTestUser!2021@localhost:5432/snippetbox_test?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	script, err := ioutil.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		script, err := ioutil.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	}
}
