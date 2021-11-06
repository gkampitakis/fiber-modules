package healthcheck

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

var upTime = time.Now()

type HealthCheckResponse struct {
	Service      string            `json:"service,omitempty"`
	Uptime       string            `json:"uptime"`
	Memory       MemoryMetrics     `json:"memory"`
	GoRoutines   int               `json:"go_routines"`
	HealthChecks map[string]string `json:"health_checks,omitempty"`
}

type MemoryMetrics struct {
	ResidentSetSize  uint64 `json:"rss"`
	TotalAlloc       uint64 `json:"total_alloc"`
	HeapAlloc        uint64 `json:"heap_alloc"`
	HeapObjectsCount uint64 `json:"heap_objects_count"`
}

type CheckResult struct {
	msg   string
	label string
}

func New(options ...func(*Config)) fiber.Handler {
	cfg := Config{}

	for _, o := range options {
		o(&cfg)
	}

	return registerHealthcheck(cfg)
}

// @Description Route reporting health of service
// @Summary Healthcheck route
// @Tags health
// @Accept text/plain
// @Product json/application
// @Success 200 {object} map[string]string "This can be dynamic and add more fields in checks"
// @Failure 500 {object} map[string]string "The route can return 500 in case of failed check,timeouts or panic"
// @Router /health [get]
func registerHealthcheck(cfg Config) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		if cfg.HealthChecks == nil {
			return ctx.Status(http.StatusOK).JSON(prepareResponse(cfg.ServiceName, nil))
		}

		status := http.StatusOK
		response := prepareResponse(cfg.ServiceName, map[string]string{})
		c := make(chan CheckResult)
		checksLength := len(cfg.HealthChecks)

		for label, control := range cfg.HealthChecks {
			go check(
				label,
				control,
				&status,
				c,
				cfg,
			)
		}

		for i := 0; i < checksLength; i++ {
			checkResponse := <-c
			response.HealthChecks[checkResponse.label] = checkResponse.msg
		}

		return ctx.Status(status).JSON(response)
	}
}

func prepareResponse(serviceName string, checks map[string]string) *HealthCheckResponse {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return &HealthCheckResponse{
		Uptime:  time.Since(upTime).String(),
		Service: serviceName,
		Memory: MemoryMetrics{
			ResidentSetSize:  mem.HeapSys,
			TotalAlloc:       mem.TotalAlloc,
			HeapAlloc:        mem.HeapAlloc,
			HeapObjectsCount: mem.HeapObjects,
		},
		GoRoutines:   runtime.NumGoroutine(),
		HealthChecks: checks,
	}
}

func check(
	label string,
	control func() error,
	status *int,
	c chan<- CheckResult,
	cfg Config,
) {
	internalChan := make(chan CheckResult, 1)

	go func() {
		defer func() {
			if e := recover(); e != nil {
				internalChan <- CheckResult{msg: fmt.Errorf("paniced with error: %v", e).Error(), label: label}
				if *status == http.StatusOK {
					*status = http.StatusInternalServerError
				}
			}
		}()

		err := control()
		if err == nil {
			internalChan <- CheckResult{msg: "healthy", label: label}
			return
		}

		if *status == http.StatusOK {
			*status = http.StatusInternalServerError
		}

		msg := "unhealthy"

		if cfg.ShowErrors {
			msg = err.Error()
		}

		internalChan <- CheckResult{msg, label}
	}()

	if cfg.TimeoutEnabled {
		select {
		case tmp := <-internalChan:
			c <- tmp
		case <-time.After(time.Second * cfg.TimeoutPeriod):
			c <- CheckResult{msg: fmt.Sprintf("Timeout after %d seconds", cfg.TimeoutPeriod), label: label}
		}

		return
	}

	c <- <-internalChan
}
