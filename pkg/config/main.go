package config

import (
	"fmt"
	"os"
	"strconv"
)

// Predefined log levels
const (
	logDebugLevel uint8 = iota + 1
	logInfoLevel
	logWarnLevel
	logErrorLevel
	logOffLevel
	logPanicLevel
	logFatalLevel
)

// New config from .env file
func New() *Config {
	return &Config{}
}

// Config structure
type Config struct{}

// ListenPort returns application port
func (c Config) ListenPort() string {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	return fmt.Sprintf(":%s", port)
}

// APIVersion returns current API version
func (c Config) APIVersion() string {
	return fmt.Sprintf("v%s", os.Getenv("API_VERSION"))
}

// DebugMode determines whether debug mode is enabled or not
func (c Config) DebugMode() bool {
	mode, err := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	if err != nil {
		return false
	}
	return mode
}

// AutoTLSMode says to use authogenerated LetsEncrypt certificate
func (c Config) AutoTLSMode() bool {
	mode, err := strconv.ParseBool(os.Getenv("TLS_AUTO"))
	if err != nil {
		return false
	}
	return mode
}

// CrtDir returns path to store generated LetsEncrypt certificate
func (c Config) CrtDir() string {
	return os.Getenv("TLS_DIR")
}

// CustomCrtPath returns path to custom TLS certificate
func (c Config) CustomCrtPath() string {
	return os.Getenv("TLS_CRT_PATH")
}

// CustomCrtKeyPath returns path to custom TLS certificate key
func (c Config) CustomCrtKeyPath() string {
	return os.Getenv("TLS_KEY_PATH")
}

// GzipLevel func returns response gzip level
func (c Config) GzipLevel() int {
	lvl, err := strconv.Atoi(os.Getenv("GZIP_LEVEL"))
	if err != nil {
		return -1
	}
	return lvl
}

// LogLevel returns logs level
func (c Config) LogLevel() int {
	lvl, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return 1
	}
	return lvl
}

// APIDomain returns API domain name
func (c Config) APIDomain() string {
	return os.Getenv("API_DOMAIN")
}

// Env returns current application environment
func (c Config) Env() string {
	return os.Getenv("APP_ENV")
}

// AppName returns the application name
func (c Config) AppName() string {
	return os.Getenv("APP_NAME")
}

// JWTSigningKey returns JWT signing key
func (c Config) JWTSigningKey() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

// JWTTTL returns JWT life time in seconds
func (c Config) JWTTTL() int64 {
	lvl, err := strconv.Atoi(os.Getenv("JWT_TTL"))
	if err != nil {
		return 1
	}
	return int64(lvl)
}

// DBLogMode determines whether the DB logging mode is enabled
func (c Config) DBLogMode() bool {
	mode := os.Getenv("DB_LOGGING")
	return mode == "true"
}

// DBConnection returns db engine name, e.g.: mysql, postgres, sqlite
func (c Config) DBConnection() string {
	return os.Getenv("DB_CONNECTION")
}

// DBHost returns db host
func (c Config) DBHost() string {
	return os.Getenv("DB_HOST")
}

// DBPort returns db port
func (c Config) DBPort() string {
	return os.Getenv("DB_PORT")
}

// DBName returns db name
func (c Config) DBName() string {
	return os.Getenv("DB_DATABASE")
}

// DBUsername returns db username
func (c Config) DBUsername() string {
	return os.Getenv("DB_USERNAME")
}

// DBPassword returns db password
func (c Config) DBPassword() string {
	return os.Getenv("DB_PASSWORD")
}
