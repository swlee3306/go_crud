package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// TestDatabaseConnection tests the database connection
func TestDatabaseConnection() error {
	log.Println("Testing database connection...")
	
	// Create database manager
	manager, err := NewDatabaseManager()
	if err != nil {
		return fmt.Errorf("failed to create database manager: %w", err)
	}
	defer manager.Close()
	
	// Test basic connection
	if err := manager.Ping(); err != nil {
		return fmt.Errorf("database ping test failed: %w", err)
	}
	log.Println("✓ Database ping test passed")
	
	// Test connection pool
	if err := testConnectionPool(manager); err != nil {
		return fmt.Errorf("connection pool test failed: %w", err)
	}
	log.Println("✓ Connection pool test passed")
	
	// Test transaction
	if err := testTransaction(manager); err != nil {
		return fmt.Errorf("transaction test failed: %w", err)
	}
	log.Println("✓ Transaction test passed")
	
	// Test query execution
	if err := testQueryExecution(manager); err != nil {
		return fmt.Errorf("query execution test failed: %w", err)
	}
	log.Println("✓ Query execution test passed")
	
	log.Println("All database connection tests passed successfully!")
	return nil
}

// testConnectionPool tests the connection pool
func testConnectionPool(manager *DatabaseManager) error {
	sqlDB, err := manager.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	
	// Test multiple connections
	for i := 0; i < 5; i++ {
		if err := sqlDB.Ping(); err != nil {
			return fmt.Errorf("connection pool test failed at iteration %d: %w", i, err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	
	return nil
}

// testTransaction tests database transactions
func testTransaction(manager *DatabaseManager) error {
	return manager.DB.Transaction(func(tx *gorm.DB) error {
		// Simple transaction test
		var result int
		if err := tx.Raw("SELECT 1").Scan(&result).Error; err != nil {
			return fmt.Errorf("transaction query failed: %w", err)
		}
		
		if result != 1 {
			return fmt.Errorf("unexpected transaction result: got %d, expected 1", result)
		}
		
		return nil
	})
}

// testQueryExecution tests basic query execution
func testQueryExecution(manager *DatabaseManager) error {
	var result int
	if err := manager.DB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}
	
	if result != 1 {
		return fmt.Errorf("unexpected query result: got %d, expected 1", result)
	}
	
	return nil
}

// TestDatabaseHealth checks database health status
func TestDatabaseHealth(manager *DatabaseManager) map[string]interface{} {
	health := map[string]interface{}{
		"connected": false,
		"driver":    manager.Config.Driver,
		"host":      manager.Config.Host,
		"port":      manager.Config.Port,
		"database":  manager.Config.DBName,
		"timestamp": time.Now().Unix(),
	}
	
	if manager.IsConnected() {
		health["connected"] = true
		
		// Get connection pool stats
		sqlDB, err := manager.DB.DB()
		if err == nil {
			stats := sqlDB.Stats()
			health["pool_stats"] = map[string]interface{}{
				"max_open_connections":     stats.MaxOpenConnections,
				"open_connections":         stats.OpenConnections,
				"in_use":                   stats.InUse,
				"idle":                     stats.Idle,
				"wait_count":               stats.WaitCount,
				"wait_duration":            stats.WaitDuration.String(),
				"max_idle_closed":          stats.MaxIdleClosed,
				"max_idle_time_closed":     stats.MaxIdleTimeClosed,
				"max_lifetime_closed":      stats.MaxLifetimeClosed,
			}
		}
	}
	
	return health
}
