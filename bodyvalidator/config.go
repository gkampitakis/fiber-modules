package bodyvalidator

import "github.com/qri-io/jsonschema"

type Config struct {
	SchemaLiteral string
	SchemaPath    string
}

type Keywords map[string]jsonschema.KeyMaker

type globalConfig struct {
	customKeywords Keywords
	exposeErrors   bool
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
