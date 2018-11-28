package address

import (
	"database/sql"

	"github.com/labstack/echo"
)

type (
	config interface {
	}

	dbClient interface {
		Select(dest interface{}, query string, args ...interface{}) error
		Get(dest interface{}, query string, args ...interface{}) error
		Exec(query string, args ...interface{}) (sql.Result, error)
	}
)

// SetUp address package
func SetUp(cnf config, router *echo.Echo, db dbClient) error {
	return nil
}
