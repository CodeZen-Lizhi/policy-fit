package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Storage  StorageConfig
	LLM      LLMConfig
	Parser   ParserConfig
	Security SecurityConfig
	Log      LogConfig
	Worker   WorkerConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type StorageConfig struct {
	Type      string
	Path      string
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
}

type LLMConfig struct {
	Provider string
	APIKey   string
	BaseURL  string
	Model    string
	Timeout  int
}

type ParserConfig struct {
	PDFParser         string
	PythonServiceURL  string
}

type SecurityConfig struct {
	JWTSecret         string
	DataRetentionDays int
}

type LogConfig struct {
	Level  string
	Format string
}

type WorkerConfig struct {
	Concurrency int
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: viper.GetInt("API_PORT"),
			Mode: viper.GetString("GIN_MODE"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetInt("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		Storage: StorageConfig{
			Type:      viper.GetString("STORAGE_TYPE"),
			Path:      viper.GetString("STORAGE_PATH"),
			Endpoint:  viper.GetString("S3_ENDPOINT"),
			Bucket:    viper.GetString("S3_BUCKET"),
			AccessKey: viper.GetString("S3_ACCESS_KEY"),
			SecretKey: viper.GetString("S3_SECRET_KEY"),
		},
		LLM: LLMConfig{
			Provider: viper.GetString("LLM_PROVIDER"),
			APIKey:   viper.GetString("LLM_API_KEY"),
			BaseURL:  viper.GetString("LLM_BASE_URL"),
			Model:    viper.GetString("LLM_MODEL"),
			Timeout:  viper.GetInt("LLM_TIMEOUT"),
		},
		Parser: ParserConfig{
			PDFParser:        viper.GetString("PDF_PARSER"),
			PythonServiceURL: viper.GetString("PYTHON_SERVICE_URL"),
		},
		Security: SecurityConfig{
			JWTSecret:         viper.GetString("JWT_SECRET"),
			DataRetentionDays: viper.GetInt("DATA_RETENTION_DAYS"),
		},
		Log: LogConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
		Worker: WorkerConfig{
			Concurrency: viper.GetInt("WORKER_CONCURRENCY"),
		},
	}

	// 设置默认值
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Worker.Concurrency == 0 {
		cfg.Worker.Concurrency = 5
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.Format == "" {
		cfg.Log.Format = "json"
	}

	return cfg, nil
}
