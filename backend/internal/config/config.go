package config

import (
	"strings"

	"github.com/spf13/viper"
)

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
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	// This allows mapping env vars like DB_HOST to DBHost struct field
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set defaults
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("JWT_SECRET", "super_secret_development_key")
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "*")
	viper.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	viper.SetDefault("MINIO_ACCESS_KEY", "admin")
	viper.SetDefault("MINIO_SECRET_KEY", "password123")
	viper.SetDefault("MINIO_BUCKET", "cloudstorex-default")
	viper.SetDefault("MINIO_USE_SSL", false)

	viper.SetDefault("AWS_REGION", "us-east-1")
	viper.SetDefault("AWS_ENDPOINT", "")
	viper.SetDefault("AWS_BUCKET_PREFIX", "cloudstorex-")
	
	viper.SetDefault("MAX_UPLOAD_SIZE_MB", 100)

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
