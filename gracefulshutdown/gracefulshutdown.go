package gracefulshutdown

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

var defaultSignals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}

/* Starts fiber server
Calls app.Listen if "Enabled" intercepts os.Signals
and provides some period for handling requests and calling closing functions
and closing connections
*/
func Listen(app *fiber.App, addr string, cfg Config) {
	errorHandler, signals := getDefaults(cfg)

	if !cfg.Enabled {
		startApp(app, addr, errorHandler)
		return
	}

	graceStartApp(
		app,
		addr,
		errorHandler,
		signals,
		cfg.Period,
		cfg.ShutdownFns,
	)
}

func getDefaults(c Config) (func(error), []os.Signal) {
	signals := defaultSignals

	if c.Signals != nil {
		signals = c.Signals
	}

	errorHandler := c.ListenErrorHandler
	if errorHandler != nil {
		errorHandler = func(err error) {
			log.Fatal(err)
		}
	}

	return errorHandler, signals
}

func startApp(app *fiber.App, addr string, errorHandler func(error)) {
	if err := app.Listen(addr); err != nil {
		errorHandler(err)
	}
}

func graceStartApp(
	app *fiber.App,
	addr string,
	errorHandler func(error),
	signals []os.Signal,
	period time.Duration,
	shutdownFns []func() error,
) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)

	go func() {
		<-c

		log.Println("server gracefully shutting down")
		err := app.Shutdown()
		if err != nil {
			errorHandler(err)
		}

		log.Println("calling shutdown functions")

		for _, fn := range shutdownFns {
			go executeFn(fn)
		}
	}()

	startApp(app, addr, errorHandler)

	select {
	case <-c:
		log.Println("received second signal")
	case <-time.After(time.Second * period):
	}
}

func executeFn(fn func() error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	if err := fn(); err != nil {
		log.Println(err)
	}
}
