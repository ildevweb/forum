package handlers

import (
	"fmt"
	"net/http"

	"forum/database"
)
type Comment struct {
	CommentID         int
	Content    string
	Likes int
	Deslikes int
	user_id    int
	Post_id    int
	CreatedAt  string
	Created_by string
}

func GetComments(w http.ResponseWriter, r *http.Request, postid int) ([]Comment, error) {
	var Comments []Comment
	rows, err := database.DB.Query("SELECT * FROM comments WHERE post_id = ?", postid)
	if err != nil {
		Eroors(w, r, http.StatusInternalServerError)
		return Comments, err
	}
	for rows.Next() {
		var Comment Comment
		err := rows.Scan(&Comment.CommentID, &Comment.Content, &Comment.Likes, &Comment.Deslikes, &Comment.user_id, &Comment.Post_id, &Comment.CreatedAt)
		if err != nil {
			Eroors(w, r, http.StatusInternalServerError)
			fmt.Println("Error fetching comments:", err)
			return Comments, err
		}
		Comment.Created_by, err = Getusernamebyid(Comment.user_id)
		if err != nil {
			Eroors(w, r, http.StatusInternalServerError)
			fmt.Println("Error fetching username:", err)
			return Comments, err
		}
		Comments = append(Comments, Comment)

	}
	if  len(Comments) == 0 {
		return nil, nil
	}
	return Comments, nil
}

func Getusernamebyid(id int) (string, error) {
	var username string
	err := database.DB.QueryRow("SELECT username FROM users WHERE id = ?", id).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
