package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"forum/database"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		Eroors(w, r, http.StatusMethodNotAllowed)
		return
	}
	sesionid :=  Checksession(w, r)	
	if !sesionid {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Validate the session
	session  , _:= r.Cookie("session_id")
	var userID int
	err := database.DB.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", session.Value).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			Eroors(w, r, http.StatusInternalServerError)
			fmt.Println("Error fetching user ID:", err)
		}
		return
	}
	//  parse  json data 
	var postContent  struct { 
		Title string `json:"title"`
		Content  string `json:"content"`
	}
	
	err = json.NewDecoder(r.Body).Decode(&postContent)
	if err != nil {
		Eroors(w, r, http.StatusBadRequest)
		return
	}
	switch {
	case postContent.Title == "":
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	case postContent.Content == "":
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}
	// Insert the post into the database
	_, err = database.DB.Exec(`
        INSERT INTO posts (title, content, user_id) 
        VALUES (?, ?, ?)`, postContent.Title, postContent.Content, userID)
	if err != nil {
		Eroors(w, r, http.StatusInternalServerError)
		return
	}

	// Redirect back to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
