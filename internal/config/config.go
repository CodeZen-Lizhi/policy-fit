package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv     string
	ConfigFile string
	Server     ServerConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	Storage    StorageConfig
	LLM        LLMConfig
	Parser     ParserConfig
	Security   SecurityConfig
	Log        LogConfig
	Worker     WorkerConfig
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
	PDFParser        string
	PythonServiceURL string
}

type SecurityConfig struct {
	JWTSecret         string
	DataRetentionDays int
	AdminUserIDs      []int64
}

type LogConfig struct {
	Level  string
	Format string
}

type WorkerConfig struct {
	Concurrency int
}

func Load() (*Config, error) {
	appEnv := detectAppEnv()
	configFile, err := resolveConfigFile(appEnv)
	if err != nil {
		return nil, err
	}

	return loadFromFile(configFile, appEnv, true)
}

func ValidateEnvFile(path string) error {
	if _, err := loadFromFile(path, "", false); err != nil {
		return err
	}
	return nil
}

func loadFromFile(path string, appEnv string, useEnvOverride bool) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if useEnvOverride {
		v.AutomaticEnv()
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	if appEnv == "" {
		appEnv = strings.TrimSpace(strings.ToLower(v.GetString("APP_ENV")))
	}
	if appEnv == "" {
		appEnv = "dev"
	}

	cfg := &Config{
		AppEnv:     appEnv,
		ConfigFile: path,
		Server: ServerConfig{
			Port: v.GetInt("API_PORT"),
			Mode: v.GetString("GIN_MODE"),
		},
		Database: DatabaseConfig{
			Host:     v.GetString("DB_HOST"),
			Port:     v.GetInt("DB_PORT"),
			User:     v.GetString("DB_USER"),
			Password: v.GetString("DB_PASSWORD"),
			DBName:   v.GetString("DB_NAME"),
			SSLMode:  v.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     v.GetString("REDIS_HOST"),
			Port:     v.GetInt("REDIS_PORT"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
		Storage: StorageConfig{
			Type:      v.GetString("STORAGE_TYPE"),
			Path:      v.GetString("STORAGE_PATH"),
			Endpoint:  v.GetString("S3_ENDPOINT"),
			Bucket:    v.GetString("S3_BUCKET"),
			AccessKey: v.GetString("S3_ACCESS_KEY"),
			SecretKey: v.GetString("S3_SECRET_KEY"),
		},
		LLM: LLMConfig{
			Provider: v.GetString("LLM_PROVIDER"),
			APIKey:   v.GetString("LLM_API_KEY"),
			BaseURL:  v.GetString("LLM_BASE_URL"),
			Model:    v.GetString("LLM_MODEL"),
			Timeout:  v.GetInt("LLM_TIMEOUT"),
		},
		Parser: ParserConfig{
			PDFParser:        v.GetString("PDF_PARSER"),
			PythonServiceURL: v.GetString("PYTHON_SERVICE_URL"),
		},
		Security: SecurityConfig{
			JWTSecret:         v.GetString("JWT_SECRET"),
			DataRetentionDays: v.GetInt("DATA_RETENTION_DAYS"),
			AdminUserIDs:      parseInt64List(v.GetString("ADMIN_USER_IDS")),
		},
		Log: LogConfig{
			Level:  v.GetString("LOG_LEVEL"),
			Format: v.GetString("LOG_FORMAT"),
		},
		Worker: WorkerConfig{
			Concurrency: v.GetInt("WORKER_CONCURRENCY"),
		},
	}

	applyDefaults(cfg)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.AppEnv == "" {
		cfg.AppEnv = "dev"
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.Mode == "" {
		if cfg.AppEnv == "prod" {
			cfg.Server.Mode = "release"
		} else {
			cfg.Server.Mode = "debug"
		}
	}
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}
	if cfg.Storage.Type == "" {
		cfg.Storage.Type = "local"
	}
	if cfg.Storage.Path == "" {
		cfg.Storage.Path = "./storage"
	}
	if cfg.LLM.Timeout == 0 {
		cfg.LLM.Timeout = 120
	}
	if cfg.Parser.PDFParser == "" {
		cfg.Parser.PDFParser = "pdftotext"
	}
	if cfg.Security.DataRetentionDays == 0 {
		cfg.Security.DataRetentionDays = 30
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
}

func (c *Config) Validate() error {
	var missing []string

	validateRequired(&missing, c.Database.Host, "DB_HOST")
	validateRequiredInt(&missing, c.Database.Port, "DB_PORT")
	validateRequired(&missing, c.Database.User, "DB_USER")
	validateRequired(&missing, c.Database.DBName, "DB_NAME")
	validateRequired(&missing, c.Redis.Host, "REDIS_HOST")
	validateRequiredInt(&missing, c.Redis.Port, "REDIS_PORT")
	validateRequired(&missing, c.Security.JWTSecret, "JWT_SECRET")
	validateRequiredInt(&missing, c.Security.DataRetentionDays, "DATA_RETENTION_DAYS")
	validateRequired(&missing, c.LLM.Provider, "LLM_PROVIDER")
	validateRequired(&missing, c.LLM.APIKey, "LLM_API_KEY")
	validateRequired(&missing, c.LLM.BaseURL, "LLM_BASE_URL")
	validateRequired(&missing, c.LLM.Model, "LLM_MODEL")
	validateRequiredInt(&missing, c.LLM.Timeout, "LLM_TIMEOUT")
	validateRequired(&missing, c.Parser.PDFParser, "PDF_PARSER")
	validateRequiredInt(&missing, c.Server.Port, "API_PORT")
	validateRequiredInt(&missing, c.Worker.Concurrency, "WORKER_CONCURRENCY")

	switch c.Storage.Type {
	case "local":
		validateRequired(&missing, c.Storage.Path, "STORAGE_PATH")
	case "s3":
		validateRequired(&missing, c.Storage.Endpoint, "S3_ENDPOINT")
		validateRequired(&missing, c.Storage.Bucket, "S3_BUCKET")
		validateRequired(&missing, c.Storage.AccessKey, "S3_ACCESS_KEY")
		validateRequired(&missing, c.Storage.SecretKey, "S3_SECRET_KEY")
	default:
		return fmt.Errorf("invalid STORAGE_TYPE: %s (allowed: local, s3)", c.Storage.Type)
	}

	switch c.Parser.PDFParser {
	case "pdftotext":
	case "python-service":
		validateRequired(&missing, c.Parser.PythonServiceURL, "PYTHON_SERVICE_URL")
	default:
		return fmt.Errorf("invalid PDF_PARSER: %s (allowed: pdftotext, python-service)", c.Parser.PDFParser)
	}

	switch c.AppEnv {
	case "dev", "test", "prod":
	default:
		return fmt.Errorf("invalid APP_ENV: %s (allowed: dev, test, prod)", c.AppEnv)
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config: %s", strings.Join(missing, ", "))
	}

	return nil
}

func validateRequired(missing *[]string, value, key string) {
	if strings.TrimSpace(value) == "" {
		*missing = append(*missing, key)
	}
}

func validateRequiredInt(missing *[]string, value int, key string) {
	if value <= 0 {
		*missing = append(*missing, key)
	}
}

func detectAppEnv() string {
	appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if appEnv == "" {
		return "dev"
	}
	return appEnv
}

func resolveConfigFile(appEnv string) (string, error) {
	if envFile := strings.TrimSpace(os.Getenv("ENV_FILE")); envFile != "" {
		if exists(envFile) {
			absPath, err := filepath.Abs(envFile)
			if err != nil {
				return "", fmt.Errorf("failed to resolve config path %s: %w", envFile, err)
			}
			return absPath, nil
		}
		return "", fmt.Errorf("ENV_FILE not found: %s", envFile)
	}

	candidates := []string{
		fmt.Sprintf(".env.%s.local", appEnv),
		fmt.Sprintf(".env.%s", appEnv),
		".env",
	}

	for _, candidate := range candidates {
		if exists(candidate) {
			absPath, err := filepath.Abs(candidate)
			if err != nil {
				return "", fmt.Errorf("failed to resolve config path %s: %w", candidate, err)
			}
			return absPath, nil
		}
	}

	return "", errors.New("no config file found. expected one of: .env.<env>.local, .env.<env>, .env")
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func parseInt64List(raw string) []int64 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	values := make([]int64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		v, err := strconv.ParseInt(part, 10, 64)
		if err != nil || v <= 0 {
			continue
		}
		values = append(values, v)
	}
	if len(values) == 0 {
		return nil
	}
	return values
}
