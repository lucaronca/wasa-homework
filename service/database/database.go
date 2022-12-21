/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Conn() *sql.DB
	Ping() error
}

type appdbimpl struct {
	connectionInstance *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	// Users table
	sqlStmt := `CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY, username TEXT);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	// Users tokens table
	sqlStmt = `
		CREATE TABLE IF NOT EXISTS user_tokens (
			user_id INTEGER NOT NULL,
			token TEXT NOT NULL,
			UNIQUE(user_id, token),
			PRIMARY KEY(user_id, token),
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	// Users bans table
	sqlStmt = `
		CREATE TABLE IF NOT EXISTS user_bans (
			user_id INTEGER NOT NULL,
			banned_id INTEGER NOT NULL,
			UNIQUE(user_id, banned_id),
			PRIMARY KEY(user_id, banned_id),
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(banned_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	// Follows table
	sqlStmt = `
		CREATE TABLE IF NOT EXISTS follows (
			follower_id INTEGER NOT NULL,
			following_id INTEGER NOT NULL,
			UNIQUE(follower_id, following_id),
			PRIMARY KEY (following_id, follower_id),
			FOREIGN KEY(follower_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(following_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	// Photos table
	sqlStmt = `
		CREATE TABLE IF NOT EXISTS photos (
			id INTEGER NOT NULL PRIMARY KEY,
			url TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			upload_date TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	// Likes table
	sqlStmt = `
		CREATE TABLE IF NOT EXISTS likes (
			id INTEGER NOT NULL PRIMARY KEY,
			photo_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			date TEXT NOT NULL,
			UNIQUE(photo_id, user_id),
			FOREIGN KEY(photo_id) REFERENCES photos(id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	// Likes table
	sqlStmt = `
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER NOT NULL PRIMARY KEY,
			photo_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			date TEXT NOT NULL,
			content TEXT NOT NULL,
			FOREIGN KEY(photo_id) REFERENCES photos(id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	return &appdbimpl{
		db,
	}, nil
}

func (db *appdbimpl) Conn() *sql.DB {
	return db.connectionInstance
}

func (db *appdbimpl) Ping() error {
	return db.connectionInstance.Ping()
}
