package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"forum/database"

	"github.com/gofrs/uuid"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	Sesion  , code := Checksession(w, r)
	if code == http.StatusInternalServerError{ 
		Eroors(w, r, code)
		return 
	}
	if  Sesion {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/login.html")
		return
	}
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if data.Username == "" || data.Password == "" {
			http.Error(w, "Both username and password are required", http.StatusBadRequest)
			return
		}
	}

	username := data.Username
	password := data.Password

	var userID int
	var storedPassword string
	Dberror := database.DB.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&userID, &storedPassword)
	if Dberror != nil {
		http.Error(w, "Invalid username ", http.StatusUnauthorized)
		return
	}

	if password != storedPassword {
		http.Error(w, "Invalid  password", http.StatusUnauthorized)
		return
	}

	sessionID, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	_, err = database.DB.Exec(`
    INSERT INTO sessions (session_id, user_id, expires_at, ip_address, user_agent)
    VALUES (?, ?, ?, ?, ?)`, sessionID, userID, time.Now().Add(24*time.Hour), r.RemoteAddr, r.UserAgent())
	if err != nil {
		http.Error(w, "Failed to store session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   sessionID.String(),
		Expires: time.Now().Add(24 * time.Hour),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
