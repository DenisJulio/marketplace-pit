package utils

type Logger interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}