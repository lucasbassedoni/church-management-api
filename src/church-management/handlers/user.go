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
