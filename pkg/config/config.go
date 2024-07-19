package config

import (
	"os"
)

type AppConfig struct {
	ListenAddr     string
	ZipkinProtocol string
	ZipkinHost     string
	ZipkinPort     string
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func LoadConfig() *AppConfig {
	return &AppConfig{
		ListenAddr:     getEnvWithDefault("LISTEN_ADDR", ":8126"),
		ZipkinProtocol: getEnvWithDefault("ZIPKIN_PROTOCOL", "http"),
		ZipkinHost:     getEnvWithDefault("ZIPKIN_HOST", "localhost"),
		ZipkinPort:     getEnvWithDefault("ZIPKIN_PORT", "9411"),
	}
}
