package models

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Type       string `json:"type"`
	BirthDate  string `json:"birthDate"`
	IsBaptized string `json:"isBaptized"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
	Status     string `json:"status"`
}

// Credenciais para login
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JWT Token Claim
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
