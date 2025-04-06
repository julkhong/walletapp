package config

import (
	"os"
	"testing"
)

func TestGetEnvWithValue(t *testing.T) {
	err := os.Setenv("TEST_ENV_VAR", "actual_value")
	if err != nil {
		t.Fatalf("failed to set TEST_ENV_VAR: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("TEST_ENV_VAR"); err != nil {
			panic("failed to unset TEST_ENV_VAR: " + err.Error())
		}
	}()

	value := getEnv("TEST_ENV_VAR", "default_value")
	if value != "actual_value" {
		t.Errorf("expected 'actual_value', got '%s'", value)
	}
}

func TestGetEnvFallback(t *testing.T) {
	defer func() {
		if err := os.Unsetenv("NON_EXISTENT_VAR"); err != nil {
			panic("failed to unset NON_EXISTENT_VAR: " + err.Error())
		}
	}()
	value := getEnv("NON_EXISTENT_VAR", "default_value")
	if value != "default_value" {
		t.Errorf("expected 'default_value', got '%s'", value)
	}
}

func TestLoadConfigWithDefaults(t *testing.T) {
	// Unset all known env vars to test fallback
	defer func() {
		if err := os.Unsetenv("DB_HOST"); err != nil {
			panic("failed to unset DB_HOST: " + err.Error())
		}
		if err := os.Unsetenv("DB_PORT"); err != nil {
			panic("failed to unset DB_PORT: " + err.Error())
		}
		if err := os.Unsetenv("DB_USER"); err != nil {
			panic("failed to unset DB_USER: " + err.Error())
		}
		if err := os.Unsetenv("DB_PASSWORD"); err != nil {
			panic("failed to unset DB_PASSWORD: " + err.Error())
		}
		if err := os.Unsetenv("DB_NAME"); err != nil {
			panic("failed to unset DB_NAME: " + err.Error())
		}
		if err := os.Unsetenv("REDIS_HOST"); err != nil {
			panic("failed to unset REDIS_HOST: " + err.Error())
		}
		if err := os.Unsetenv("REDIS_PORT"); err != nil {
			panic("failed to unset REDIS_PORT: " + err.Error())
		}
	}()

	cfg := LoadConfig()

	if cfg.DBHost != "localhost" || cfg.DBPort != "5432" || cfg.DBUser != "wallet_user" ||
		cfg.DBPassword != "wallet_pass" || cfg.DBName != "wallet" {
		t.Errorf("unexpected default DB config: %+v", cfg)
	}

	if cfg.RedisHost != "localhost" || cfg.RedisPort != "6379" {
		t.Errorf("unexpected default Redis config: %+v", cfg)
	}

	expectedURL := "postgres://wallet_user:wallet_pass@localhost:5432/wallet?sslmode=disable"
	if cfg.DBURL != expectedURL {
		t.Errorf("expected DBURL '%s', got '%s'", expectedURL, cfg.DBURL)
	}
}
