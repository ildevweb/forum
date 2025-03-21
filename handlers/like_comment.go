package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"forum/database"
)

func check_comment_like(user_id int, comment_id int) bool {
	rows, err := database.DB.Query(`SELECT 1 FROM comment_likes WHERE user_id = ? AND comment_id = ?`, user_id, comment_id)
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func check_comment_deslike(user_id int, comment_id int) bool {
	rows, err := database.DB.Query(`SELECT 1 FROM comment_deslikes WHERE user_id = ? AND comment_id = ?`, user_id, comment_id)
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func Like_comment_handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Eroors(w, r, http.StatusMethodNotAllowed)
		return
	}
	commentIDStr := r.URL.Path[len("/like-comment/"):]
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil || commentID < 0 {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	userID := get_user_id(r)
	if userID <= 0 {
		fmt.Println("can't like")
		return
	}
	var likes int
	var deslikes int
	rows, _ := database.DB.Query(`SELECT Likes, Deslikes FROM comments WHERE id=?`, commentID)

	for rows.Next() {
		er := rows.Scan(&likes, &deslikes)
		if er != nil {
			fmt.Println(er)
			return
		}
	}

	if check_comment_like(userID, commentID) {
		_, err := database.DB.Exec(`DELETE FROM comment_likes WHERE user_id = ? AND comment_id = ?`, userID, commentID)
		if err != nil {
			Eroors(w, r, 500)
			return
		}
		likes--
	} else {
		if check_comment_deslike(userID, commentID) {
			_, err := database.DB.Exec(`DELETE FROM comment_deslikes WHERE user_id = ? AND comment_id = ?`, userID, commentID)
			if err != nil {
				Eroors(w, r, 500)
				return
			}
			deslikes--
		}

		likes++
		_, err4 := database.DB.Exec(`INSERT INTO comment_likes (user_id, comment_id) VALUES (?,?)`, userID, commentID)
		if err4 != nil {
			Eroors(w, r, 500)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":          "success",
		"newLikeCount":    likes,
		"newDeslikeCount": deslikes,
	}
	json.NewEncoder(w).Encode(response)

	_, err2 := database.DB.Exec(`UPDATE comments SET Likes=?, Deslikes=? WHERE ID=?`, likes, deslikes, commentID)
	if err2 != nil {
		Eroors(w, r, 500)
		return
	}
}


func Deslike_comment_handle(w http.ResponseWriter, r *http.Request) {
	commentIDStr := r.URL.Path[len("/deslike-comment/"):]
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil || commentID < 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID := get_user_id(r)
	if userID <= 0 {
		fmt.Println("can't like")
		return
	}

	var likes int
	var deslikes int
	rows, _ := database.DB.Query(`SELECT Likes, Deslikes FROM comments WHERE id=?`, commentID)

	for rows.Next() {
		er := rows.Scan(&likes, &deslikes)
		if er != nil {
			fmt.Println(er)
			return
		}
	}

	if check_comment_deslike(userID, commentID) {
		_, err := database.DB.Exec(`DELETE FROM comment_deslikes WHERE user_id = ? AND comment_id = ?`, userID, commentID)
		if err != nil {
			Eroors(w, r, 500)
			return
		}
		deslikes--
	} else {
		if check_comment_like(userID, commentID) {
			_, err := database.DB.Exec(`DELETE FROM comment_likes WHERE user_id = ? AND comment_id = ?`, userID, commentID)
			if err != nil {
				Eroors(w, r, 500)
				return
			}

			likes--
		}

		deslikes++

		_, err4 := database.DB.Exec(`INSERT INTO comment_deslikes (user_id, comment_id) VALUES (?,?)`, userID, commentID)
		if err4 != nil {
			Eroors(w, r, 500)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":          "success",
		"newLikeCount":    likes,
		"newDeslikeCount": deslikes,
	}
	json.NewEncoder(w).Encode(response)

	_, err2 := database.DB.Exec(`UPDATE comments SET Likes=?, Deslikes=? WHERE ID=?`, likes, deslikes, commentID)
	if err2 != nil {
		Eroors(w, r, 500)
		return
	}
}
