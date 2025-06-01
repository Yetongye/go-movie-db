-- movies table
CREATE TABLE IF NOT EXISTS movies (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  year INTEGER,
  rating REAL
);

-- genres table
CREATE TABLE IF NOT EXISTS genres (
  movie_id INTEGER,
  genre TEXT,
  FOREIGN KEY (movie_id) REFERENCES movies(id)
);

-- actors table 
CREATE TABLE IF NOT EXISTS actors (
  id INTEGER PRIMARY KEY,
  first_name TEXT,
  last_name TEXT,
  gender TEXT
);

-- roles table
CREATE TABLE IF NOT EXISTS roles (
  actor_id INTEGER,
  movie_id INTEGER,
  role TEXT,
  FOREIGN KEY (actor_id) REFERENCES actors(id),
  FOREIGN KEY (movie_id) REFERENCES movies(id)
);

-- directors table
CREATE TABLE IF NOT EXISTS directors (
  id INTEGER PRIMARY KEY,
  first_name TEXT,
  last_name TEXT
);

-- director_genres table
CREATE TABLE IF NOT EXISTS director_genres (
  director_id INTEGER,
  genre TEXT,
  FOREIGN KEY (director_id) REFERENCES directors(id)
);

-- my_collection table
CREATE TABLE IF NOT EXISTS my_collection (
  movie_id INTEGER PRIMARY KEY,
  location TEXT,
  my_rating REAL,
  note TEXT,
  FOREIGN KEY (movie_id) REFERENCES movies(id)
);

-- sqlite3 movies.db < db/schema.sql
