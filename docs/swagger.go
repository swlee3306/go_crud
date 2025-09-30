package docs

// User represents a user in the system
// @Description User account information
type UserResponse struct {
	ID        uint   `json:"id" example:"1"`
	Username  string `json:"username" example:"john_doe"`
	Email     string `json:"email" example:"john@example.com"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	IsActive  bool   `json:"is_active" example:"true"`
	CreatedAt string `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// UserCreateRequest represents the request payload for creating a user
// @Description User creation request
type UserCreateRequest struct {
	Username  string `json:"username" binding:"required" example:"john_doe" minLength:"3" maxLength:"50"`
	Email     string `json:"email" binding:"required" example:"john@example.com" format:"email"`
	Password  string `json:"password" binding:"required" example:"SecurePass123!" minLength:"8"`
	FirstName string `json:"first_name" example:"John" maxLength:"50"`
	LastName  string `json:"last_name" example:"Doe" maxLength:"50"`
}

// UserUpdateRequest represents the request payload for updating a user
// @Description User update request
type UserUpdateRequest struct {
	Username  *string `json:"username,omitempty" example:"john_doe" minLength:"3" maxLength:"50"`
	Email     *string `json:"email,omitempty" example:"john@example.com" format:"email"`
	FirstName *string `json:"first_name,omitempty" example:"John" maxLength:"50"`
	LastName  *string `json:"last_name,omitempty" example:"Doe" maxLength:"50"`
	IsActive  *bool   `json:"is_active,omitempty" example:"true"`
}

// LoginRequest represents the login request payload
// @Description User login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"john@example.com" format:"email"`
	Password string `json:"password" binding:"required" example:"SecurePass123!"`
}

// LoginResponse represents the login response
// @Description Login response with token
type LoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

// ErrorResponse represents an error response
// @Description Error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Code    string `json:"code,omitempty" example:"VALIDATION_ERROR"`
	Details string `json:"details,omitempty" example:"Username is required"`
}

// SuccessResponse represents a success response
// @Description Success response
type SuccessResponse struct {
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginationResponse represents a paginated response
// @Description Paginated response
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination represents pagination information
// @Description Pagination information
type Pagination struct {
	Page       int  `json:"page" example:"1"`
	PerPage    int  `json:"per_page" example:"10"`
	Total      int  `json:"total" example:"100"`
	TotalPages int  `json:"total_pages" example:"10"`
	HasNext    bool `json:"has_next" example:"true"`
	HasPrev    bool `json:"has_prev" example:"false"`
}

// Post represents a blog post
// @Description Blog post information
type PostResponse struct {
	ID         uint   `json:"id" example:"1"`
	Title      string `json:"title" example:"My First Post"`
	Slug       string `json:"slug" example:"my-first-post"`
	Content    string `json:"content" example:"This is the content of my first post"`
	Excerpt    string `json:"excerpt" example:"This is an excerpt..."`
	Status     string `json:"status" example:"published"`
	UserID     uint   `json:"user_id" example:"1"`
	CategoryID uint   `json:"category_id" example:"1"`
	CreatedAt  string `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt  string `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// PostCreateRequest represents the request payload for creating a post
// @Description Post creation request
type PostCreateRequest struct {
	Title      string `json:"title" binding:"required" example:"My First Post" minLength:"1" maxLength:"200"`
	Content    string `json:"content" binding:"required" example:"This is the content of my first post"`
	Excerpt    string `json:"excerpt" example:"This is an excerpt..." maxLength:"500"`
	Status     string `json:"status" example:"draft" enums:"draft,published,archived"`
	CategoryID uint   `json:"category_id" binding:"required" example:"1"`
}

// Comment represents a comment on a post
// @Description Comment information
type CommentResponse struct {
	ID        uint   `json:"id" example:"1"`
	Content   string `json:"content" example:"This is a great post!"`
	UserID    uint   `json:"user_id" example:"1"`
	PostID    uint   `json:"post_id" example:"1"`
	CreatedAt string `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CommentCreateRequest represents the request payload for creating a comment
// @Description Comment creation request
type CommentCreateRequest struct {
	Content string `json:"content" binding:"required" example:"This is a great post!" minLength:"1" maxLength:"1000"`
	PostID  uint   `json:"post_id" binding:"required" example:"1"`
}

// Category represents a post category
// @Description Category information
type CategoryResponse struct {
	ID        uint   `json:"id" example:"1"`
	Name      string `json:"name" example:"Technology"`
	Slug      string `json:"slug" example:"technology"`
	CreatedAt string `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CategoryCreateRequest represents the request payload for creating a category
// @Description Category creation request
type CategoryCreateRequest struct {
	Name string `json:"name" binding:"required" example:"Technology" minLength:"1" maxLength:"100"`
	Slug string `json:"slug" binding:"required" example:"technology" minLength:"1" maxLength:"100"`
}

// HealthCheck represents the health check response
// @Description Health check response
type HealthCheck struct {
	Status    string `json:"status" example:"healthy"`
	Timestamp string `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Version   string `json:"version" example:"1.0.0"`
	Database  string `json:"database" example:"connected"`
}
