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

func GetFinancialRecords(w http.ResponseWriter, r *http.Request) {
	db := utils.DB
	rows, err := db.Query("SELECT id, type, name, amount, date FROM financial_records ORDER BY date DESC LIMIT 10")
	if err != nil {
		http.Error(w, "Error fetching financial records", http.StatusInternalServerError)
		log.Printf("Error fetching financial records: %v", err)
		return
	}
	defer rows.Close()

	var records []models.FinancialRecord
	for rows.Next() {
		var record models.FinancialRecord
		err := rows.Scan(&record.ID, &record.Type, &record.Name, &record.Amount, &record.Date)
		if err != nil {
			http.Error(w, "Error scanning financial records", http.StatusInternalServerError)
			log.Printf("Error scanning financial records: %v", err)
			return
		}
		records = append(records, record)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func CreateFinancialRecord(w http.ResponseWriter, r *http.Request) {
	var record models.FinancialRecord

	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding financial record input: %v", err)
		return
	}

	db := utils.DB
	sqlStatement := `INSERT INTO financial_records (type, name, amount, date) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, record.Type, record.Name, record.Amount, record.Date)
	if err != nil {
		http.Error(w, "Error inserting financial record", http.StatusInternalServerError)
		log.Printf("Error inserting financial record: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteFinancialRecord(w http.ResponseWriter, r *http.Request) {
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
	sqlStatement := "DELETE FROM financial_records WHERE id=$1"
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, "Error deleting financial record", http.StatusInternalServerError)
		log.Printf("Error deleting financial record: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
