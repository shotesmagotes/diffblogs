package sql

import (
	"database/sql"
	"fmt"
	"errors"
	"strings"
)

const (
	dbname = "./meta.db"
	dbdriver = "sqlite3"
)

type Article struct {
	id uint
	title string
	atype string
	filename string
	published string
	modified string
}

func New(
	title string,
	atype string,
	filename string,
	published string,
	modified string,
) Article {
	a := Article {
		0,
		title,
		atype,
		filename,
		published,
		modified,
	}

	stmt, _ := insertQuery()
	defer stmt.Close()

	res, _ := stmt.Exec(title, atype, filename, published, modified)
	id, _ := res.LastInsertId()

	a.id = uint(id)
	return a
}

func Get(
	kv map[string]interface{},
) (*Article, error) {
	allowed := map[string]bool {
		"id": true,
		"title": true,
		"type": true,
		"filename": true,
		"published": true,
		"modified": true,
	}

	colnames := make([]string, 0)
	values := make([]interface{}, 0)

	for key, value := range(kv) {
		if !allowed[key] {
			panic(fmt.Sprintf("Column %s not present.", key))
		}
		colnames = append(colnames, key)
		values = append(values, value)
	}

	stmt, err := readQuery(colnames)
	defer stmt.Close()

	rows, err := stmt.Query(values...)
	if err != nil {
		err := errors.New(
			fmt.Sprintf(
				"Could not query %v",
				kv,
			),
		)
		return nil, err
	}

	// returns first row retrieved; for multiple rows use GetAll
	rows.Next()
	defer rows.Close()
	a := scan(rows)

	return a, nil
}

// TODO: GetAll function and refactor Get to use get method which cleans the input and returns *sql.Rows

func (a *Article) Set(
	kv map[string]interface{},
) (*Article, error) {
	allowed := map[string]bool {
		"id": true,
		"title": true,
		"type": true,
		"filename": true,
		"published": true,
		"modified": true,
	}

	wcols := make([]string, 0)
	ucols := make([]string, 0)
	uvals := make([]interface{}, 0)

	for key, value := range(kv) {
		if !allowed[key] {
			panic(fmt.Sprintf("Column %s not present.", key))
		}
		ucols = append(ucols, key)
		uvals = append(uvals, value)
	}

	uvals = append(uvals, a.id)
	wcols = append(wcols, "id")
	stmt, err := updateQuery(wcols, ucols)
	defer stmt.Close()

	if err != nil {
		err := errors.New(
			fmt.Sprintf(
				"Could not prepare with values %v",
				uvals,
			),
		)
		return nil, err
	}

	res, err := stmt.Exec(uvals...)
	if err != nil {
		err := errors.New(
			fmt.Sprintf(
				"Could not update rows with values %v",
				uvals,
			),
		)
		return nil, err
	}

	rws, err := res.RowsAffected()
	if rws == 0 {
		return nil, errors.New("Update failed.")
	} else if err != nil {
		return nil, err
	}

	kv["id"] = a.id
	return Get(kv)
}

func readQuery(cols []string) (*sql.Stmt, error) {
	db := Start(dbdriver, dbname)

	qs := `
		SELECT
			id,
			title,
			type,
			filename,
			published,
			modified
		FROM
			articles
		WHERE
	`

	var colvals string
	for _, col := range(cols) {
		colvals += col + "=? AND "
	}
	i := strings.LastIndex(colvals, "AND ")
	colvals = colvals[:i]
	qs += " " + colvals

	qs += `
		ORDER BY id DESC LIMIT 1;
	`
	stmt, err := db.Prepare(qs)
	if err != nil {
		err = errors.New("Could not prepare query.")
		return nil, err
	}

	return stmt, nil
}

func updateQuery(wcols []string, ucols []string) (*sql.Stmt, error) {
	db := Start(dbdriver, dbname)
	qs := `
		UPDATE
			articles
		SET
	`

	// create the update SET string
	var colvals string
	for _, col := range(ucols) {
		colval := col + "=?, "
		colvals += colval
	}
	i := strings.LastIndex(colvals, ", ")
	colvals = colvals[:i]
	qs += (colvals + " WHERE ")

	// create the WHERE string
	colvals = ""
	for _, col := range(wcols) {
		colval := col + "=? AND "
		colvals += colval
	}
	i = strings.LastIndex(colvals, "AND ")
	colvals = colvals[:i]
	qs += " " + colvals + ";"

	stmt, err := db.Prepare(qs)
	if err != nil {
		err = errors.New("Could not prepare update.")
		return nil, err
	}
	return stmt, nil
}

func insertQuery() (*sql.Stmt, error) {
	db := Start(dbdriver, dbname)
	qs := `
		INSERT INTO
			articles (
				title,
				type,
				filename,
				published,
				modified
			)
		VALUES (
			?, ?, ?, ?, ?
		);
	`
	stmt, err := db.Prepare(qs)
	if err != nil {
		return nil, errors.New("Could not prepare insert.")
	}
	return stmt, nil
}

func scan(r *sql.Rows) *Article {
	a := Article{}
	err := r.Scan(
		&a.id,
		&a.title,
		&a.atype,
		&a.filename,
		&a.published,
		&a.modified,
	)

	if err != nil {
		panic(err)
	}

	return &a
}
