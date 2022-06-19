package main

import (
	"log"
	"net"
	"net/http"

	"github.com/nyaneet/music-streaming-backend/db"
	"github.com/nyaneet/music-streaming-backend/handler"
)

type App struct {
	Database db.Database
}

func (app *App) Initialize(user, password, dbname string) {
	var err error
	app.Database, err = db.Initialize(user, password, dbname)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
}

func (app *App) Run(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error occurred: %s", err.Error())
	}

	httpHandler := handler.NewHandler(app.Database)
	server := &http.Server{
		Handler: httpHandler,
	}

	log.Printf("Started server on %s", addr)

	// go func() {
	server.Serve(listener)
	// }()

}
