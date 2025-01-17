package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"forum/database"

	_ "github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var Interface, InterfaceError = template.ParseGlob("./templates/*.html")

type Post struct {
	ID        int
	Title     string
	Content   string
	Likes     int
	Deslikes  int
	Username  string
	CreatedAt string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	Sesion := Checksession(w, r)

	if !Sesion {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Clear LocalStorage</title>
			</head>
			<body>
				<script>
					// Clear localStorage
					localStorage.clear();
				</script>
			</body>
			</html>
		`)
	}

	if r.URL.Path != "/" && r.URL.Path != "/home" {
		Eroors(w, r, http.StatusNotFound)
		return
	}
	rows, err := database.DB.Query(`
        SELECT posts.id, posts.title, posts.content, posts.Likes, posts.Deslikes, users.username, posts.created_at
        FROM posts
        JOIN users ON posts.user_id = users.id
        ORDER BY posts.created_at DESC
    `)
	if err != nil {
		http.Error(w, "Failed to fetch posts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Likes, &post.Deslikes, &post.Username, &post.CreatedAt)
		if err != nil {
			http.Error(w, "Something went wrong: ", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	// Check for errors after iterating rows
	if err := rows.Err(); err != nil {
		http.Error(w, "Something went wrong: ", http.StatusInternalServerError)
		return
	}
	rows.Close()
	// If no posts are found, set to nil
	if len(posts) == 0 {
		posts = nil
	}
	// Render the home page template with posts
	Interface.ExecuteTemplate(w, "home.html", map[string]interface{}{
		"Posts":    posts,
		"LoggedIn": Sesion,
	})
}

func Checksession(w http.ResponseWriter, r *http.Request) bool {
	userSession, err := r.Cookie("session_id")
	if err != nil {
		return false
	}
	Db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		Eroors(w, r, http.StatusInternalServerError)
		return false
	}
	defer Db.Close()
	row := Db.QueryRow("SELECT username FROM sessions WHERE session_id = ?", userSession.Value)
	return row != nil
}

func Eroors(w http.ResponseWriter, r *http.Request, code int) {
	Interface.ExecuteTemplate(w, "Error.html", map[string]interface{}{
		"Code":    code,
		"Message": http.StatusText(code),
	})
}
