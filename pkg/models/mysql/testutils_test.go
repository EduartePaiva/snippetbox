package mysql_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	dsn := "test_web:localhost@/test_snippetbox?parseTime=true&multiStatements=true"
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}

	if _, err = db.Exec(string(script)); err != nil {
		t.Fatal(err)
	}

	return db, func() {
		defer db.Close()

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		if _, err = db.Exec(string(script)); err != nil {
			t.Fatal(err)
		}
	}
}
