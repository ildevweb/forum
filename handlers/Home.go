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
	Type      string
	ID        int
	Title     string
	Content   string
	Category  string
	Likes     int
	Deslikes  int
	Username  string
	CreatedAt string
	Comments  []Comment
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	Sesion, code := Checksession(w, r)
	if code == http.StatusInternalServerError {
		Eroors(w, r, code)
		return
	}
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
        SELECT posts.id, posts.title, posts.content, posts.category, posts.Likes, posts.Deslikes, users.username, posts.created_at
        FROM posts
        JOIN users ON posts.user_id = users.id
        ORDER BY posts.created_at DESC
    `)
	if err != nil {
		http.Error(w, "Failed to fetch posts: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("Error fetching posts:", err)
		return
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Likes, &post.Deslikes, &post.Username, &post.CreatedAt)
		if err != nil {
			http.Error(w, "Something went wrong: ", http.StatusInternalServerError)
			return
		}
		user_id := get_user_id(r)
		username, _ := Getusernamebyid(user_id)
		if post.Username == username {
			post.Type = "created"
		}
		if Check(user_id, post.ID) {
			post.Type = post.Type + " liked"
		}
		//  get  Comment  of   posts

		post.Comments, err = GetComments(w, r, post.ID)
		if err != nil {
			http.Error(w, "Something went wrong: ", http.StatusInternalServerError)
			fmt.Println("Error fetching comments:", err)
			return
		}
		posts = append(posts, post)
	}

	// Check for errors after iterating rows
	if err := rows.Err(); err != nil {
		http.Error(w, "Something went wrong: ", http.StatusInternalServerError)
		fmt.Println("Error iterating rows:", err)
		return
	}
	rows.Close()
	// If no posts are found, set to nil
	if len(posts) == 0 {
		posts = nil
	}
	// for _ , post  := range posts{
	// 	fmt.Println(post.Type)
	// }
	// Render the home page template with posts
	Interface.ExecuteTemplate(w, "home.html", map[string]interface{}{
		"Posts":    posts,
		"LoggedIn": Sesion,
	})
}

func Checksession(w http.ResponseWriter, r *http.Request) (bool, int) {
	userSession, err := r.Cookie("session_id")
	if err != nil {
		return false, http.StatusUnauthorized
	}
	Db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return false, http.StatusInternalServerError
	}
	defer Db.Close()
	username := ""
	Db.QueryRow("SELECT user_id  FROM sessions WHERE session_id = ?", userSession.Value).Scan(&username)
	return username != "", http.StatusOK
}


func Eroors(w http.ResponseWriter, r *http.Request, code int) {
	Interface.ExecuteTemplate(w, "Error.html", map[string]interface{}{
		"Code":    code,
		"Message": http.StatusText(code),
	})
}
