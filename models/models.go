package models

import (
	"time"
)

type (
	Role string
)

type User struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"first_name"`
	Surname   string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone 	  string 	`json:"phone"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct{
	SessionID string `json:"session_id"`
	UserID 		string `json:"user_id"`
	Token 		string `json:"token"`
	IPAdress	string `json:"ip_address"`
	UserAgent	string `json:"user_agent"`
	CreatedAt	time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	LastAccess time.Time `json:"last_access"`
	IsActive bool `json:"is_active"`
	Location string `json:"location"`
}