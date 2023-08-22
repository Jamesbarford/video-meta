package database

import (
	"database/sql"
	"fmt"
)

type DbConnection struct {
	db *sql.DB
}

/* Prepare SQL statement to prevent sql injection */
func (conn *DbConnection) PrepareStmt(sqlQuery string, args ...interface{}) (*sql.Rows, error) {
	stmt, prepareErr := conn.db.Prepare(sqlQuery)
	fmt.Printf("%s\n", stmt)
	if prepareErr != nil {
		return nil, prepareErr
	}
	defer stmt.Close()

	rows, queryErr := stmt.Query(args...)
	if queryErr != nil {
		return nil, queryErr
	}

	return rows, nil
}

/**
 * QUERY PASSED IN MUST RESOLVE TO JSON
 * As we are using postgres to generate the JSON we only need the first row
 * and do not need to parse datatypes. This should be fairly safe.
 */
func (conn *DbConnection) ExecGetJsonString(sqlQuery string, args ...interface{}) (string, error) {
	rows, rowsError := conn.PrepareStmt(sqlQuery, args...)
	if rowsError != nil {
		return "", rowsError
	}

	rowResult := ""
	defer rows.Close()
	rows.Next()
	rows.Scan(&rowResult)

	if rowResult == "" {
		rowResult = "{}"
	}

	return rowResult, nil
}

/* A pass through */
func (conn *DbConnection) QueryRow(query string, args ...interface{}) *sql.Row {
	return conn.db.QueryRow(query, args...)
}

func (conn *DbConnection) Exec(query string) (sql.Result, error) {
	return conn.db.Exec(query)
}

func (conn *DbConnection) Query(query string, args ...any) (*sql.Rows, error) {
	return conn.db.Query(query, args...)
}
