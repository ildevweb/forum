package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"forum/database"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Eroors(w, r, http.StatusMethodNotAllowed)
		return
	}
	sesionid ,  code:=  Checksession(w, r)
	if code != http.StatusOK {
		Eroors(w, r, code)
		return
	}else  if !sesionid {
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
		Category string `json:"category"`
	}
	err = json.NewDecoder(r.Body).Decode(&postContent)
	if err != nil {
		Eroors(w, r, http.StatusBadRequest)
		return
	}
	title :=  strings.TrimSpace(postContent.Title) ; Content :=  strings.TrimSpace(postContent.Content)

	switch {
	case title== "":
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	case Content == "":
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}
	// ---- Check  Lenght  of  The Content 
	if  len(postContent.Title) >  50  ||  len(postContent.Content) >  300{
		http.Error(w ,  "Something  Whas  Wrong !" ,  http.StatusBadRequest)
		return 
	}
	// ---- Insert the post into the database
	_, err = database.DB.Exec(`
        INSERT INTO posts (title, content, category, user_id) 
        VALUES (?, ?, ?, ?)`, postContent.Title, postContent.Content, postContent.Category, userID)
	if err != nil {
		Eroors(w, r, http.StatusInternalServerError)
		return
	}

	// Redirect back to the home page
	post  , _ :=  GetPostsByUser(w , r )
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status" : "success",
		"post" : post[len(post)-1] ,
	}
	json.NewEncoder(w).Encode(response)
}
