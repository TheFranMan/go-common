package testing

import (
	"database/sql"
	"os"
	"strings"
	"testing"
)

func ImportFile(t *testing.T, db *sql.DB, filename string) {
	b, err := os.ReadFile("./testdata/" + filename)
	if nil != err {
		t.Fatalf("cannot open SQL file: %s", err.Error())
		return
	}

	statements := strings.Split(strings.TrimSpace(string(b)), ";")

	for _, statement := range statements {
		if 0 == len(statement) {
			continue
		}

		_, err := db.Exec(statement + ";")
		if nil != err {
			t.Fatalf("cannot run SQL statement: %s", err.Error())
			return
		}
	}
}
