package sql

import (
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"database/sql"
)

var dbs map[string]*sql.DB = make(map[string]*sql.DB)

func Start(dn string, db string) *sql.DB {
	if val, ok := dbs[db]; ok {
		return val
	}

	if dn != "sqlite3" {
		panic(fmt.Sprintf("Driver name %s not recognized", dn))
	}

	// Driver name
	datab, err := sql.Open("sqlite3", db)
	if err == nil {
		dbs[db] = datab
		return datab
	} else {
		panic(fmt.Sprintf("Error in opening database %s", db))
	}
}
