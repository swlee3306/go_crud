package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseManager manages database connections
type DatabaseManager struct {
	DB     *gorm.DB
	Config *DatabaseConfig
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager() (*DatabaseManager, error) {
	config := LoadDatabaseConfig()
	
	// Validate driver
	if err := ValidateDriver(config.Driver); err != nil {
		return nil, fmt.Errorf("database driver validation failed: %w", err)
	}
	
	// Get driver configuration
	driverConfig, err := GetDriverConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get driver config: %w", err)
	}
	
	// Configure GORM logger
	var logLevel logger.LogLevel
	switch getEnv("LOG_LEVEL", "info") {
	case "debug":
		logLevel = logger.Info
	case "info":
		logLevel = logger.Warn
	case "warn":
		logLevel = logger.Error
	default:
		logLevel = logger.Silent
	}
	
	// Open database connection
	db, err := gorm.Open(driverConfig.GetDialector(), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}
	
	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	manager := &DatabaseManager{
		DB:     db,
		Config: config,
	}
	
	// Test connection
	if err := manager.Ping(); err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}
	
	log.Printf("Database connection established successfully (Driver: %s)", config.Driver)
	return manager, nil
}

// Ping tests the database connection
func (dm *DatabaseManager) Ping() error {
	sqlDB, err := dm.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	
	return nil
}

// Close closes the database connection
func (dm *DatabaseManager) Close() error {
	sqlDB, err := dm.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	
	log.Println("Database connection closed")
	return nil
}

// GetDB returns the GORM database instance
func (dm *DatabaseManager) GetDB() *gorm.DB {
	return dm.DB
}

// GetConfig returns the database configuration
func (dm *DatabaseManager) GetConfig() *DatabaseConfig {
	return dm.Config
}

// IsConnected checks if the database is connected
func (dm *DatabaseManager) IsConnected() bool {
	return dm.Ping() == nil
}
