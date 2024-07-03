package main

import (
	"log"

	"church-management/router"
	"church-management/utils"
)

func main() {
	_, err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	router.StartServer()
}
