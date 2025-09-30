package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-crud/config"
)

func main() {
	log.Println("Starting Go-CRUD application...")
	
	// Test database connection
	if err := config.TestDatabaseConnection(); err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	}
	
	// Create database manager
	dbManager, err := config.NewDatabaseManager()
	if err != nil {
		log.Fatalf("Failed to create database manager: %v", err)
	}
	defer dbManager.Close()
	
	// Print database health
	health := config.TestDatabaseHealth(dbManager)
	log.Printf("Database health: %+v", health)
	
	log.Println("Go-CRUD application started successfully!")
	
	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	<-c
	log.Println("Shutting down Go-CRUD application...")
}
