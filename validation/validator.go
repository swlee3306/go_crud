package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", ve.Field, ve.Message)
}

type Validator struct {
	errors []ValidationError
}

func NewValidator() *Validator {
	return &Validator{
		errors: make([]ValidationError, 0),
	}
}

func (v *Validator) AddError(field, message, code string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
		Code:    code,
	})
}

func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

func (v *Validator) GetErrors() []ValidationError {
	return v.errors
}

func (v *Validator) Clear() {
	v.errors = make([]ValidationError, 0)
}

// String validation
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "This field is required", "REQUIRED")
	}
	return v
}

func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.AddError(field, fmt.Sprintf("Minimum length is %d characters", min), "MIN_LENGTH")
	}
	return v
}

func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.AddError(field, fmt.Sprintf("Maximum length is %d characters", max), "MAX_LENGTH")
	}
	return v
}

func (v *Validator) Email(field, value string) *Validator {
	if value != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(value) {
			v.AddError(field, "Invalid email format", "INVALID_EMAIL")
		}
	}
	return v
}

func (v *Validator) Username(field, value string) *Validator {
	if value != "" {
		// Username should be 3-20 characters, alphanumeric and underscores only
		usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
		if !usernameRegex.MatchString(value) {
			v.AddError(field, "Username must be 3-20 characters, alphanumeric and underscores only", "INVALID_USERNAME")
		}
	}
	return v
}

func (v *Validator) Password(field, value string) *Validator {
	if value != "" {
		if len(value) < 8 {
			v.AddError(field, "Password must be at least 8 characters", "PASSWORD_TOO_SHORT")
		}
		
		var hasUpper, hasLower, hasDigit, hasSpecial bool
		for _, char := range value {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsDigit(char):
				hasDigit = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
			}
		}
		
		if !hasUpper {
			v.AddError(field, "Password must contain at least one uppercase letter", "PASSWORD_NO_UPPER")
		}
		if !hasLower {
			v.AddError(field, "Password must contain at least one lowercase letter", "PASSWORD_NO_LOWER")
		}
		if !hasDigit {
			v.AddError(field, "Password must contain at least one digit", "PASSWORD_NO_DIGIT")
		}
		if !hasSpecial {
			v.AddError(field, "Password must contain at least one special character", "PASSWORD_NO_SPECIAL")
		}
	}
	return v
}

// Numeric validation
func (v *Validator) Min(field string, value, min int) *Validator {
	if value < min {
		v.AddError(field, fmt.Sprintf("Value must be at least %d", min), "MIN_VALUE")
	}
	return v
}

func (v *Validator) Max(field string, value, max int) *Validator {
	if value > max {
		v.AddError(field, fmt.Sprintf("Value must be at most %d", max), "MAX_VALUE")
	}
	return v
}

// Custom validation
func (v *Validator) Custom(field string, value interface{}, fn func(interface{}) bool, message string) *Validator {
	if !fn(value) {
		v.AddError(field, message, "CUSTOM_VALIDATION")
	}
	return v
}

// Sanitization
func (v *Validator) SanitizeString(value string) string {
	// Remove HTML tags and trim whitespace
	value = strings.TrimSpace(value)
	// Basic XSS prevention
	value = strings.ReplaceAll(value, "<", "&lt;")
	value = strings.ReplaceAll(value, ">", "&gt;")
	value = strings.ReplaceAll(value, "\"", "&quot;")
	value = strings.ReplaceAll(value, "'", "&#x27;")
	return value
}
