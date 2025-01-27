package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"forum/database"
)

type GivenComment struct {
	Post_id string `json:"post_id"`
	Comment string `json:"comment"`
}

type Comments struct {
	ID        int
	Content   string
	Likes     int
	Deslikes  int
	Username  string
	PostID    int
	CreatedAt string
}

func SaveComments(w http.ResponseWriter, r *http.Request) {
	sesion, _ := Checksession(w, r)
	if !sesion {
		http.Error(w, "<a href='/register'>Register</a> first Or Login <a href='/login'>Login</a>", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		Eroors(w, r, http.StatusMethodNotAllowed)
		return
	}
	var comment GivenComment


	/*err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}*/


	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}


	if strings.Contains(string(body), "{") || r.Header.Get("Content-Type") == "application/json" {
		err := json.Unmarshal(body, &comment)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}
	} else {
		bodySplit := strings.Split(string(body), "&")
		if len(bodySplit) < 2 {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		postID := bodySplit[0][len("post_id="):]
		commentText := bodySplit[1][len("comment="):]

		comment = GivenComment{
			Post_id:  postID,
			Comment: commentText,
		}
	}


	/*_, err2 := json.Marshal(comment)
	if err2 != nil {
		http.Error(w, "Error converting struct to JSON", http.StatusInternalServerError)
		return
	}*/

	
	if comment.Comment == "" {
		http.Error(w, "Comment is required", http.StatusBadRequest)
		return
	}

	UserId := get_user_id(r)
	// var postId string
	// database.DB.QueryRow("SELECT id FROM posts WHERE content = ?", comment).Scan(&postId)

	_, err = database.DB.Exec(
		`INSERT INTO comments (Post_id, User_id, Content) VALUES (?, ?, ?)`,
		comment.Post_id, UserId, comment.Comment,
	)
	if err != nil {
		http.Error(w, "Failed to add  your  Comment ! ", http.StatusInternalServerError)
		fmt.Println("Error adding comment:", err)
		return
	}

	// http.Redirect(w, r, "/home", http.StatusSeeOther)
	comments, _ := GetCommentsByUser(w, r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":  "success",
		"comment": comments[len(comments)-1],
	}
	json.NewEncoder(w).Encode(response)
}




func GetCommentsByUser(w http.ResponseWriter, r *http.Request) ([]Comments, bool) {
	var comments []Comments
	Sesionid, _ := r.Cookie("session_id")
	var UserID int
	database.DB.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", Sesionid.Value).Scan(&UserID)
	row, er := database.DB.Query(`SELECT comments.id, comments.content, comments.Likes, comments.Deslikes, users.username, comments.post_id, comments.created_at
        FROM comments
        JOIN users ON comments.user_id = users.id
		WHERE comments.user_id=?`, UserID)
	if er != nil {
		Eroors(w, r, http.StatusInternalServerError)
		return nil, true
	}
	for row.Next() {
		var comment Comments
		err := row.Scan(&comment.ID, &comment.Content, &comment.Likes, &comment.Deslikes, &comment.Username, &comment.PostID, &comment.CreatedAt)
		if err != nil {
			Eroors(w, r, http.StatusInternalServerError)
			return nil, true
		}
		comments = append(comments, comment)
	}
	return comments, false
}
