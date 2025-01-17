package main

import (
	"fmt"
	"log"
	"net/http"
	"forum/database"
	"forum/handlers"
)

func main() {
	database.InitDB()
	http.Handle("/scripts/",  http.StripPrefix("/scripts/", http.FileServer(http.Dir("./scripts/"))))
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/profile", handlers.ProfileHandler)
	http.HandleFunc("/make-post", handlers.CreatePostHandler)
	http.HandleFunc("/like-post/", handlers.Like_handle)
	http.HandleFunc("/deslike-post/", handlers.Deslike_handle)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	
	fmt.Println("server starting at: http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("%v", err)
	}
}
