package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DB is a global variable that holds the database connection
var DB *sql.DB

// InitDB initializes the SQLite database connection and creates necessary tables
// if they don't already exist. It uses a local file "forum.db" as the database.
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createPostsTable := `
	CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
	category TEXT NOT NULL,
    user_id INTEGER NOT NULL,
	Likes INTEGER DEFAULT 0,
	Deslikes INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	const Sessions = `
	CREATE TABLE IF NOT EXISTS sessions (
    session_id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at DATETIME,
    ip_address TEXT,
    user_agent TEXT
);`

	createCommentsTable := `
	CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
	Likes INTEGER DEFAULT 0,
	Deslikes INTEGER DEFAULT 0,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`

	createCommentLikesTable := `
	CREATE TABLE IF NOT EXISTS comment_likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    UNIQUE (user_id, comment_id)
	);`

	createCommentDeslikesTable := `
	CREATE TABLE IF NOT EXISTS comment_deslikes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    UNIQUE (user_id, comment_id)
	);`

	createLikesTable := `
	CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    UNIQUE (user_id, post_id)
	);`

	createDeslikesTable := `
	CREATE TABLE IF NOT EXISTS deslikes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    UNIQUE (user_id, post_id)
	);`

	// Execute table creation statements
	// If any statement fails, the program will terminate
	_, err = DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(createPostsTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(createCommentsTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(createCommentLikesTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(createCommentDeslikesTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(createLikesTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(createDeslikesTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(Sessions)
	if err != nil {
		log.Fatal(err)
	}
}
