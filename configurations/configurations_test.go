package configurations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigurations(t *testing.T) {
	config := DefaultConfigurations()

	if config.MaxMRUDisplay != -1 {
		t.Fatalf(`DefaultConfigurations().MaxMRUDisplay = %d, want -1`, config.MaxMRUDisplay)
	}

	if config.FileFallbackBehavior != true {
		t.Fatalf(`DefaultConfigurations().FileFallbackBehavior = %v, want true`, config.FileFallbackBehavior)
	}
}

func TestGetConfigurations(t *testing.T) {
	tests := []struct {
		name             string
		existingConfig   *Configuration
		expectedMRU      int
		expectedFallback bool
	}{
		{
			name:             "creates default when no config exists",
			existingConfig:   nil,
			expectedMRU:      -1,
			expectedFallback: true,
		},
		{
			name: "loads existing config",
			existingConfig: &Configuration{
				MaxMRUDisplay:        10,
				FileFallbackBehavior: false,
			},
			expectedMRU:      10,
			expectedFallback: false,
		},
		{
			name: "handles partial config with defaults",
			existingConfig: &Configuration{
				MaxMRUDisplay:        50,
				FileFallbackBehavior: true,
			},
			expectedMRU:      50,
			expectedFallback: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalHome := os.Getenv("HOME")
			tempDir, _ := os.MkdirTemp("", "ucd-config-test")
			os.Setenv("HOME", tempDir)
			defer func() {
				os.Setenv("HOME", originalHome)
				os.RemoveAll(tempDir)
			}()

			if tt.existingConfig != nil {
				configPath := filepath.Join(tempDir, ".config", "ucd")
				os.MkdirAll(configPath, 0700)
				configFile := filepath.Join(configPath, "ucd.conf")
				data, _ := json.MarshalIndent(tt.existingConfig, "", "\t")
				os.WriteFile(configFile, data, 0644)
			}

			c := Configuration{}
			result := c.GetConfigurations()

			if result.MaxMRUDisplay != tt.expectedMRU {
				t.Fatalf(`MaxMRUDisplay = %d, want %d`, result.MaxMRUDisplay, tt.expectedMRU)
			}

			if result.FileFallbackBehavior != tt.expectedFallback {
				t.Fatalf(`FileFallbackBehavior = %v, want %v`, result.FileFallbackBehavior, tt.expectedFallback)
			}

			configPath := filepath.Join(tempDir, ".config", "ucd", "ucd.conf")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				t.Fatalf(`config file was not created at %s`, configPath)
			}
		})
	}
}

func TestGetConfigurationsCreatesDirectory(t *testing.T) {
	originalHome := os.Getenv("HOME")
	tempDir, _ := os.MkdirTemp("", "ucd-config-test")
	os.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
		os.RemoveAll(tempDir)
	}()

	configDir := filepath.Join(tempDir, ".config", "ucd")
	if _, err := os.Stat(configDir); !os.IsNotExist(err) {
		t.Fatalf(`config directory should not exist before test`)
	}

	c := Configuration{}
	_ = c.GetConfigurations()

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Fatalf(`config directory was not created at %s`, configDir)
	}
}

func TestGetConfigurationsPreservesExistingConfig(t *testing.T) {
	originalHome := os.Getenv("HOME")
	tempDir, _ := os.MkdirTemp("", "ucd-config-test")
	os.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
		os.RemoveAll(tempDir)
	}()

	configPath := filepath.Join(tempDir, ".config", "ucd")
	os.MkdirAll(configPath, 0700)

	originalConfig := Configuration{
		MaxMRUDisplay:        25,
		FileFallbackBehavior: false,
	}
	configFile := filepath.Join(configPath, "ucd.conf")
	data, _ := json.MarshalIndent(originalConfig, "", "\t")
	os.WriteFile(configFile, data, 0644)

	c := Configuration{}
	result := c.GetConfigurations()

	if result.MaxMRUDisplay != 25 {
		t.Fatalf(`MaxMRUDisplay = %d, want 25`, result.MaxMRUDisplay)
	}

	if result.FileFallbackBehavior != false {
		t.Fatalf(`FileFallbackBehavior = %v, want false`, result.FileFallbackBehavior)
	}
}
