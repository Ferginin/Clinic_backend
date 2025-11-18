package entity

import (
	"errors"
	"time"
)

type User struct {
	ID                  int       `json:"id"`
	Username            string    `json:"username" binding:"required"`
	Email               string    `json:"email" binding:"required,email"`
	Provider            *string   `json:"provider,omitempty"`
	Password            string    `json:"password,omitempty" binding:"required,min=6"`
	ResetPasswordToken  *string   `json:"-"`
	ConfirmationToken   *string   `json:"-"`
	Confirmed           bool      `json:"confirmed"`
	Blocked             bool      `json:"blocked"`
	RoleID              *int      `json:"role_id"`
	RoleName            string    `json:"role_name,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Confirmed bool      `json:"confirmed"`
	Blocked   bool      `json:"blocked"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Confirmed: u.Confirmed,
		Blocked:   u.Blocked,
		RoleName:  u.RoleName,
		CreatedAt: u.CreatedAt,
	}
}
