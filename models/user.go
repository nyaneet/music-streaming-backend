package models

import (
	"net/http"
	"time"
)

type User struct {
	Id               int       `json: "id"`
	Name             string    `json: "name"`
	Email            string    `json: "email"`
	RegistrationDate time.Time `json: "registration_date"`
}

func (*User) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}
