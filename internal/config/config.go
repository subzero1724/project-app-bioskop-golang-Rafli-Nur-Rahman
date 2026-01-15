package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name string
	Port string
	Env  string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level    string
	Encoding string
}

// Load reads configuration from .env file and environment variables
func Load() (*Config, error) {
	// Set config file settings
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Port: viper.GetString("APP_PORT"),
			Env:  viper.GetString("APP_ENV"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		JWT: JWTConfig{
			Secret:          viper.GetString("JWT_SECRET"),
			ExpirationHours: viper.GetInt("JWT_EXPIRATION_HOURS"),
		},
		Log: LogConfig{
			Level:    viper.GetString("LOG_LEVEL"),
			Encoding: viper.GetString("LOG_ENCODING"),
		},
	}

	// Set defaults if not provided
	if config.App.Port == "" {
		config.App.Port = "8080"
	}
	if config.JWT.ExpirationHours == 0 {
		config.JWT.ExpirationHours = 24
	}
	if config.Log.Level == "" {
		config.Log.Level = "info"
	}
	if config.Log.Encoding == "" {
		config.Log.Encoding = "json"
	}

	return config, nil
}

// GetDatabaseDSN returns PostgreSQL connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetJWTExpiration returns JWT expiration duration
func (c *Config) GetJWTExpiration() time.Duration {
	return time.Duration(c.JWT.ExpirationHours) * time.Hour
}
