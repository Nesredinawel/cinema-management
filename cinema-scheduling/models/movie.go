package models

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Movie struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Duration       int       `json:"duration"`
	ReleaseYear    int       `json:"release_year"`
	Rating         *float64  `json:"rating"`           // nullable
	ImagePosterURL *string   `json:"image_poster_url"` // nullable
	Genres         []string  `json:"genres,omitempty"` // list of genre names
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

var DB *pgxpool.Pool

// ---------------- Create Movie ----------------
func CreateMovie(movie *Movie, genreIDs []int) error {
	tx, err := DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Insert movie
	err = tx.QueryRow(context.Background(),
		`INSERT INTO movies (title, description, duration, release_year, rating, image_poster_url, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW()) RETURNING id`,
		movie.Title, movie.Description, movie.Duration, movie.ReleaseYear, movie.Rating, movie.ImagePosterURL,
	).Scan(&movie.ID)
	if err != nil {
		log.Printf("❌ CreateMovie error: %v", err)
		return err
	}

	// Insert into movie_genres
	for _, gid := range genreIDs {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO movie_genres (movie_id, genre_id, created_at, updated_at) VALUES ($1,$2,NOW(),NOW())`,
			movie.ID, gid)
		if err != nil {
			log.Printf("❌ Insert movie_genre error: %v", err)
			return err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}

// ---------------- List Movies ----------------
func GetAllMovies() ([]*Movie, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT m.id, m.title, m.description, m.duration, m.release_year, m.rating, m.image_poster_url, m.created_at, m.updated_at,
		        COALESCE(array_agg(g.name) FILTER (WHERE g.id IS NOT NULL), '{}') AS genres
		 FROM movies m
		 LEFT JOIN movie_genres mg ON m.id = mg.movie_id
		 LEFT JOIN genres g ON mg.genre_id = g.id
		 GROUP BY m.id
		 ORDER BY m.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*Movie
	for rows.Next() {
		m := &Movie{}
		var genres []string
		var rating *float64
		if err := rows.Scan(
			&m.ID, &m.Title, &m.Description, &m.Duration, &m.ReleaseYear,
			&rating, &m.ImagePosterURL, &m.CreatedAt, &m.UpdatedAt, &genres,
		); err != nil {
			log.Printf("❌ Scan movie error: %v", err)
			return nil, err
		}
		m.Rating = rating
		m.Genres = genres
		movies = append(movies, m)
	}
	return movies, nil
}

// ---------------- Get Movie By ID ----------------
func GetMovieByID(id int) (*Movie, error) {
	m := &Movie{}
	var genres []string
	var rating *float64

	err := DB.QueryRow(context.Background(),
		`SELECT m.id, m.title, m.description, m.duration, m.release_year, m.rating, m.image_poster_url, m.created_at, m.updated_at,
		        COALESCE(array_agg(g.name) FILTER (WHERE g.id IS NOT NULL), '{}') AS genres
		 FROM movies m
		 LEFT JOIN movie_genres mg ON m.id = mg.movie_id
		 LEFT JOIN genres g ON mg.genre_id = g.id
		 WHERE m.id=$1
		 GROUP BY m.id`, id,
	).Scan(
		&m.ID, &m.Title, &m.Description, &m.Duration, &m.ReleaseYear,
		&rating, &m.ImagePosterURL, &m.CreatedAt, &m.UpdatedAt, &genres,
	)
	if err != nil {
		return nil, err
	}
	m.Rating = rating
	m.Genres = genres
	return m, nil
}

// ---------------- Update Movie ----------------
func UpdateMovie(movie *Movie, genreIDs []int) error {
	tx, err := DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		`UPDATE movies 
		 SET title=$1, description=$2, duration=$3, release_year=$4, rating=$5, image_poster_url=$6, updated_at=NOW() 
		 WHERE id=$7`,
		movie.Title, movie.Description, movie.Duration, movie.ReleaseYear, movie.Rating, movie.ImagePosterURL, movie.ID)
	if err != nil {
		log.Printf("❌ UpdateMovie error: %v", err)
		return err
	}

	// Refresh genres
	_, err = tx.Exec(context.Background(), `DELETE FROM movie_genres WHERE movie_id=$1`, movie.ID)
	if err != nil {
		return err
	}

	for _, gid := range genreIDs {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO movie_genres (movie_id, genre_id, created_at, updated_at) VALUES ($1,$2,NOW(),NOW())`,
			movie.ID, gid)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}

// ---------------- Delete Movie ----------------
func DeleteMovie(id int) error {
	_, err := DB.Exec(context.Background(),
		`DELETE FROM movies WHERE id=$1`, id)
	if err != nil {
		log.Printf("❌ DeleteMovie error: %v", err)
		return err
	}
	return nil
}
