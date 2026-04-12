package config_test

import (
	"os"
	"testing"

	"toir-app/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD",
		"DB_NAME", "DB_SSLMODE", "SERVER_PORT", "JWT_SECRET",
	} {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}
}

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_NAME", "toir_test")
	t.Setenv("JWT_SECRET", "test-secret")
}

func TestLoad_DefaultValues(t *testing.T) {
	clearEnv(t)
	setRequiredEnv(t)

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "postgres", cfg.DBUser)
	assert.Equal(t, "", cfg.DBPassword)
	assert.Equal(t, "toir_test", cfg.DBName)
	assert.Equal(t, "disable", cfg.DBSSLMode)
	assert.Equal(t, "8080", cfg.ServerPort)
	assert.Equal(t, "test-secret", cfg.JWTSecret)
}

func TestLoad_MissingRequired(t *testing.T) {
	tests := []struct {
		name    string
		setEnv  func(t *testing.T)
		wantErr string
	}{
		{
			name: "missing DB_HOST",
			setEnv: func(t *testing.T) {
				t.Helper()
				t.Setenv("DB_NAME", "toir")
				t.Setenv("JWT_SECRET", "secret")
			},
			wantErr: "DB_HOST",
		},
		{
			name: "missing DB_NAME",
			setEnv: func(t *testing.T) {
				t.Helper()
				t.Setenv("DB_HOST", "localhost")
				t.Setenv("JWT_SECRET", "secret")
			},
			wantErr: "DB_NAME",
		},
		{
			name: "missing JWT_SECRET",
			setEnv: func(t *testing.T) {
				t.Helper()
				t.Setenv("DB_HOST", "localhost")
				t.Setenv("DB_NAME", "toir")
			},
			wantErr: "JWT_SECRET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			tt.setEnv(t)

			_, err := config.Load()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestConfig_DSN(t *testing.T) {
	cfg := config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "pass",
		DBName:     "toir",
		DBSSLMode:  "disable",
	}

	expected := "host=localhost user=postgres password=pass dbname=toir port=5432 sslmode=disable"
	assert.Equal(t, expected, cfg.DSN())
}
