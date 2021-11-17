package bodyvalidator

import "github.com/qri-io/jsonschema"

type Config struct {
	// String literal representing json schema
	SchemaLiteral string
	// Path to json file containing json schema ( relative to root path )
	SchemaPath string
}

type Keywords map[string]jsonschema.KeyMaker

type BadRequestResponse func([]jsonschema.KeyError) interface{}

type globalConfig struct {
	customKeywords     Keywords
	exposeErrors       bool
	badRequestResponse BadRequestResponse
}

/*
Expose a description where the validation error occurred
*/
func ExposeErrors() func(*globalConfig) {
	return func(cfg *globalConfig) {
		cfg.exposeErrors = true
	}
}

/*
Register custom keywords to jsonschema (https://github.com/qri-io/jsonschema)
*/
func RegisterKeywords(keywords Keywords) func(*globalConfig) {
	return func(cfg *globalConfig) {
		cfg.customKeywords = keywords
	}
}

/*
Override default response with your custom response
*/
func SetResponse(response BadRequestResponse) func(*globalConfig) {
	return func(cfg *globalConfig) {
		cfg.badRequestResponse = response
	}
}
