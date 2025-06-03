package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Adapter         string `yaml:"adapter"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Database        string `yaml:"database"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	SSLMode         string `yaml:"sslmode"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`
}

type Config struct {
	Database map[string]DatabaseConfig `yaml:",inline"`
}

var (
	envVarRegex = regexp.MustCompile(`\$\{([^}]+)\}`)
)

// LoadConfig loads configuration from YAML file with environment variable substitution
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Substitute environment variables
	content := string(data)
	content = envVarRegex.ReplaceAllStringFunc(content, func(match string) string {
		// Extract variable name and default value
		varExpr := strings.Trim(match, "${}")
		parts := strings.SplitN(varExpr, ":", 2)
		varName := parts[0]
		defaultValue := ""
		if len(parts) > 1 {
			defaultValue = parts[1]
		}

		// Get environment variable or use default
		if value := os.Getenv(varName); value != "" {
			return value
		}
		return defaultValue
	})

	var config Config
	config.Database = make(map[string]DatabaseConfig)

	// Unmarshal directly into the Database map
	if err := yaml.Unmarshal([]byte(content), &config.Database); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// GetEnvironment returns the current environment (development, test, production)
func GetEnvironment() string {
	env := os.Getenv("ENV")
	if env == "" {
		env = os.Getenv("GO_ENV")
	}
	if env == "" {
		env = "development"
	}
	return env
}

// GetDatabaseConfig returns the database configuration for the current environment
func (c *Config) GetDatabaseConfig(env string) (*DatabaseConfig, error) {
	if env == "" {
		env = GetEnvironment()
	}

	dbConfig, exists := c.Database[env]
	if !exists {
		return nil, fmt.Errorf("database configuration for environment '%s' not found", env)
	}

	return &dbConfig, nil
}

// ConnectDatabase establishes database connection using the configuration
func (c *Config) ConnectDatabase(env string) (*gorm.DB, error) {
	dbConfig, err := c.GetDatabaseConfig(env)
	if err != nil {
		return nil, err
	}

	var dialector gorm.Dialector

	switch dbConfig.Adapter {
	case "postgres", "postgresql":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Database, dbConfig.Port, dbConfig.SSLMode)
		dialector = postgres.Open(dsn)

	default:
		return nil, fmt.Errorf("unsupported database adapter: %s", dbConfig.Adapter)
	}

	// Configure GORM logger based on environment
	logLevel := logger.Silent
	if env == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if dbConfig.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	}
	if dbConfig.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	}

	// Parse duration string
	if dbConfig.ConnMaxLifetime != "" {
		duration, err := time.ParseDuration(dbConfig.ConnMaxLifetime)
		if err == nil {
			sqlDB.SetConnMaxLifetime(duration)
		}
	}

	fmt.Printf("Connected to %s database '%s' in %s environment\n",
		dbConfig.Adapter, dbConfig.Database, env)

	return db, nil
}
