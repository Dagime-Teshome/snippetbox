package mysql

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	// Because setup and teardown contain multiple SQL statements,
	// we must enable multiStatements=true.
	dsn := "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}

	// Force an actual connection.
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}

	// Run setup script.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(string(script)); err != nil {
		t.Fatal(err)
	}

	// Return cleanup function.
	return db, func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = db.Exec(string(script)); err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
}
