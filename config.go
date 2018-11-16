package main

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

type config struct{}

// ListenPort returns application port
func (c config) ListenPort() string {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	return fmt.Sprintf(":%s", port)
}

// APIVersion returns current API version
func (c config) APIVersion() string {
	return fmt.Sprintf("v%s", os.Getenv("API_VERSION"))
}

// DebugMode determines whether debug mode is enabled or not
func (c config) DebugMode() bool {
	mode, err := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	if err != nil {
		return false
	}
	return mode
}

// AutoTLSMode says to use authogenerated LetsEncrypt certificate
func (c config) AutoTLSMode() bool {
	mode, err := strconv.ParseBool(os.Getenv("TLS_AUTO"))
	if err != nil {
		return false
	}
	return mode
}

// CrtDir returns path to store generated LetsEncrypt certificate
func (c config) CrtDir() string {
	return os.Getenv("TLS_DIR")
}

// CustomCrtPath returns path to custom TLS certificate
func (c config) CustomCrtPath() string {
	return os.Getenv("TLS_CRT_PATH")
}

// CustomCrtKeyPath returns path to custom TLS certificate key
func (c config) CustomCrtKeyPath() string {
	return os.Getenv("TLS_KEY_PATH")
}

// GzipLevel func returns response gzip level
func (c config) GzipLevel() int {
	lvl, err := strconv.Atoi(os.Getenv("GZIP_LEVEL"))
	if err != nil {
		return -1
	}
	return lvl
}

// LogLevel returns logs level
func (c config) LogLevel() int {
	lvl, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return 1
	}
	return lvl
}

// APIDomain returns API domain name
func (c config) APIDomain() string {
	return os.Getenv("API_DOMAIN")
}

// Env returns current application environment
func (c config) Env() string {
	return os.Getenv("APP_ENV")
}

// AppName returns the application name
func (c config) AppName() string {
	return os.Getenv("APP_NAME")
}

// JWTSigningKey returns JWT signing key
func (c config) JWTSigningKey() string {
	return os.Getenv("JWT_SECRET")
}
