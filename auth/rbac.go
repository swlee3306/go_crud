package auth

import (
	"fmt"
	"go-crud/models"

	"gorm.io/gorm"
)

type Role string

const (
	Admin  Role = "admin"
	User   Role = "user"
	Guest  Role = "guest"
)

type Permission string

const (
	CreateUser Permission = "create:user"
	ReadUser   Permission = "read:user"
	UpdateUser Permission = "update:user"
	DeleteUser Permission = "delete:user"
	
	CreatePost Permission = "create:post"
	ReadPost   Permission = "read:post"
	UpdatePost Permission = "update:post"
	DeletePost Permission = "delete:post"
	
	CreateComment Permission = "create:comment"
	ReadComment   Permission = "read:comment"
	UpdateComment Permission = "update:comment"
	DeleteComment Permission = "delete:comment"
	
	ManageSystem Permission = "manage:system"
)

type UserRole struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserID uint   `gorm:"not null" json:"user_id"`
	Role   string `gorm:"size:50;not null" json:"role"`
	
	User models.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

type RolePermission struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Role         string `gorm:"size:50;not null" json:"role"`
	Permission   string `gorm:"size:100;not null" json:"permission"`
}

type RBACService struct {
	db *gorm.DB
}

func NewRBACService(db *gorm.DB) *RBACService {
	return &RBACService{db: db}
}

func (rs *RBACService) AssignRole(userID uint, role Role) error {
	userRole := &UserRole{
		UserID: userID,
		Role:   string(role),
	}
	
	return rs.db.Create(userRole).Error
}

func (rs *RBACService) RemoveRole(userID uint, role Role) error {
	return rs.db.Where("user_id = ? AND role = ?", userID, string(role)).Delete(&UserRole{}).Error
}

func (rs *RBACService) GetUserRoles(userID uint) ([]string, error) {
	var roles []string
	err := rs.db.Model(&UserRole{}).Where("user_id = ?", userID).Pluck("role", &roles).Error
	return roles, err
}

func (rs *RBACService) HasRole(userID uint, role Role) (bool, error) {
	var count int64
	err := rs.db.Model(&UserRole{}).Where("user_id = ? AND role = ?", userID, string(role)).Count(&count).Error
	return count > 0, err
}

func (rs *RBACService) AddPermission(role Role, permission Permission) error {
	rolePermission := &RolePermission{
		Role:       string(role),
		Permission: string(permission),
	}
	
	return rs.db.Create(rolePermission).Error
}

func (rs *RBACService) RemovePermission(role Role, permission Permission) error {
	return rs.db.Where("role = ? AND permission = ?", string(role), string(permission)).Delete(&RolePermission{}).Error
}

func (rs *RBACService) GetRolePermissions(role Role) ([]string, error) {
	var permissions []string
	err := rs.db.Model(&RolePermission{}).Where("role = ?", string(role)).Pluck("permission", &permissions).Error
	return permissions, err
}

func (rs *RBACService) HasPermission(userID uint, permission Permission) (bool, error) {
	// Get user roles
	roles, err := rs.GetUserRoles(userID)
	if err != nil {
		return false, err
	}
	
	// Check if any role has the permission
	for _, role := range roles {
		var count int64
		err := rs.db.Model(&RolePermission{}).Where("role = ? AND permission = ?", role, string(permission)).Count(&count).Error
		if err != nil {
			return false, err
		}
		if count > 0 {
			return true, nil
		}
	}
	
	return false, nil
}

func (rs *RBACService) RequirePermission(userID uint, permission Permission) error {
	hasPermission, err := rs.HasPermission(userID, permission)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}
	
	if !hasPermission {
		return fmt.Errorf("insufficient permissions: %s required", permission)
	}
	
	return nil
}

func (rs *RBACService) RequireRole(userID uint, role Role) error {
	hasRole, err := rs.HasRole(userID, role)
	if err != nil {
		return fmt.Errorf("failed to check role: %w", err)
	}
	
	if !hasRole {
		return fmt.Errorf("insufficient role: %s required", role)
	}
	
	return nil
}

func (rs *RBACService) RequireAnyRole(userID uint, roles ...Role) error {
	for _, role := range roles {
		hasRole, err := rs.HasRole(userID, role)
		if err != nil {
			return fmt.Errorf("failed to check role: %w", err)
		}
		if hasRole {
			return nil
		}
	}
	
	return fmt.Errorf("insufficient role: one of %v required", roles)
}

func (rs *RBACService) RequireOwnershipOrRole(userID uint, resourceUserID uint, role Role) error {
	// Check if user owns the resource
	if userID == resourceUserID {
		return nil
	}
	
	// Check if user has the required role
	return rs.RequireRole(userID, role)
}

func (rs *RBACService) InitializeDefaultRoles() error {
	// Initialize default role permissions
	defaultPermissions := map[Role][]Permission{
		Admin: {
			CreateUser, ReadUser, UpdateUser, DeleteUser,
			CreatePost, ReadPost, UpdatePost, DeletePost,
			CreateComment, ReadComment, UpdateComment, DeleteComment,
			ManageSystem,
		},
		User: {
			CreatePost, ReadPost, UpdatePost, DeletePost,
			CreateComment, ReadComment, UpdateComment, DeleteComment,
		},
		Guest: {
			ReadPost, ReadComment,
		},
	}
	
	for role, permissions := range defaultPermissions {
		for _, permission := range permissions {
			// Check if permission already exists
			var count int64
			rs.db.Model(&RolePermission{}).Where("role = ? AND permission = ?", string(role), string(permission)).Count(&count)
			
			if count == 0 {
				if err := rs.AddPermission(role, permission); err != nil {
					return fmt.Errorf("failed to add permission %s to role %s: %w", permission, role, err)
				}
			}
		}
	}
	
	return nil
}
