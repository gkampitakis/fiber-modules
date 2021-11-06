# Healthcheck

`fiber-modules/healthcheck` is a module for creating a health route with custom 
evaluations.

## Usage

```go
import hc "github.com/gkampitakis/fiber-modules/healthcheck"

app := fiber.New()
app.Get("/health", hc.New())
```

### Options

`fiber-modules/healthcheck` can be configured by passing options on `New`

- `EnableTimeout()` if a healthcheck takes more than `TimeoutPeriod` 
will return timeout error. `Default: false`
- `ShowError()` if enabled when a healthcheck errs the response will output 
the error message. `Default: false`
- `SetServiceName` Name of service that will be included in the healthcheck 
response. `Default: ""`
- `SetTimeoutPeriod` Threshold before marking a healthcheck as timed out.
- `RegisterHealthChecks` Map containing healthcheck functions for asserting server
health.

### Example

```go
import hc "github.com/gkampitakis/fiber-modules/healthcheck"

app := fiber.New()
app.Get("/health", hc.New(
  hc.ShowErrors(),
  hc.RegisterHealthChecks(hc.HealthchecksMap{
		"db-connection": func() error {
			time.Sleep(4 * time.Second)
			return errors.New("mock-error")
		})
))
```

The response with this configuration will be
```json
{
  "uptime": "2m11.996690979s",
  "memory": {
    "rss": 24608768,
    "total_alloc": 20907760,
    "heap_alloc": 20609832,
    "heap_objects_count": 3527
  },
  "go_routines": 7,
  "health_checks": {
    "db-connection": "mock-error"
  }
}
```

and the status code will be `500`.