package main

import (
	"encoding/json"
	"github.com/aarsinh/auth-golang/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var creds models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Could not decode credentials")
		return
	}

	user, err := app.db.GetUserByEmail(creds.Email)
	if err != nil {
		app.respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		app.respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := app.generateToken(user)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Could not generate token for user: "+err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	var details models.User
	err := json.NewDecoder(r.Body).Decode(&details)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Could not decode signup details")
		app.errorLog.Printf("Failed to decode during signup")
		return
	}

	_, err = app.db.CreateUser(details)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, err.Error())
		app.errorLog.Printf("Failed to create user")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"response": "User created successfully",
	})

}

func (app *application) accessProtected(w http.ResponseWriter, r *http.Request) {
	app.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Welcome to the protected route"})
}
