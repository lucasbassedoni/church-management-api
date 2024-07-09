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

// Credentials represents the credentials for login
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Claims represents the JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
