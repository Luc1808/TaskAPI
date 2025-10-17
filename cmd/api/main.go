package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Luc1808/TaskAPI/internal/api"
	"github.com/Luc1808/TaskAPI/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (probably running in prod)")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("missing required env var: PORT")
	}

	db, err := repository.InitDB()
	if err != nil {
		log.Fatalf("database init error: %v", err)
	}
	defer db.Close()

	r := api.NewRouter()

	log.Printf("server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
