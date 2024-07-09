package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"church-management/models"
	"church-management/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding user input: %v", err)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		log.Printf("Error hashing password: %v", err)
		return
	}
	user.Password = hashedPassword

	db := utils.DB

	// Verificar se o email já está cadastrado
	var existingUser models.User
	err = db.QueryRow("SELECT id FROM users WHERE email=$1", user.Email).Scan(&existingUser.ID)
	if err == nil {
		http.Error(w, "Email já cadastrado", http.StatusConflict)
		return
	}

	if user.BaptismDate == "" {
		user.BaptismDate = "0001-01-01"
	}

	sqlStatement := `
        INSERT INTO users (name, email, password, type, birth_date, baptism_date, address, phone) 
        VALUES ($1, $2, $3, $4, to_date($5, 'YYYY-MM-DD'), to_date($6, 'YYYY-MM-DD'), $7, $8) 
        RETURNING id`
	err = db.QueryRow(sqlStatement, user.Name, user.Email, user.Password, user.Type, user.BirthDate, user.BaptismDate, user.Address, user.Phone).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
		log.Printf("Error inserting user into database: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	// Obtém o usuário autenticado (esta parte pode variar dependendo de como você implementa a autenticação)
	// Aqui estamos assumindo que você consegue obter o email do usuário autenticado
	email := r.Context().Value("email").(string)

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
