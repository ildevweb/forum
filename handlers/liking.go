package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"forum/database"
)

func GetPosts(w http.ResponseWriter, r *http.Request) []Post {
	var posts []Post
	rows, _ := database.DB.Query(`
		SELECT posts.id, posts.title, posts.content, posts.Likes, posts.Deslikes, users.username, posts.created_at
		FROM posts
		JOIN users ON posts.user_id = users.id
	`)
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Likes, &post.Deslikes, &post.Username, &post.CreatedAt)
		if err != nil {
			http.Error(w, "Something went wrong: ", http.StatusInternalServerError)
			return nil
		}
		posts = append(posts, post)
	}
	return posts
}

func get_user_id(r *http.Request) int {
	session, err := r.Cookie("session_id")
	if err != nil {
		return 0
	}
	userID := 0
	err3 := database.DB.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", session.Value).Scan(&userID)
	if err3 != nil {
		return 0
	}

	return userID
}

func Check(user_id int, post_id int) bool {
	rows, err := database.DB.Query(`SELECT 1 FROM likes WHERE user_id = ? AND post_id = ?`, user_id, post_id)
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}


func Like_handle(w http.ResponseWriter, r *http.Request) {
	posts := GetPosts(w, r)
	postIDStr := r.URL.Path[len("/like-post/"):]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID < 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID := get_user_id(r)
	if userID <= 0 {
		fmt.Println("can't like")
		return
	}

	postLikes := posts[postID-1]

	if Check(userID, postID) {
		_, err := database.DB.Exec(`DELETE FROM likes WHERE user_id = ? AND post_id = ?`, userID, postID)
		if err != nil {
			Eroors(w, r, 500)
			return
		}
		postLikes.Likes--
		posts[postID-1].Likes = postLikes.Likes
	} else {
		if Check_deslike(userID, postID) {
			postLikes.Deslikes--
			posts[postID-1].Deslikes = postLikes.Deslikes

			_, err := database.DB.Exec(`DELETE FROM deslikes WHERE user_id = ? AND post_id = ?`, userID, postID)
			if err != nil {
				Eroors(w, r, 500)
				return
			}
		}
		postLikes.Likes++
		posts[postID-1].Likes = postLikes.Likes

		_, err4 := database.DB.Exec(`INSERT INTO likes (user_id, post_id) VALUES (?,?)`, userID, postID)
		if err4 != nil {
			Eroors(w, r, 500)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":          "success",
		"newLikeCount":    postLikes.Likes,
		"newDeslikeCount": postLikes.Deslikes,
	}
	json.NewEncoder(w).Encode(response)

	_, err2 := database.DB.Exec(`UPDATE posts SET Likes=?, Deslikes=? WHERE ID=?`, postLikes.Likes, postLikes.Deslikes, postID)
	if err2 != nil {
		Eroors(w, r, 500)
		return
	}
}

func Check_deslike(user_id int, post_id int) bool {
	rows, err := database.DB.Query(`SELECT 1 FROM deslikes WHERE user_id = ? AND post_id = ?`, user_id, post_id)
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next() 
}

func Deslike_handle(w http.ResponseWriter, r *http.Request) {
	posts := GetPosts(w, r)

	postIDStr := r.URL.Path[len("/deslike-post/"):]
	postID, err := strconv.Atoi(postIDStr)
	post := posts[postID-1]
	if err != nil || postID < 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	userID := get_user_id(r)
	if userID <= 0 {
		fmt.Println("can't like")
		return
	}

	if Check_deslike(userID, postID) {
		_, err := database.DB.Exec(`DELETE FROM deslikes WHERE user_id = ? AND post_id = ?`, userID, postID)
		if err != nil {
			Eroors(w, r, 500)
			return
		}

		post.Deslikes--
		posts[postID-1].Deslikes = post.Deslikes
	} else {
		if Check(userID, postID) {
			post.Likes--
			posts[postID-1].Likes = post.Likes

			_, err := database.DB.Exec(`DELETE FROM likes WHERE user_id = ? AND post_id = ?`, userID, postID)
			if err != nil {
				Eroors(w, r, 500)
				return
			}
		}
		post.Deslikes++
		posts[postID-1].Deslikes = post.Deslikes

		_, err4 := database.DB.Exec(`INSERT INTO deslikes (user_id, post_id) VALUES (?,?)`, userID, postID)
		if err4 != nil {
			Eroors(w, r, 500)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":          "success",
		"newLikeCount":    post.Likes,
		"newDeslikeCount": post.Deslikes,
	}
	json.NewEncoder(w).Encode(response)

	_, err2 := database.DB.Exec(`UPDATE posts SET Deslikes=? , Likes=? WHERE ID=?`, post.Deslikes, post.Likes, postID)
	if err2 != nil {
		Eroors(w, r, 500)
		return
	}
}
