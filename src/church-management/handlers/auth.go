package handlers

import (
	"encoding/json"
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
	err = db.QueryRow("SELECT id, email, password FROM users WHERE email=$1", creds.Email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
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

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
