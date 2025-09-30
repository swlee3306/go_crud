package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a post category
type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Slug      string         `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// One-to-Many relationship
	Posts []Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"posts,omitempty"`
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "categories"
}

// Post represents a blog post
type Post struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Slug       string         `gorm:"uniqueIndex;size:250;not null" json:"slug"`
	Content    string         `gorm:"type:longtext;not null" json:"content"`
	Excerpt    string         `gorm:"type:text" json:"excerpt"`
	Status     string         `gorm:"size:20;default:'draft'" json:"status"`
	UserID     uint           `gorm:"not null" json:"user_id"`
	CategoryID uint           `gorm:"not null" json:"category_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// One-to-Many relationships
	User     User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Category Category  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category,omitempty"`
	Comments []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"comments,omitempty"`
}

// TableName specifies the table name for Post
func (Post) TableName() string {
	return "posts"
}

// Comment represents a comment on a post
type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	PostID    uint           `gorm:"not null" json:"post_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// One-to-Many relationships
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Post Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"post,omitempty"`
}

// TableName specifies the table name for Comment
func (Comment) TableName() string {
	return "comments"
}
