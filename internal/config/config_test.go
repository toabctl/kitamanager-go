package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		setEnv       bool
		want         string
	}{
		{
			name:         "returns default when env not set",
			key:          "TEST_CONFIG_UNSET",
			defaultValue: "default_value",
			setEnv:       false,
			want:         "default_value",
		},
		{
			name:         "returns env value when set",
			key:          "TEST_CONFIG_SET",
			defaultValue: "default_value",
			envValue:     "env_value",
			setEnv:       true,
			want:         "env_value",
		},
		{
			name:         "returns default when env is empty string",
			key:          "TEST_CONFIG_EMPTY",
			defaultValue: "default_value",
			envValue:     "",
			setEnv:       true,
			want:         "default_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before test
			os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnv(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnv(%q, %q) = %q, want %q", tt.key, tt.defaultValue, got, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// Save original env values and restore after test
	envVars := []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
		"SERVER_PORT", "JWT_SECRET", "RBAC_MODEL_PATH",
		"SEED_ADMIN_EMAIL", "SEED_ADMIN_PASSWORD", "SEED_ADMIN_NAME",
		"CORS_ALLOW_ORIGINS", "CORS_ALLOW_CREDENTIALS",
		"LOG_LEVEL", "LOG_FORMAT",
	}
	originalValues := make(map[string]string)
	for _, key := range envVars {
		originalValues[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	defer func() {
		for key, value := range originalValues {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	t.Run("returns default values when env not set", func(t *testing.T) {
		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.DBHost != "localhost" {
			t.Errorf("DBHost = %q, want %q", cfg.DBHost, "localhost")
		}
		if cfg.DBPort != "5432" {
			t.Errorf("DBPort = %q, want %q", cfg.DBPort, "5432")
		}
		if cfg.DBUser != "kitamanager" {
			t.Errorf("DBUser = %q, want %q", cfg.DBUser, "kitamanager")
		}
		if cfg.DBPassword != "kitamanager" {
			t.Errorf("DBPassword = %q, want %q", cfg.DBPassword, "kitamanager")
		}
		if cfg.DBName != "kitamanager" {
			t.Errorf("DBName = %q, want %q", cfg.DBName, "kitamanager")
		}
		if cfg.ServerPort != "8080" {
			t.Errorf("ServerPort = %q, want %q", cfg.ServerPort, "8080")
		}
		if cfg.JWTSecret != "default-secret-key" {
			t.Errorf("JWTSecret = %q, want %q", cfg.JWTSecret, "default-secret-key")
		}
		if cfg.LogLevel != "info" {
			t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "info")
		}
		if cfg.LogFormat != "json" {
			t.Errorf("LogFormat = %q, want %q", cfg.LogFormat, "json")
		}
	})

	t.Run("parses CORS origins correctly", func(t *testing.T) {
		os.Setenv("CORS_ALLOW_ORIGINS", "http://localhost:3000, http://example.com , http://test.com")
		defer os.Unsetenv("CORS_ALLOW_ORIGINS")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		expectedOrigins := []string{"http://localhost:3000", "http://example.com", "http://test.com"}
		if len(cfg.CORSAllowOrigins) != len(expectedOrigins) {
			t.Fatalf("CORSAllowOrigins length = %d, want %d", len(cfg.CORSAllowOrigins), len(expectedOrigins))
		}

		for i, origin := range cfg.CORSAllowOrigins {
			if origin != expectedOrigins[i] {
				t.Errorf("CORSAllowOrigins[%d] = %q, want %q", i, origin, expectedOrigins[i])
			}
		}
	})

	t.Run("parses CORS credentials correctly", func(t *testing.T) {
		os.Setenv("CORS_ALLOW_CREDENTIALS", "false")
		defer os.Unsetenv("CORS_ALLOW_CREDENTIALS")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.CORSAllowCredentials != false {
			t.Errorf("CORSAllowCredentials = %v, want false", cfg.CORSAllowCredentials)
		}
	})

	t.Run("uses custom values from env", func(t *testing.T) {
		os.Setenv("DB_HOST", "custom-host")
		os.Setenv("DB_PORT", "5433")
		os.Setenv("SERVER_PORT", "9090")
		os.Setenv("JWT_SECRET", "super-secret")
		defer func() {
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("SERVER_PORT")
			os.Unsetenv("JWT_SECRET")
		}()

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.DBHost != "custom-host" {
			t.Errorf("DBHost = %q, want %q", cfg.DBHost, "custom-host")
		}
		if cfg.DBPort != "5433" {
			t.Errorf("DBPort = %q, want %q", cfg.DBPort, "5433")
		}
		if cfg.ServerPort != "9090" {
			t.Errorf("ServerPort = %q, want %q", cfg.ServerPort, "9090")
		}
		if cfg.JWTSecret != "super-secret" {
			t.Errorf("JWTSecret = %q, want %q", cfg.JWTSecret, "super-secret")
		}
	})
}
