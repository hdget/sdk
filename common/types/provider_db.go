package types

import "database/sql"

type DbProvider interface {
	Provider
	My() DbClient
	Master() DbClient
	Slave(i int) DbClient
	By(name string) DbClient
}

type DbExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type DbClient interface {
	DbExecutor
	
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Rebind(query string) string
	Close() error
}
