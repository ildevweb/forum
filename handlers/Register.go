package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"forum/database"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{7,29}$`)
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	Session  ,  code:= Checksession(w, r)
	if code == http.StatusInternalServerError{ 
		Eroors(w, r, code)
		return 
	}else if Session {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/register.html")
		return
	}

	var data struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if data.Email == "" || data.Username == "" || data.Password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}
	}

	email := data.Email
	username := data.Username
	password := data.Password

	if !isValidEmail(email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if !usernameRegex.MatchString(username) {
		http.Error(w, "Invalid username format. Must start with a letter and be 8-30 characters long.", http.StatusBadRequest)
		return
	}

	if !isValidPassword(password) {
		http.Error(w, "Password must be at least 8 characters long, with at least 1 letter and 1 number.", http.StatusBadRequest)
		return
	}

	hashedPassword := password

	_, err := database.DB.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", email, username, hashedPassword)
	if err != nil {
		http.Error(w, "Failed to register user: Email already exists ", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func isValidEmail(email string) bool {
	if !emailRegex.MatchString(email) {
		return false
	}
	localPart, domainPart := splitEmail(email)
	return len(email) <= 320 && len(localPart) <= 64 && len(domainPart) <= 255
}

func splitEmail(email string) (string, string) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// elach hadii adrarii, password khaso ykon hashedPassword
func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLetter := false
	hasNumber := false
	for _, char := range password {
		if char >= 'A' && char <= 'Z' || char >= 'a' && char <= 'z' {
			hasLetter = true
		}
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
		if hasLetter && hasNumber {
			return true
		}
	}
	return false
}
