package config

import (
	"os"
)

type Config struct {
	Port            string
	MongoURI        string
	MongoDB         string
	MongoCollection string

	OtelServiceName string
	OtelEndpoint    string
	OtelSampler     string
}

func GetEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func Load() Config {
	return Config{
		Port:            GetEnv("PORT", "9000"),
		MongoURI:        GetEnv("MONGO_URI", ""),
		MongoDB:         GetEnv("MONGO_DB", "simple_otel"),
		MongoCollection: GetEnv("MONGO_COLLECTION", "todos"),

		OtelServiceName: GetEnv("OTEL_SERVICE_NAME", "todo-api-go"),
		OtelEndpoint:    GetEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
		OtelSampler:     GetEnv("OTEL_TRACES_SAMPLER", "always_on"),
	}
}
