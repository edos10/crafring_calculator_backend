package main

import (
	_ "github.com/joho/godotenv/autoload"

	"api_service/internal/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router, err := handlers.SetupRoutes()
	if err != nil {
		log.Fatal("Failed to setup router ", err)
	}

	fmt.Println("API server is running on port 5051")
	log.Fatal(http.ListenAndServe(":5051", router))
}
