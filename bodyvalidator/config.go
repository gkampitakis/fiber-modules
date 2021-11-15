package bodyvalidator

import "github.com/qri-io/jsonschema"

type Config struct {
	SchemaLiteral string
	SchemaPath    string
}

type Keywords map[string]jsonschema.KeyMaker

type BadRequestResponse func([]jsonschema.KeyError) interface{}

type globalConfig struct {
	customKeywords     Keywords
	exposeErrors       bool
	badRequestResponse BadRequestResponse
}

func ExposeErrors() func(*globalConfig) {
	return func(cfg *globalConfig) {
		cfg.exposeErrors = true
	}
}

func RegisterKeywords(keywords Keywords) func(*globalConfig) {
	return func(cfg *globalConfig) {
		cfg.customKeywords = keywords
	}
}

func SetResponse(response BadRequestResponse) func(*globalConfig) {
	return func(cfg *globalConfig) {
		cfg.badRequestResponse = response
	}
}
