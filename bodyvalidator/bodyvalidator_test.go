package bodyvalidator

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	jptr "github.com/qri-io/jsonpointer"
	"github.com/qri-io/jsonschema"
	"github.com/stretchr/testify/assert"
)

// This is for testing custom keyword
type IsFoo bool

func newIsFoo() jsonschema.Keyword {
	return new(IsFoo)
}

func (f *IsFoo) Register(uri string, registry *jsonschema.SchemaRegistry) {}

func (f *IsFoo) Resolve(pointer jptr.Pointer, uri string) *jsonschema.Schema {
	return nil
}

// ValidateKeyword implements jsonschema.Keyword
func (f IsFoo) ValidateKeyword(
	ctx context.Context,
	currentState *jsonschema.ValidationState,
	data interface{},
) {
	if str, ok := data.(string); ok {
		if str != "foo" {
			currentState.AddError(data, "invalid foo")
		}
	}
}

func validationRequest(t *testing.T, app *fiber.App, reqBody string) ([]byte, int) {
	req, err := http.NewRequest(
		"POST",
		"/validation",
		strings.NewReader(reqBody),
	)
	if err != nil {
		t.Fatal(err)
	}

	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return body, res.StatusCode
}

func setupValidationRoute(app *fiber.App, middleware func(c *fiber.Ctx) error) {
	app.Post("/validation", middleware, func(c *fiber.Ctx) error {
		return c.Status(200).Send([]byte("Success"))
	})
}

func TestBodyValidator(t *testing.T) {
	t.Run("should use json file for validation", func(t *testing.T) {
		bodyValidator := New()
		v := bodyValidator(Config{
			SchemaPath: "testdata/body-schema.json",
		})

		app := fiber.New()
		setupValidationRoute(app, v)

		payload := `{"user":{"name":"gkampitakis","age":10}}`
		body, statusCode := validationRequest(t, app, payload)
		expectedBody := "Success"

		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, expectedBody, string(body))
	})

	t.Run("should use json file for invalidation", func(t *testing.T) {
		bodyValidator := New()
		v := bodyValidator(Config{
			SchemaPath: "testdata/body-schema.json",
		})

		app := fiber.New()
		setupValidationRoute(app, v)

		payload := `{"user":{"name":"gkampitakis","age":"invalid"}}`
		body, statusCode := validationRequest(t, app, payload)
		expectedBody :=
			`{"message":"Bad request","statusCode":400}`

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, expectedBody, string(body))
	})

	t.Run("should use string literal for validation", func(t *testing.T) {
		bodyValidator := New()
		v := bodyValidator(Config{
			SchemaLiteral: `{
				"type":"object",
				"properties": {
					"user": {
						"type": "object",
						"properties": {
							"name": {
								"type": "string"
							},
							"age": {
								"type": "number"
							}
						},
						"required": ["name"]
					}
				}
			}`,
		})

		app := fiber.New()
		setupValidationRoute(app, v)

		payload := `{"user":{"name":"gkampitakis","age":10}}`
		body, statusCode := validationRequest(t, app, payload)
		expectedBody := "Success"

		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, expectedBody, string(body))
	})

	t.Run("should use string literal for invalidation", func(t *testing.T) {
		bodyValidator := New()
		v := bodyValidator(Config{
			SchemaLiteral: `{
				"type":"object",
				"properties": {
					"user": {
						"type": "object",
						"properties": {
							"name": {
								"type": "string"
							},
							"age": {
								"type": "number"
							}
						},
						"required": ["name"]
					}
				}
			}`,
		})

		app := fiber.New()
		setupValidationRoute(app, v)

		payload := `{"user":{"name":"gkampitakis","age":"invalid"}}`
		body, statusCode := validationRequest(t, app, payload)
		expectedBody :=
			`{"message":"Bad request","statusCode":400}`

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, expectedBody, string(body))
	})

	t.Run("should expose errors in validation", func(t *testing.T) {
		bodyValidator := New(ExposeErrors())
		v := bodyValidator(Config{
			SchemaPath: "testdata/body-schema.json",
		})

		app := fiber.New()
		setupValidationRoute(app, v)

		payload := `{"user":{"name":"gkampitakis","age":"invalid"}}`
		body, statusCode := validationRequest(t, app, payload)
		expectedBody :=
			`{"description":[{"propertyPath":"/user/age","invalidValue":"invalid","message":"type should be number, got string"}],"message":"Bad request","statusCode":400}`

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, expectedBody, string(body))
	})

	t.Run("should handle non json body", func(t *testing.T) {
		bodyValidator := New()
		v := bodyValidator(Config{
			SchemaPath: "testdata/body-schema.json",
		})

		app := fiber.New()
		setupValidationRoute(app, v)

		body, statusCode := validationRequest(t, app, "")
		expectedBody :=
			`{"message":"Bad request: unexpected end of JSON input","statusCode":400}`

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, expectedBody, string(body))
	})

	t.Run("should register 'foo' keyword", func(t *testing.T) {
		bodyValidator := New(ExposeErrors(), RegisterKeywords(Keywords{"foo": newIsFoo}))
		v := bodyValidator(Config{
			SchemaLiteral: `{ "foo": true }`,
		})

		app := fiber.New()
		setupValidationRoute(app, v)

		payload := `"bar"`
		body, statusCode := validationRequest(t, app, payload)
		expectedBody :=
			`{"description":[{"propertyPath":"/","invalidValue":"bar","message":"invalid foo"}],"message":"Bad request","statusCode":400}`

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, expectedBody, string(body))
	})

	t.Run("should exit when invalid schema provided", func(t *testing.T) {
		expectedError := "testdata/missing.json: no such file or directory"
		logFatal = func(l ...interface{}) {
			err, _ := l[0].(error)
			assert.Contains(t, err.Error(), expectedError)
		}

		bodyValidator := New()
		v := bodyValidator(Config{
			SchemaPath: "testdata/missing.json",
		})

		app := fiber.New()
		setupValidationRoute(app, v)
	})

	t.Run("should override default bad request response", func(t *testing.T) {
		response := func(ke []jsonschema.KeyError) interface{} {
			return map[string]interface{}{
				"message": "oops you messed it up there",
			}
		}
		bodyValidator := New(SetResponse(response))
		v := bodyValidator(Config{
			SchemaPath: "testdata/body-schema.json",
		})

		app := fiber.New()
		setupValidationRoute(app, v)

		payload := `{"user":{"name":"gkampitakis","age":"invalid"}}`
		body, statusCode := validationRequest(t, app, payload)
		expectedBody :=
			`{"message":"oops you messed it up there"}`

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, expectedBody, string(body))
	})
}
