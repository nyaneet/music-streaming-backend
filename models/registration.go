package models

import (
	"fmt"
	"net/http"
)

type RegistrationData struct {
	Password          string `json:"password"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Role              string `json:"role"`
	ArtistName        string `json:"artist_name"`
	ArtistDescription string `json:"artist_description"`
}

func validateUserRole(userRole string) error {
	switch userRole {
	case
		"COMMON",
		"ARTIST":
		return nil
	default:
		return fmt.Errorf("Invalid user role.")
	}
}

func (u *RegistrationData) Bind(req *http.Request) error {
	if u.Username == "" {
		return fmt.Errorf("Username is required.")
	}
	if u.Email == "" {
		return fmt.Errorf("Email is required.")
	}
	if u.Password == "" {
		return fmt.Errorf("Password is required.")
	}
	if err := validateUserRole(u.Role); err != nil {
		return err
	}
	return nil
}
