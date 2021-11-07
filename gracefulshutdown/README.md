# Gracefulshutdown

`fiber-modules/gracefulshutdown` is a module for providing a grace period between
signaling server to shutdown and actual server closing. In between you can 
handle pending requests and call closing functions `ShutdownFns`.

## Usage

```go
import  "github.com/gkampitakis/fiber-modules/gracefulshutdown"

app := fiber.New()

...

gracefulshutdown.Listen(app, "localhost:8080", gracefulshutdown.WithShutdownFns([]func() error{
  func() error {
      return db.Close()
  }
}))
```

### Options

`fiber-modules/gracefulshutdown` can be configured with three ways:

```go
// Registers gracefulshutdown with default config
gracefulshutdown.Listen(app, "localhost:8080", gracefulshutdown.Default())

// Registers gracefulshutdown with default config and passes ShutdownFns
gracefulshutdown.Listen(app, "localhost:8080", gracefulshutdown.WithShutdownFns())

// or the last method where you can pass the config as you please
gracefulshutdown.Listen(app, "localhost:8080", gracefulshutdown.Config{...})
```

The default values that get applied are: 

- `Default()` sets `Period` to 15 seconds and `Enabled` to true
- `WithShutdownFns()` sets `Period` to 15 seconds and `Enabled` 
to true and allows you to pass

With `gracefulshutdown.Config{}` you can pass:
- `Period` Time to wait before shutting down the server
- `Enabled` Intercept os Signals and provide graceful shutdown when true
- `ShutdownFns` Functions that will execute before shutting down server

  > Make sure Period time is enough for functions to run. Each functions is 
executed in a separate goroutine
- `ListenErrorHandler` Function that handles the error returned from app.Listen.
By default calls log.Fatal 
- `Signals` Signals that get intercepted. Default syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM,