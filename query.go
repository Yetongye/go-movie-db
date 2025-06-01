package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// QueryTopGenres retrieves the top genres by average movie rating
func QueryTopGenres(db *sql.DB) {
	rows, err := db.Query(`
		SELECT g.genre, AVG(m.rating) as avg_rating
		FROM genres g
		JOIN movies m ON g.movie_id = m.id
		GROUP BY g.genre
		ORDER BY avg_rating DESC
		LIMIT 10;
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fmt.Println("Top Genres by Average Rating:")
	for rows.Next() {
		var genre string
		var avg float64
		rows.Scan(&genre, &avg)
		fmt.Printf("Genre: %-15s  Avg Rating: %.2f\n", genre, avg)
	}
}

// QueryTopActors retrieves the top actors by number of roles
func QueryTopActors(db *sql.DB) {
	rows, err := db.Query(`
		SELECT a.first_name || ' ' || a.last_name AS actor, COUNT(*) AS roles
		FROM roles r
		JOIN actors a ON r.actor_id = a.id
		GROUP BY actor
		ORDER BY roles DESC
		LIMIT 10;
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fmt.Println("\nTop 10 Actors by Role Count:")
	for rows.Next() {
		var actor string
		var count int
		rows.Scan(&actor, &count)
		fmt.Printf("Actor: %-30s Roles: %d\n", actor, count)
	}
}

// QueryTopDirectors retrieves the top directors by number of movies directed
func QueryDirectorPreferences(db *sql.DB) {
	rows, err := db.Query(`
		SELECT d.first_name || ' ' || d.last_name AS director, dg.genre, COUNT(*) as count
		FROM director_genres dg
		JOIN directors d ON dg.director_id = d.id
		GROUP BY director, dg.genre
		ORDER BY count DESC
		LIMIT 10;
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fmt.Println("\nDirector Genre Preferences:")
	for rows.Next() {
		var director, genre string
		var count int
		rows.Scan(&director, &genre, &count)
		fmt.Printf("Director: %-25s Genre: %-12s Count: %d\n", director, genre, count)
	}
}

// QueryTopMoviesByGenre retrieves the top 5 movies by rating for a given genre
func QueryTopMoviesByGenre(db *sql.DB, genre string) {
	stmt := `
		SELECT m.title, m.rating
		FROM movies m
		JOIN genres g ON m.id = g.movie_id
		WHERE g.genre = ?
		ORDER BY m.rating DESC
		LIMIT 5;
	`
	rows, err := db.Query(stmt, genre)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fmt.Printf("\nTop 5 '%s' Movies by Rating:\n", genre)
	for rows.Next() {
		var title string
		var rating float64
		rows.Scan(&title, &rating)
		fmt.Printf("Title: %-40s Rating: %.2f\n", title, rating)
	}
}

func DebugTableCounts(db *sql.DB) {
	var count int

	db.QueryRow("SELECT COUNT(*) FROM actors").Scan(&count)
	fmt.Printf("Total actors: %d\n", count)

	db.QueryRow("SELECT COUNT(*) FROM roles").Scan(&count)
	fmt.Printf("Total roles: %d\n", count)

	db.QueryRow("SELECT COUNT(*) FROM directors").Scan(&count)
	fmt.Printf("Total directors: %d\n", count)

	db.QueryRow("SELECT COUNT(*) FROM director_genres").Scan(&count)
	fmt.Printf("Total director_genres: %d\n", count)
}

func FindMovieID(db *sql.DB, title string) {
	rows, err := db.Query("SELECT id, title, year FROM movies WHERE title LIKE ?", "%"+title+"%")
	if err != nil {
		log.Printf("Query failed: %v", err)
		return
	}
	defer rows.Close()

	fmt.Println("Search Results:")
	for rows.Next() {
		var id int
		var name string
		var year int
		rows.Scan(&id, &name, &year)
		fmt.Printf("ID: %d | Title: %s (%d)\n", id, name, year)
	}
}

func AddToCollection(db *sql.DB, movieID int, location string, rating float64, note string) {
	stmt, err := db.Prepare("INSERT OR REPLACE INTO my_collection (movie_id, location, my_rating, note) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatalf("Prepare failed: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(movieID, location, rating, note)
	if err != nil {
		log.Printf("Insert failed: %v", err)
	} else {
		fmt.Println("Movie added to collection.")
	}
}

func QueryMyCollection(db *sql.DB) {
	rows, err := db.Query(`
		SELECT m.title, c.location, c.my_rating, c.note
		FROM my_collection c
		JOIN movies m ON c.movie_id = m.id
	`)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return
	}
	defer rows.Close()

	fmt.Println("My Movie Collection:")
	for rows.Next() {
		var title, location, note string
		var rating float64
		rows.Scan(&title, &location, &rating, &note)
		fmt.Printf("Title: %-40s Rating: %.1f | Location: %-12s Note: %s\n", title, rating, location, note)
	}
}

func PromptAddFavorite(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter movie ID to add to collection: ")
	idInput, _ := reader.ReadString('\n')
	idInput = strings.TrimSpace(idInput)
	movieID, err := strconv.Atoi(idInput)
	if err != nil {
		fmt.Println("Invalid ID. Please enter a numeric movie ID.")
		return
	}

	fmt.Print("Enter location: ")
	location, _ := reader.ReadString('\n')
	location = strings.TrimSpace(location)
	if location == "" {
		fmt.Println("Location cannot be empty.")
		return
	}

	fmt.Print("Enter your personal rating (0.0 - 10.0): ")
	ratingInput, _ := reader.ReadString('\n')
	ratingInput = strings.TrimSpace(ratingInput)
	rating, err := strconv.ParseFloat(ratingInput, 64)
	if err != nil || rating < 0 || rating > 10 {
		fmt.Println("Invalid rating. Please enter a number between 0 and 10.")
		return
	}

	fmt.Print("Enter your personal note about the movie: ")
	note, _ := reader.ReadString('\n')
	note = strings.TrimSpace(note)

	AddToCollection(db, movieID, location, rating, note)
}
