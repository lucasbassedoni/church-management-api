package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

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

	sqlStatement := `
        INSERT INTO users (name, email, password, type, birth_date, is_baptized, address, phone, status) 
        VALUES ($1, $2, $3, $4, to_date($5, 'YYYY-MM-DD'), $6, $7, $8, $9) 
        RETURNING id`
	err = db.QueryRow(sqlStatement, user.Name, user.Email, user.Password, user.Type, user.BirthDate, user.IsBaptized, user.Address, user.Phone, user.Status).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
		log.Printf("Error inserting user into database: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := utils.DB
	rows, err := db.Query("SELECT id, name, email, phone, status FROM users ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		log.Printf("Error fetching users: %v", err)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Status)
		if err != nil {
			http.Error(w, "Error scanning users", http.StatusInternalServerError)
			log.Printf("Error scanning users: %v", err)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var member models.User
	err := json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding member input: %v", err)
		return
	}

	db := utils.DB
	sqlStatement := `
        UPDATE users SET email=$1, phone=$2, status=$3 WHERE id=$4`
	_, err = db.Exec(sqlStatement, member.Email, member.Phone, member.Status, member.ID)
	if err != nil {
		http.Error(w, "Error updating member", http.StatusInternalServerError)
		log.Printf("Error updating member: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	db := utils.DB
	sqlStatement := "DELETE FROM users WHERE id=$1"
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, "Error deleting member", http.StatusInternalServerError)
		log.Printf("Error deleting member: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
