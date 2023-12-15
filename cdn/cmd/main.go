package main

import (
	_ "github.com/joho/godotenv/autoload"

	"cdn_service/internal/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router, err := handlers.SetupRoutes()
	if err != nil {
		log.Fatal("Failed to setup router ", err)
	}

	fmt.Println("CDN server is running on port 5053")
	log.Fatal(http.ListenAndServe(":5053", router))
}
