package main

import (
	_ "github.com/joho/godotenv/autoload"

	"api_service/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	router := handlers.SetupRoutes()
	fmt.Println(os.Getenv("DB_PORT"))
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
