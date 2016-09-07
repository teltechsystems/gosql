package gosql

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type DB interface {
	Begin() (*sql.Tx, error)
	Close() error
	Driver() driver.Driver
	Exec(query string, args ...interface{}) (sql.Result, error)
	Ping() error
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	Stats() sql.DBStats
}
