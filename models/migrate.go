package models

import (
	"log"

	"gorm.io/gorm"
)

// AutoMigrate runs database migrations for all models
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database migration...")
	
	// Migrate all models
	if err := db.AutoMigrate(
		&User{},
		&UserProfile{},
		&Category{},
		&Post{},
		&Comment{},
		&Tag{},
		&PostTag{},
	); err != nil {
		return err
	}
	
	log.Println("Database migration completed successfully!")
	return nil
}

// TestRelationships tests all model relationships
func TestRelationships(db *gorm.DB) error {
	log.Println("Testing model relationships...")
	
	// Test One-to-One relationship
	var user User
	if err := db.Preload("Profile").First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	
	// Test One-to-Many relationships
	var posts []Post
	if err := db.Preload("User").Preload("Category").Preload("Comments").Find(&posts).Error; err != nil {
		return err
	}
	
	// Test Many-to-Many relationships
	var tags []Tag
	if err := db.Preload("Posts").Find(&tags).Error; err != nil {
		return err
	}
	
	log.Println("All relationship tests passed!")
	return nil
}
