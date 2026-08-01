package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudstorex/backend/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Defaults(t *testing.T) {
	cfg, err := config.LoadConfigWithFile("/path/that/does/not/exist/config.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "development", cfg.AppEnv)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, int64(100), cfg.MaxUploadSizeMB)
}

func TestLoadConfig_ConfigFileFallback(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	content := []byte(`
APP_ENV: staging
PORT: "9090"
DB_HOST: "staging-db"
MAX_UPLOAD_SIZE_MB: 250
`)
	err := os.WriteFile(configPath, content, 0600)
	require.NoError(t, err)

	cfg, err := config.LoadConfigWithFile(configPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "staging", cfg.AppEnv)
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "staging-db", cfg.DBHost)
	assert.Equal(t, int64(250), cfg.MaxUploadSizeMB)
}

func TestLoadConfig_EnvVarOverridePrecedence(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	content := []byte(`
APP_ENV: staging
PORT: "9090"
`)
	err := os.WriteFile(configPath, content, 0600)
	require.NoError(t, err)

	t.Setenv("APP_ENV", "production")
	t.Setenv("PORT", "8443")
	t.Setenv("DB_HOST", "prod-postgres")

	cfg, err := config.LoadConfigWithFile(configPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Env vars should override config file values
	assert.Equal(t, "production", cfg.AppEnv)
	assert.Equal(t, "8443", cfg.Port)
	assert.Equal(t, "prod-postgres", cfg.DBHost)
}
