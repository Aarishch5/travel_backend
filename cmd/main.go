package main

import (
	"TravelBackend/database"
	"TravelBackend/handlers"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	err := database.OpenConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.DB.Close()

	log.Println("connected to database")

	registerRoutes(database.DB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func registerRoutes(db *sqlx.DB) {
	// driver's apis
	http.HandleFunc("/create-driver", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateDriver(w, r, db)
	})
	http.HandleFunc("/delete-driver", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteDriver(w, r, db)
	})

	//rider's apis
	http.HandleFunc("/create-rider", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateRider(w, r, db)
	})
	http.HandleFunc("/delete-rider", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteRider(w, r, db)
	})
}
