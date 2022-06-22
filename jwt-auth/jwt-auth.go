package jwtauth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nyaneet/music-streaming-backend/models"
)

const JWT_LIFETIME = 24 * time.Hour

var JWTKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

type Credentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Credentials) Bind(req *http.Request) error {
	if c.Username == "" {
		return fmt.Errorf("Username field is required.")
	}
	if c.Password == "" {
		return fmt.Errorf("Password field is required.")
	}
	return nil
}

func (c *Credentials) Render(req *http.Request) error {
	return nil
}

type Payload struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func (p *Payload) Bind(req *http.Request) error {
	return nil
}

func (p *Payload) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func GetToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(JWT_LIFETIME)
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(JWTKey)
	if err != nil {
		return "", err
	}

	return token, err
}
