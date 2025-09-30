package models

import (
	"time"

	"gorm.io/gorm"
)

// Tag represents a post tag
type Tag struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Slug      string         `gorm:"uniqueIndex;size:50;not null" json:"slug"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Many-to-Many relationship
	Posts []Post `gorm:"many2many:post_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"posts,omitempty"`
}

// TableName specifies the table name for Tag
func (Tag) TableName() string {
	return "tags"
}

// PostTag represents the many-to-many relationship between posts and tags
type PostTag struct {
	PostID uint `gorm:"primaryKey" json:"post_id"`
	TagID  uint `gorm:"primaryKey" json:"tag_id"`
}

// TableName specifies the table name for PostTag
func (PostTag) TableName() string {
	return "post_tags"
}
