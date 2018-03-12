package sql

import (
	"testing"
)


func TestStart(t *testing.T){
	defer func() {
		if r := recover(); r == nil {
			t.Error("Panic not registered.")
		}
	}()
	Start("postgres", "./.meta.db")

	db := Start("sqlite3", "./.meta.db")
	err := db.Ping()
	if err != nil {
		t.Error("Was not able to ping the database.")
	}
}
