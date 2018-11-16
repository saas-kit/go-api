package http

// Server config interface
type serverConfig interface {
	routesConfig
	APIDomain() string
	APIVersion() string
	AppName() string
	AutoTLSMode() bool
	CrtDir() string
	CustomCrtPath() string
	CustomCrtKeyPath() string
	DebugMode() bool
	Env() string
	GzipLevel() int
	ListenPort() string
	LogLevel() int
}

// Default routes config interface
type routesConfig interface {
	JWTSigningKey() string
}
