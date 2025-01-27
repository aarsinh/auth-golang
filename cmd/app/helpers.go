package main

import (
	"encoding/json"
	"github.com/aarsinh/auth-golang/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"
)

// Respond with an error json
func (app *application) respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		app.errorLog.Printf("responding with a 5XX code")
	}

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	app.respondWithJSON(w, code, ErrorResponse{
		Error: msg,
	})
}

// Function to respond with any json payload
func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		app.errorLog.Printf("Could not marshal payload")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		app.errorLog.Printf("Could not write data")
		return
	}
}

func (app *application) generateToken(user *models.User) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
