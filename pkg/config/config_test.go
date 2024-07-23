package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set up environment variables
	os.Setenv("LISTEN_ADDR", ":8080")
	os.Setenv("ZIPKIN_PROTOCOL", "https")
	os.Setenv("ZIPKIN_HOST", "zipkin.io")
	os.Setenv("ZIPKIN_PORT", "9412")

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

	// Clean up environment variables
	os.Unsetenv("LISTEN_ADDR")
	os.Unsetenv("ZIPKIN_PROTOCOL")
	os.Unsetenv("ZIPKIN_HOST")
	os.Unsetenv("ZIPKIN_PORT")

	// Load the config again to check the default values
	config = LoadConfig()

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
