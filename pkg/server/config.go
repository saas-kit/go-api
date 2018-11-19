package server

// Server config interface
type serverConfig interface {
	APIDomain() string
	AutoTLSMode() bool
	CrtDir() string
	CustomCrtPath() string
	CustomCrtKeyPath() string
	DebugMode() bool
	GzipLevel() int
	ListenPort() string
	LogLevel() int
}
