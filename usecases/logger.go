package usecases

// Logger interface
type Logger interface {
	Debug(i ...interface{})
	Info(i ...interface{})
	Warn(i ...interface{})
	Error(i ...interface{})
}
