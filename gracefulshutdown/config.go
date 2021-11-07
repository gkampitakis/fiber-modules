package gracefulshutdown

import (
	"os"
	"time"
)

type Config struct {
	// Time to wait before shutting down the server
	Period time.Duration
	// Intercept os Signals and provide graceful shutdown
	Enabled bool
	// Functions that will execute before shutting down server
	// E.g. Close db connections
	ShutdownFns []func() error
	// Function that handles the error returned from app.Listen
	// By default calls log.Fatal
	ListenErrorHandler func(error)
	// Signals that get intercepted
	// Default syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM,
	Signals []os.Signal
}

/* Default values
Period: 15

Enabled: true
*/
func Default() Config {
	return Config{
		Period:  15,
		Enabled: true,
	}
}

/*
Returns a default Config and sets ShutdownFns

Period: 15

Enabled: true
*/
func WithShutdownFns(fns []func() error) Config {
	return Config{
		Period:      15,
		Enabled:     true,
		ShutdownFns: fns,
	}
}
