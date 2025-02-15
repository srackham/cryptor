package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadConfig_ValidYAML(t *testing.T) {
	// Create a temporary file with valid YAML content.
	tmpFile, err := os.CreateTemp("", "load_config_test")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.WriteString(`xrates-appid: 1234`)
	if err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}

	// Load the test config file.
	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Assert that the config is valid.
	expectedConfig := Config{
		XratesAppId: "1234",
	}

	if diff := cmp.Diff(expectedConfig, *config); diff != "" {
		t.Fatalf("unexpected config: %s", diff)
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Create a temporary file with invalid YAML content.
	tmpFile, err := os.CreateTemp("", "load_config_test")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.WriteString(`INVALID YAML`)
	if err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}

	// Load the YAML file.
	_, err = LoadConfig(tmpFile.Name())
	if err == nil {
		t.Fatalf("expected error when loading invalid YAML")
	}
}

func TestLoadConfig_MissingConfig(t *testing.T) {
	// Load the test config file.
	_, err := LoadConfig("MISSING CONFIG")
	if err == nil {
		t.Fatalf("expected error when loading missing config file")
	}
}
