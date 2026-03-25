package model

// UserResponse is the standard user response struct
type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"username"`
	Email     string `json:"email"`
	Token     string `json:"token,omitempty"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// VerifyUserRequest carries a JWT token string for verification
type VerifyUserRequest struct {
	Token string `json:"token,omitempty"`
}

// RegisterUserRequest is the request body for registration
type RegisterUserRequest struct {
	ID       string `json:"id" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=100"`
}

// LoginUserRequest is the request body for login
type LoginUserRequest struct {
	ID       string `json:"id" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

// LogoutUserRequest is the request body for logout
type LogoutUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

// GetUserRequest requests a user by ID
type GetUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
