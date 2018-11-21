package v1

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"

	"saas-kit-api/api/v1/delivery/http"
	"saas-kit-api/api/v1/repositories/mysql"
	"saas-kit-api/api/v1/usecases"
	"saas-kit-api/api/v1/notifications/email"
)

const apiVersion = "v1"

// Logger interface
type logger interface {
	Debug(i ...interface{})
	Info(i ...interface{})
	Warn(i ...interface{})
	Error(i ...interface{})
}

// SetUp API v1 services
func SetUp(e *echo.Echo, db *sqlx.DB, l logger, developers map[string]struct{}, signingKey string, signingDataTTL int64) {
	// Create API route group
	apiGroup := e.Group(apiVersion)

	// Init repositories
	userRepo := mysql.NewUserRepository(db)
	projRepo := mysql.NewProjectRepository(db, userRepo)
	invRepo := mysql.NewInvitationsRepository(db)

	// User notifications provider
	userNotif := email.NewUserNotification()

	// Set up auth web service
	authWs := http.NewAuthWebService(usecases.NewAuthInteractor(
		userRepo, invRepo, projRepo, userNotif, l, developers,signingKey, signingDataTTL,
	))
	authWs.SetUp(apiGroup)
}
