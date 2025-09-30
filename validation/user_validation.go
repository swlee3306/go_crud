package validation

import (
	"go-crud/models"
)

type UserCreateRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserUpdateRequest struct {
	Username  *string `json:"username,omitempty"`
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

type UserProfileCreateRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
	Website  string `json:"website"`
	Location string `json:"location"`
}

func ValidateUserCreate(req *UserCreateRequest) []ValidationError {
	validator := NewValidator()
	
	// Sanitize inputs
	req.Username = validator.SanitizeString(req.Username)
	req.Email = validator.SanitizeString(req.Email)
	req.FirstName = validator.SanitizeString(req.FirstName)
	req.LastName = validator.SanitizeString(req.LastName)
	
	// Validate required fields
	validator.Required("username", req.Username).
		Required("email", req.Email).
		Required("password", req.Password)
	
	// Validate field lengths and formats
	validator.MinLength("username", req.Username, 3).
		MaxLength("username", req.Username, 50).
		Username("username", req.Username).
		Email("email", req.Email).
		MaxLength("email", req.Email, 100).
		Password("password", req.Password).
		MaxLength("first_name", req.FirstName, 50).
		MaxLength("last_name", req.LastName, 50)
	
	return validator.GetErrors()
}

func ValidateUserUpdate(req *UserUpdateRequest) []ValidationError {
	validator := NewValidator()
	
	// Validate only provided fields
	if req.Username != nil {
		*req.Username = validator.SanitizeString(*req.Username)
		validator.MinLength("username", *req.Username, 3).
			MaxLength("username", *req.Username, 50).
			Username("username", *req.Username)
	}
	
	if req.Email != nil {
		*req.Email = validator.SanitizeString(*req.Email)
		validator.Email("email", *req.Email).
			MaxLength("email", *req.Email, 100)
	}
	
	if req.FirstName != nil {
		*req.FirstName = validator.SanitizeString(*req.FirstName)
		validator.MaxLength("first_name", *req.FirstName, 50)
	}
	
	if req.LastName != nil {
		*req.LastName = validator.SanitizeString(*req.LastName)
		validator.MaxLength("last_name", *req.LastName, 50)
	}
	
	return validator.GetErrors()
}

func ValidateUserProfileCreate(req *UserProfileCreateRequest) []ValidationError {
	validator := NewValidator()
	
	// Sanitize inputs
	req.Bio = validator.SanitizeString(req.Bio)
	req.Avatar = validator.SanitizeString(req.Avatar)
	req.Website = validator.SanitizeString(req.Website)
	req.Location = validator.SanitizeString(req.Location)
	
	// Validate field lengths
	validator.MaxLength("bio", req.Bio, 500).
		MaxLength("avatar", req.Avatar, 255).
		MaxLength("website", req.Website, 255).
		MaxLength("location", req.Location, 100)
	
	// Validate website URL format if provided
	if req.Website != "" {
		validator.Custom("website", req.Website, func(value interface{}) bool {
			url := value.(string)
			return url == "" || (len(url) > 7 && (url[:7] == "http://" || url[:8] == "https://"))
		}, "Website must be a valid URL starting with http:// or https://")
	}
	
	return validator.GetErrors()
}

func ValidateUser(user *models.User) []ValidationError {
	validator := NewValidator()
	
	// Sanitize inputs
	user.Username = validator.SanitizeString(user.Username)
	user.Email = validator.SanitizeString(user.Email)
	user.FirstName = validator.SanitizeString(user.FirstName)
	user.LastName = validator.SanitizeString(user.LastName)
	
	// Validate required fields
	validator.Required("username", user.Username).
		Required("email", user.Email)
	
	// Validate field lengths and formats
	validator.MinLength("username", user.Username, 3).
		MaxLength("username", user.Username, 50).
		Username("username", user.Username).
		Email("email", user.Email).
		MaxLength("email", user.Email, 100).
		MaxLength("first_name", user.FirstName, 50).
		MaxLength("last_name", user.LastName, 50)
	
	return validator.GetErrors()
}
