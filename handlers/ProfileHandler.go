package handlers

import (
	"net/http"

	"forum/database"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	Sesion ,  code  := Checksession(w, r)
	if !Sesion {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}else  if code != http.StatusOK {
		Eroors(w , r , code )
		return 
	}

	Sesionid, _ := r.Cookie("session_id")
	row := database.DB.QueryRow("SELECT username, email, created_at FROM users WHERE id = (SELECT user_id FROM sessions WHERE session_id = ?)", Sesionid.Value)
	if row == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var userinfo struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
	}
	err := row.Scan(&userinfo.Username, &userinfo.Email, &userinfo.CreatedAt)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	posts, Eroor := GetPostsByUser(w, r)
	if Eroor {
		return
	}
	ERR := Interface.ExecuteTemplate(w, "profile.html", map[string]interface{}{
		"Username":  userinfo.Username,
		"Email":     userinfo.Email,
		"CreatedAt": userinfo.CreatedAt,
		"Posts":     posts,
	})
	if  ERR != nil {
		Eroors(w, r, http.StatusInternalServerError)
		return 
	}
	// http.ServeFile(w, r, "templates/profile.html")
}



func GetPostsByUser(w http.ResponseWriter, r *http.Request) ([]Post, bool) {
	var posts []Post
	Sesionid, _ := r.Cookie("session_id")
	var UserID int
	database.DB.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", Sesionid.Value).Scan(&UserID)
	row, er := database.DB.Query(`SELECT posts.id, posts.title, posts.content, posts.category, posts.Likes, posts.Deslikes, users.username, posts.created_at
        FROM posts
        JOIN users ON posts.user_id = users.id
		WHERE posts.user_id=?`, UserID)
	if er != nil {
		Eroors(w, r, http.StatusInternalServerError)
		return nil, true
	}
	for row.Next() {
		var post Post
		err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Likes, &post.Deslikes, &post.Username, &post.CreatedAt)
		if err != nil {
			Eroors(w, r, http.StatusInternalServerError)
			return nil, true
		}
		posts = append(posts, post)
	}
	return posts, false
}

func GetPostsByUserLike(w http.ResponseWriter, r *http.Request) ([]Post, bool) {
	var posts []Post
	Sesionid, _ := r.Cookie("session_id")
	var UserID int
	database.DB.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", Sesionid.Value).Scan(&UserID)

	rows, err := database.DB.Query("SELECT post_id FROM likes WHERE user_id = ?", UserID)
	if err != nil {
		Eroors(w, r, http.StatusInternalServerError)
		return nil, true
	}
	var postIds []int
	for rows.Next() {
		var postId int
		if err := rows.Scan(&postId); err != nil {
			Eroors(w, r, http.StatusInternalServerError)
			return nil, true
		}
		postIds = append(postIds, postId)
	}
	for _, id := range postIds {
		row, er := database.DB.Query(`SELECT posts.id, posts.title, posts.content, posts.category, posts.Likes, posts.Deslikes, users.username, posts.created_at
        FROM posts
        JOIN users ON posts.user_id = users.id
		WHERE posts.id = ?`, id)
		if er != nil {
			Eroors(w, r, http.StatusInternalServerError)
			return nil, true
		}
		for row.Next() {
			var post Post
			err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Likes, &post.Deslikes, &post.Username, &post.CreatedAt)
			if err != nil {
				Eroors(w, r, http.StatusInternalServerError)
				return nil, true
			}
			posts = append(posts, post)
		}
	}
	return posts, false
}
