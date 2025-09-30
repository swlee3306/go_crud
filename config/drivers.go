package config

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DatabaseDriver represents supported database drivers
type DatabaseDriver string

const (
	MySQL    DatabaseDriver = "mysql"
	Postgres DatabaseDriver = "postgres"
	SQLite   DatabaseDriver = "sqlite"
)

// DriverConfig holds driver-specific configuration
type DriverConfig struct {
	Driver   DatabaseDriver
	DSN      string
	Options  map[string]interface{}
}

// GetDriverConfig creates driver configuration based on database type
func GetDriverConfig(config *DatabaseConfig) (*DriverConfig, error) {
	driver := DatabaseDriver(strings.ToLower(config.Driver))
	
	switch driver {
	case MySQL:
		return getMySQLConfig(config)
	case Postgres:
		return getPostgresConfig(config)
	case SQLite:
		return getSQLiteConfig(config)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}
}

// getMySQLConfig creates MySQL-specific configuration
func getMySQLConfig(config *DatabaseConfig) (*DriverConfig, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)
	
	return &DriverConfig{
		Driver: MySQL,
		DSN:    dsn,
		Options: map[string]interface{}{
			"charset":   "utf8mb4",
			"parseTime": true,
			"loc":       "Local",
		},
	}, nil
}

// getPostgresConfig creates PostgreSQL-specific configuration
func getPostgresConfig(config *DatabaseConfig) (*DriverConfig, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Seoul",
		config.Host, config.User, config.Password, config.DBName, config.Port)
	
	return &DriverConfig{
		Driver: Postgres,
		DSN:    dsn,
		Options: map[string]interface{}{
			"sslmode":    "disable",
			"timezone":   "Asia/Seoul",
		},
	}, nil
}

// getSQLiteConfig creates SQLite-specific configuration
func getSQLiteConfig(config *DatabaseConfig) (*DriverConfig, error) {
	dsn := config.DBName + ".db"
	
	return &DriverConfig{
		Driver: SQLite,
		DSN:    dsn,
		Options: map[string]interface{}{
			"foreign_keys": true,
		},
	}, nil
}

// GetDialector returns the appropriate GORM dialector for the driver
func (dc *DriverConfig) GetDialector() gorm.Dialector {
	switch dc.Driver {
	case MySQL:
		return mysql.Open(dc.DSN)
	case Postgres:
		return postgres.Open(dc.DSN)
	case SQLite:
		return sqlite.Open(dc.DSN)
	default:
		return nil
	}
}

// ValidateDriver validates if the driver is supported
func ValidateDriver(driver string) error {
	supportedDrivers := []string{"mysql", "postgres", "sqlite"}
	driver = strings.ToLower(driver)
	
	for _, supported := range supportedDrivers {
		if driver == supported {
			return nil
		}
	}
	
	return fmt.Errorf("unsupported database driver: %s. Supported drivers: %v", driver, supportedDrivers)
}
