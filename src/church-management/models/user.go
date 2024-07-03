package models

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Type        string `json:"type"`
	BirthDate   string `json:"birth_date"`
	BaptismDate string `json:"baptism_date"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
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
