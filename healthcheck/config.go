package healthcheck

import "time"

type HealthchecksMap map[string]func() error

const DefaultTimeoutPeriod = 30

type config struct {
	HealthChecks   HealthchecksMap
	ServiceName    string
	TimeoutPeriod  time.Duration
	TimeoutEnabled bool
	ShowErrors     bool
}

/*
If a healthcheck takes more than "TimeoutPeriod" will return timeout error

"TimeoutPeriod" will default to 30s if not set
*/
func EnableTimeout() func(*config) {
	return func(cfg *config) {
		cfg.TimeoutEnabled = true
		if cfg.TimeoutPeriod == 0 {
			cfg.TimeoutPeriod = DefaultTimeoutPeriod
		}
	}
}

/*
When a healthcheck errs the response will output the error message.
*/
func ShowErrors() func(*config) {
	return func(cfg *config) {
		cfg.ShowErrors = true
	}
}

/*
Threshold before marking a healthcheck as timed out.
*/
func SetTimeoutPeriod(d time.Duration) func(*config) {
	return func(cfg *config) {
		cfg.TimeoutPeriod = d
		if !cfg.TimeoutEnabled {
			cfg.TimeoutEnabled = true
		}
	}
}

/*
Name of service that will be included in the healthcheck response.
*/
func SetServiceName(name string) func(*config) {
	return func(cfg *config) {
		cfg.ServiceName = name
	}
}

/*
Map containing healthcheck functions for asserting server health.
*/
func RegisterHealthChecks(hc HealthchecksMap) func(*config) {
	return func(cfg *config) {
		cfg.HealthChecks = hc
	}
}
