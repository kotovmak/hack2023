package model

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token,omitempty"`  // AccessToken
	RefreshToken string `json:"refresh_token,omitempty"` // RefreshToken
}

type Account struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
