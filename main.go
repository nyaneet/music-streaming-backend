package main

import "os"

func main() {
	app := App{}

	app.Initialize(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	defer app.Database.Conn.Close()

	app.Run(":8080")
}
