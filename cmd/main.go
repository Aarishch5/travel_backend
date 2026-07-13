package main

import (
	"TravelBackend/database"
	"TravelBackend/server"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	err := database.OpenConnection()
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	defer database.DB.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	myServer := server.SetupRoutes(database.DB)

	log.Printf("server running on port %s", port)
	myServer.Start(":" + port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := myServer.Shutdown(ctx); err != nil {
		log.Fatal("Error shutting down the server")
	}

	log.Println("Server gracefully shut down")
}
