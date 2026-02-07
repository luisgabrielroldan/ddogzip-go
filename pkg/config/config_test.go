package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set up environment variables using t.Setenv (automatically cleaned up)
	t.Setenv("LISTEN_ADDR", ":8080")
	t.Setenv("ZIPKIN_PROTOCOL", "https")
	t.Setenv("ZIPKIN_HOST", "zipkin.io")
	t.Setenv("ZIPKIN_PORT", "9412")

	// Load the config
	config := LoadConfig()

	// Check if the values match the environment variables
	if config.ListenAddr != ":8080" {
		t.Errorf("Expected ListenAddr to be ':8080', but got '%s'", config.ListenAddr)
	}
	if config.ZipkinProtocol != "https" {
		t.Errorf("Expected ZipkinProtocol to be 'https', but got '%s'", config.ZipkinProtocol)
	}
	if config.ZipkinHost != "zipkin.io" {
		t.Errorf("Expected ZipkinHost to be 'zipkin.io', but got '%s'", config.ZipkinHost)
	}
	if config.ZipkinPort != "9412" {
		t.Errorf("Expected ZipkinPort to be '9412', but got '%s'", config.ZipkinPort)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	// Load the config without setting environment variables to check defaults
	config := LoadConfig()

	// Check if the values match the default values
	if config.ListenAddr != ":8126" {
		t.Errorf("Expected ListenAddr to be ':8126', but got '%s'", config.ListenAddr)
	}
	if config.ZipkinProtocol != "http" {
		t.Errorf("Expected ZipkinProtocol to be 'http', but got '%s'", config.ZipkinProtocol)
	}
	if config.ZipkinHost != "localhost" {
		t.Errorf("Expected ZipkinHost to be 'localhost', but got '%s'", config.ZipkinHost)
	}
	if config.ZipkinPort != "9411" {
		t.Errorf("Expected ZipkinPort to be '9411', but got '%s'", config.ZipkinPort)
	}
}
