package dto

import "time"

type (
	// Requests
	UserLoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,password"`
	}

	// Responses
	UserDetail struct {
		UserID  string `json:"user_id"`
		Email   string `json:"email"`
		Name    string `json:"first_name"`
		Surname string `json:"last_name"`
	}

	GetUserResponse struct {
		UserID      string    `json:"user_id"`
		Name        string    `json:"first_name"`
		Surname     string    `json:"last_name"`
		Email       string    `json:"email"`
		CreatedAt   time.Time `json:"created_at"`
	}

	UserRegisterRequest struct{
		Name        string    `json:"first_name" validate:"required"`
		Surname     string    `json:"last_name" validate:"required"`
		Email       string    `json:"email" validate:"required"`
		Phone		string	  `json:"phone" validate:"required"`
		Password 	string	  `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirm_password"`
	}

	UserLoginResponse struct {
		UserID       string  `json:"user_id"`
		Name         string  `json:"first_name"`
		Surname      string  `json:"last_name"`
		Email        string  `json:"email"`
		CreatedAt    string  `json:"created_at"`
		AccessToken  string  `json:"access_token"`
		AccessPublic string  `json:"access_public"`
		RefreshToken string  `json:"refresh_token"`
		RefreshPublic string `json:"refresh_public"`
	}
)
