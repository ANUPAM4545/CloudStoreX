package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration.
//
// Configuration Precedence Hierarchy (highest to lowest):
//  1. CLI Arguments (if bound)
//  2. Environment Variables (including Kubernetes Secrets / ConfigMaps)
//  3. Configuration Files (YAML/JSON in standard or configurable paths)
//  4. Default Values
type Config struct {
	AppEnv     string `mapstructure:"APP_ENV"`
	Port       string `mapstructure:"PORT"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	RedisHost  string `mapstructure:"REDIS_HOST"`
	RedisPort  string `mapstructure:"REDIS_PORT"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	CORSOrigin string `mapstructure:"CORS_ALLOWED_ORIGINS"`

	MinioEndpoint  string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	MinioBucket    string `mapstructure:"MINIO_BUCKET"`
	MinioUseSSL    bool   `mapstructure:"MINIO_USE_SSL"`

	AwsAccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AwsRegion          string `mapstructure:"AWS_REGION"`
	AwsEndpoint        string `mapstructure:"AWS_ENDPOINT"`
	AwsBucketPrefix    string `mapstructure:"AWS_BUCKET_PREFIX"`

	MaxUploadSizeMB int64 `mapstructure:"MAX_UPLOAD_SIZE_MB"`

	EnablePprof bool   `mapstructure:"ENABLE_PPROF"`
	PprofToken  string `mapstructure:"PPROF_TOKEN"`
}

// LoadConfig loads configuration using standard discovery locations and environment variables.
func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_FILE")
	return LoadConfigWithFile(configPath)
}

// LoadConfigWithFile loads configuration with an optional explicit config file path.
// It enforces the precedence hierarchy: Environment Variables -> Configuration File -> Defaults.
func LoadConfigWithFile(explicitPath string) (*Config, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 1. Set Default Values (lowest precedence)
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("PORT", "8080")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("JWT_SECRET", "super_secret_development_key")
	v.SetDefault("CORS_ALLOWED_ORIGINS", "*")
	v.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	v.SetDefault("MINIO_ACCESS_KEY", "admin")
	v.SetDefault("MINIO_SECRET_KEY", "password123")
	v.SetDefault("MINIO_BUCKET", "cloudstorex-default")
	v.SetDefault("MINIO_USE_SSL", false)

	v.SetDefault("AWS_REGION", "us-east-1")
	v.SetDefault("AWS_ENDPOINT", "")
	v.SetDefault("AWS_BUCKET_PREFIX", "cloudstorex-")

	v.SetDefault("MAX_UPLOAD_SIZE_MB", 100)
	v.SetDefault("ENABLE_PPROF", false)
	v.SetDefault("PPROF_TOKEN", "")

	// 2. Discover and Read Configuration File (middle precedence)
	if explicitPath != "" {
		v.SetConfigFile(explicitPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("/etc/cloudstorex/")
		if home, err := os.UserHomeDir(); err == nil {
			v.AddConfigPath(filepath.Join(home, ".cloudstorex"))
		}
	}

	if err := v.ReadInConfig(); err != nil {
		var notFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &notFoundErr) && !os.IsNotExist(err) {
			// If a config file was explicitly specified or found but had parse/read errors, return error
			return nil, err
		}
		// If file was not found, we continue cleanly with environment variables and defaults
	}

	// 3. Unmarshal into Config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
