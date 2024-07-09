package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"church-management/models"
	"church-management/utils"
)

func GetDashboardData(w http.ResponseWriter, r *http.Request) {
	db := utils.DB

	// Número total de visitantes
	var totalUsers int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE status = 'Visitante'").Scan(&totalUsers)
	if err != nil {
		http.Error(w, "Error fetching total users count", http.StatusInternalServerError)
		log.Printf("Error fetching total users count: %v", err)
		return
	}

	// Número de membros
	var totalMembers int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE status = 'Membro'").Scan(&totalMembers)
	if err != nil {
		http.Error(w, "Error fetching total members count", http.StatusInternalServerError)
		log.Printf("Error fetching total members count: %v", err)
		return
	}

	// Número de não batizados
	var totalNonBaptized int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE is_baptized = '0'").Scan(&totalNonBaptized)
	if err != nil {
		http.Error(w, "Error fetching total non-baptized count", http.StatusInternalServerError)
		log.Printf("Error fetching total non-baptized count: %v", err)
		return
	}

	// Número de batizados
	var totalBaptized int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE is_baptized = '1'").Scan(&totalBaptized)
	if err != nil {
		http.Error(w, "Error fetching total baptized count", http.StatusInternalServerError)
		log.Printf("Error fetching total baptized count: %v", err)
		return
	}

	// Cadastros recentes
	recentRows, err := db.Query("SELECT name FROM users ORDER BY id DESC LIMIT 5")
	if err != nil {
		http.Error(w, "Error fetching recent users", http.StatusInternalServerError)
		log.Printf("Error fetching recent users: %v", err)
		return
	}
	defer recentRows.Close()

	var recentUsers []string
	for recentRows.Next() {
		var name string
		err := recentRows.Scan(&name)
		if err != nil {
			http.Error(w, "Error scanning recent users", http.StatusInternalServerError)
			log.Printf("Error scanning recent users: %v", err)
			return
		}
		recentUsers = append(recentUsers, name)
	}

	// Aniversariantes do mês
	now := time.Now()
	currentMonth := now.Month()
	currentDay := now.Day()
	birthdayRows, err := db.Query(`
		SELECT name, EXTRACT(DAY FROM birth_date) as day, 
		EXTRACT(YEAR FROM AGE(birth_date)) as age 
		FROM users 
		WHERE EXTRACT(MONTH FROM birth_date) = $1 
		AND EXTRACT(DAY FROM birth_date) >= $2 
		ORDER BY day ASC`, currentMonth, currentDay)
	if err != nil {
		http.Error(w, "Error fetching birthday users", http.StatusInternalServerError)
		log.Printf("Error fetching birthday users: %v", err)
		return
	}
	defer birthdayRows.Close()

	var birthdayUsers []models.BirthdayUser
	for birthdayRows.Next() {
		var user models.BirthdayUser
		err := birthdayRows.Scan(&user.Name, &user.Day, &user.Age)
		if err != nil {
			http.Error(w, "Error scanning birthday users", http.StatusInternalServerError)
			log.Printf("Error scanning birthday users: %v", err)
			return
		}
		birthdayUsers = append(birthdayUsers, user)
	}

	data := models.DashboardData{
		TotalUsers:       totalUsers,
		TotalMembers:     totalMembers,
		TotalNonBaptized: totalNonBaptized,
		TotalBaptized:    totalBaptized,
		RecentUsers:      recentUsers,
		BirthdayUsers:    birthdayUsers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
