package main

import (
	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func ImportActors(db *sql.DB, path string) {
	db.Exec("DELETE FROM actors")

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening actors CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	_, _ = reader.Read()

	stmt, _ := db.Prepare("INSERT INTO actors (id, first_name, last_name, gender) VALUES (?, ?, ?, ?)")
	defer stmt.Close()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 4 {
			log.Printf("Skipping actor row: %v", err)
			continue
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			continue
		}
		first := strings.Trim(record[1], "\"'")
		last := strings.Trim(record[2], "\"'")
		gender := strings.TrimSpace(record[3])
		if gender != "M" && gender != "F" {
			log.Printf("Invalid gender for actor id=%d: %s", id, gender)
			continue
		}
		_, err = stmt.Exec(id, first, last, gender)
		if err != nil {
			log.Printf("Insert failed for actor id=%d: %v", id, err)
			continue
		}
	}

	log.Println("Imported actors.")
}

func ImportRoles(db *sql.DB, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening roles CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	_, _ = reader.Read()

	stmt, _ := db.Prepare("INSERT INTO roles (actor_id, movie_id, role) VALUES (?, ?, ?)")
	defer stmt.Close()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 3 {
			continue
		}
		actorID, _ := strconv.Atoi(record[0])
		movieID, _ := strconv.Atoi(record[1])
		stmt.Exec(actorID, movieID, record[2])
	}
	log.Println("Imported roles.")
}

func ImportDirectors(db *sql.DB, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening directors CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	_, _ = reader.Read()

	stmt, _ := db.Prepare("INSERT INTO directors (id, first_name, last_name) VALUES (?, ?, ?)")
	defer stmt.Close()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 3 {
			continue
		}
		id, _ := strconv.Atoi(record[0])
		stmt.Exec(id, record[1], record[2])
	}
	log.Println("Imported directors.")
}

func ImportDirectorGenres(db *sql.DB, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening director_genres CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	_, _ = reader.Read()

	stmt, _ := db.Prepare("INSERT INTO director_genres (director_id, genre) VALUES (?, ?)")
	defer stmt.Close()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 3 {
			continue
		}
		dirID, _ := strconv.Atoi(record[0])
		genre := record[1]
		stmt.Exec(dirID, genre)
	}
	log.Println("Imported director genres.")
}

func ImportMovies(db *sql.DB, path string) {
	//db.Exec("DELETE FROM movies")

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening movies CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	_, _ = reader.Read() // skip header

	stmt, err := db.Prepare("INSERT INTO movies (id, title, year, rating) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatalf("Prepare statement failed: %v", err)
	}
	defer stmt.Close()

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Skipping row due to parse error: %v", err)
			continue
		}

		if len(record) < 4 {
			log.Printf("Skipping incomplete row (len=%d): %v", len(record), record)
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("Invalid ID: %v", record[0])
			continue
		}

		title := record[1]
		if title == "" {
			log.Printf("Empty title for movie id=%d", id)
			continue
		}

		year, err := strconv.Atoi(record[2])
		if err != nil || year < 1800 || year > 2100 {
			log.Printf("Invalid year for id=%d: %v", id, record[2])
			continue
		}

		// Handle rating, allow NULL
		var rating sql.NullFloat64
		if record[3] != "NULL" && record[3] != "" {
			r, err := strconv.ParseFloat(record[3], 64)
			if err != nil || r < 0 || r > 10 {
				log.Printf("Invalid rating for id=%d: %v", id, record[3])
				continue
			}
			rating.Valid = true
			rating.Float64 = r
		}

		_, err = stmt.Exec(id, title, year, rating.Float64)
		if err != nil {
			log.Printf("Insert failed for id=%d: %v", id, err)
			continue
		}
		count++
	}
	log.Printf("Imported %d valid movies successfully.", count)
}

func ImportGenres(db *sql.DB, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening genres CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	_, _ = reader.Read() // skip header

	stmt, err := db.Prepare("INSERT INTO genres (movie_id, genre) VALUES (?, ?)")
	if err != nil {
		log.Fatalf("Prepare statement failed: %v", err)
	}
	defer stmt.Close()

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Skipping genre row: %v", err)
			continue
		}

		if len(record) < 2 {
			log.Printf("Skipping incomplete genre row: %v", record)
			continue
		}

		movieID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("Invalid movie ID in genre row: %v", record[0])
			continue
		}

		genre := record[1]
		if genre == "" || genre == "NULL" {
			log.Printf("Empty genre in row: %v", record)
			continue
		}

		_, err = stmt.Exec(movieID, genre)
		if err != nil {
			log.Printf("Insert genre failed: %v", err)
			continue
		}
		count++
	}
	log.Printf("Imported %d valid genres successfully.", count)
}
