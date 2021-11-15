package bodyvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/qri-io/jsonschema"
)

var logFatal = log.Fatal

func New(cFns ...func(*globalConfig)) func(Config) fiber.Handler {
	globalCfg := &globalConfig{}

	for _, fn := range cFns {
		fn(globalCfg)
	}

	if globalCfg.badRequestResponse == nil {
		globalCfg.badRequestResponse = defaultResponse(globalCfg.exposeErrors)
	}

	// if keywords are given we need to load the draft again
	if len(globalCfg.customKeywords) > 0 {
		for prop, maker := range globalCfg.customKeywords {
			jsonschema.RegisterKeyword(prop, maker)
		}

		jsonschema.LoadDraft2019_09()
	}

	return func(cfg Config) fiber.Handler {
		v, err := loadValidator(cfg)
		if err != nil {
			logFatal(err)
		}

		return func(c *fiber.Ctx) error {
			validationErrors, err := v.ValidateBytes(c.Context(), c.Body())
			if err != nil {
				return c.Status(http.StatusBadRequest).
					JSON(map[string]interface{}{
						"message":    fmt.Sprintf("Bad request: %s", errors.Unwrap(err)),
						"statusCode": 400,
					})
			}

			if len(validationErrors) > 0 {
				return c.Status(http.StatusBadRequest).
					JSON(globalCfg.badRequestResponse(validationErrors))
			}

			return c.Next()
		}
	}
}

func loadValidator(cfg Config) (*jsonschema.Schema, error) {
	if cfg.SchemaLiteral == "" && cfg.SchemaPath == "" {
		warning("no schema provided")
		return nil, nil
	}
	schema := []byte(cfg.SchemaLiteral)

	if cfg.SchemaLiteral == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		schema, err = os.ReadFile(filepath.Join(wd, cfg.SchemaPath))
		if err != nil {
			return nil, err
		}
	}

	s := new(jsonschema.Schema)
	err := json.Unmarshal(schema, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func defaultResponse(exposeErrors bool) BadRequestResponse {
	return func(errors []jsonschema.KeyError) interface{} {
		res := map[string]interface{}{
			"message":    "Bad request",
			"statusCode": 400,
		}

		if exposeErrors {
			res["description"] = errors
		}

		return res
	}
}

func warning(message string) {
	fmt.Printf("\033[33m[Warning]: %s \033[0m\n", message)
}
