package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found.")
	}

	app := App{}
	app.Initialize(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	defer app.Database.Conn.Close()
	app.Run(":8080")
}
