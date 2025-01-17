package handlers

import (
	"net/http"
	"time"

	"forum/database"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	userSession, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	database.DB.Exec("DELETE FROM sessions WHERE session_id = ?", userSession.Value)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
