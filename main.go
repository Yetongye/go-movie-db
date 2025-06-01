package main

//import "os"

func main() {
	//os.Remove("movies.db") // Remove existing database file for a fresh start

	db := InitDB("movies.db")
	defer db.Close()

	//CreateTables(db)

	//ImportMovies(db, "data/IMDB-movies.csv")
	//ImportGenres(db, "data/IMDB-movies_genres.csv")
	//ImportActors(db, "data/IMDB-actors.csv")
	//ImportRoles(db, "data/IMDB-roles.csv")
	//ImportDirectors(db, "data/IMDB-directors.csv")
	//ImportDirectorGenres(db, "data/IMDB-directors_genres.csv")

	QueryTopGenres(db)
	QueryTopActors(db)
	QueryDirectorPreferences(db)
	QueryTopMoviesByGenre(db, "Family")
	DebugTableCounts(db)

	FindMovieID(db, "Matrix")
	AddToCollection(db, 41, "somewhere", 9.5, "one of the best sci-fi movies ever")

	PromptAddFavorite(db)
	QueryMyCollection(db)

}
