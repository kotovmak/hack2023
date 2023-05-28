package model

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	ID            int    `json:"id"`
	Login         string `json:"login"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	IsKNO         bool   `json:"is_kno"`
	NadzonOrganID int    `json:"nadzor_organ_id"`
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
	ID            int    `json:"id,omitempty"`
	Login         string `json:"-"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Password      string `json:"-"`
	IsKNO         bool   `json:"is_kno,omitempty"`
	NadzonOrganID int    `json:"nadzor_organ_id,omitempty"`
}

type Token struct {
	ID     int    `json:"id"`
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
}
