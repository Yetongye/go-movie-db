package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // Import SQLite driver
)

func CreateTables(db *sql.DB) {
	createMoviesTable := `
	CREATE TABLE IF NOT EXISTS movies (
	id INTEGER PRIMARY KEY,
	title TEXT NOT NULL,
	year INTEGER,
	rating REAL NULL
	);`

	createGenresTable := `
	CREATE TABLE IF NOT EXISTS genres (
		movie_id INTEGER,
		genre TEXT,
		FOREIGN KEY (movie_id) REFERENCES movies(id)
	);`

	createActorsTable := `
	CREATE TABLE IF NOT EXISTS actors (
		id INTEGER PRIMARY KEY,
		first_name TEXT,
		last_name TEXT,
		gender TEXT
	);`

	createRolesTable := `
	CREATE TABLE IF NOT EXISTS roles (
		actor_id INTEGER,
		movie_id INTEGER,
		role TEXT,
		FOREIGN KEY (actor_id) REFERENCES actors(id),
		FOREIGN KEY (movie_id) REFERENCES movies(id)
	);`

	createDirectorsTable := `
	CREATE TABLE IF NOT EXISTS directors (
		id INTEGER PRIMARY KEY,
		first_name TEXT,
		last_name TEXT
	);`

	createDirectorGenresTable := `
	CREATE TABLE IF NOT EXISTS director_genres (
		director_id INTEGER,
		genre TEXT,
		FOREIGN KEY (director_id) REFERENCES directors(id)
	);`

	createMycollection := `
	CREATE TABLE my_collection (
  movie_id INTEGER PRIMARY KEY,
  location TEXT,
  my_rating REAL,
  note TEXT,
  FOREIGN KEY (movie_id) REFERENCES movies(id)
   );`

	stmts := []string{
		createMoviesTable,
		createGenresTable,
		createActorsTable,
		createRolesTable,
		createDirectorsTable,
		createDirectorGenresTable,
		createMycollection,
	}

	for _, stmt := range stmts {
		_, err := db.Exec(stmt)
		if err != nil {
			log.Fatalf("Error executing statement: %v", err)
		}
	}
}

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
