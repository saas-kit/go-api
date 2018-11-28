package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/jmoiron/sqlx"
)

type (
	// DB struct
	DB struct {
		conn *sqlx.DB
	}

	config interface {
		DBLogMode() bool
		DBConnection() string
		DBHost() string
		DBPort() string
		DBName() string
		DBUsername() string
		DBPassword() string
	}
)

// New func is a DB struct factory
// TODO: add connection string formatting for postgres and sqlite3
func New(cnf config) (*DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=30s&parseTime=true",
		cnf.DBUsername(),
		cnf.DBPassword(),
		cnf.DBHost(),
		cnf.DBPort(),
		cnf.DBName(),
	)
	time.Sleep(15 * time.Second)
	db, err := sqlx.Connect("mysql", connString)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(100)
	return &DB{conn: db}, nil
}

// Close database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Select using this DB.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) Select(dest interface{}, query string, args ...interface{}) error {
	return db.conn.Select(dest, query, args...)
}

// Get using this DB.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	return db.conn.Get(dest, query, args...)
}

// Exec sql query
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.conn.Exec(query, args...)
}
