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

func GetEvents(w http.ResponseWriter, r *http.Request) {
	db := utils.DB
	rows, err := db.Query("SELECT id, name, date, time, location FROM events ORDER BY date DESC LIMIT 10")
	if err != nil {
		http.Error(w, "Error fetching events", http.StatusInternalServerError)
		log.Printf("Error fetching events: %v", err)
		return
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.Date, &event.Time, &event.Location)
		if err != nil {
			http.Error(w, "Error scanning events", http.StatusInternalServerError)
			log.Printf("Error scanning events: %v", err)
			return
		}
		events = append(events, event)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding event input: %v", err)
		return
	}

	db := utils.DB
	sqlStatement := `INSERT INTO events (name, date, time, location) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, event.Name, event.Date, event.Time, event.Location)
	if err != nil {
		http.Error(w, "Error inserting event", http.StatusInternalServerError)
		log.Printf("Error inserting event: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding event input: %v", err)
		return
	}

	db := utils.DB
	sqlStatement := `UPDATE events SET name=$1, date=$2, time=$3, location=$4 WHERE id=$5`
	_, err = db.Exec(sqlStatement, event.Name, event.Date, event.Time, event.Location, event.ID)
	if err != nil {
		http.Error(w, "Error updating event", http.StatusInternalServerError)
		log.Printf("Error updating event: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
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
	sqlStatement := "DELETE FROM events WHERE id=$1"
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, "Error deleting event", http.StatusInternalServerError)
		log.Printf("Error deleting event: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
