package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"church-management/models"
	"church-management/utils"

	"github.com/dgrijalva/jwt-go"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user models.User
	db := utils.DB
	err = db.QueryRow("SELECT id, name, email, password, type, birth_date, is_baptized, address, phone FROM users WHERE email=$1", creds.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Type, &user.BirthDate, &user.IsBaptized, &user.Address, &user.Phone)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(20 * time.Minute)
	claims := &models.Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(utils.JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(map[string]string{
		"token":      tokenString,
		"name":       user.Name,
		"email":      user.Email,
		"type":       user.Type,
		"birthDate":  user.BirthDate,
		"isBaptized": user.IsBaptized,
		"address":    user.Address,
		"phone":      user.Phone,
	})
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding input: %v", err)
		return
	}

	// Obtém o token do cookie
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Decodifica o token
	tknStr := c.Value
	claims := &models.Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	email := claims.Email

	var user models.User
	db := utils.DB
	err = db.QueryRow("SELECT id, password FROM users WHERE email=$1", email).Scan(&user.ID, &user.Password)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(creds.CurrentPassword, user.Password) {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := utils.HashPassword(creds.NewPassword)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		log.Printf("Error hashing password: %v", err)
		return
	}

	_, err = db.Exec("UPDATE users SET password=$1 WHERE id=$2", hashedPassword, user.ID)
	if err != nil {
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		log.Printf("Error updating password: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
