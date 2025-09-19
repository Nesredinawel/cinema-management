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
	Genre          string    `json:"genre"`
	ReleaseYear    int       `json:"release_year"`
	Rating         *float64  `json:"rating"`           // nullable
	ImagePosterURL *string   `json:"image_poster_url"` // nullable
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

var DB *pgxpool.Pool

// ---------------- Create Movie ----------------
func CreateMovie(movie *Movie) error {
	err := DB.QueryRow(context.Background(),
		`INSERT INTO movies (title, description, duration, genre, release_year, rating, image_poster_url, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW()) RETURNING id`,
		movie.Title, movie.Description, movie.Duration, movie.Genre, movie.ReleaseYear, movie.Rating, movie.ImagePosterURL,
	).Scan(&movie.ID)
	if err != nil {
		log.Printf("❌ CreateMovie error: %v", err)
		return err
	}
	return nil
}

// ---------------- List Movies ----------------
func GetAllMovies() ([]*Movie, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT id, title, description, duration, genre, release_year, rating, image_poster_url, created_at, updated_at
		 FROM movies ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*Movie
	for rows.Next() {
		m := &Movie{}
		var rating *float64 // handle nullable
		if err := rows.Scan(
			&m.ID, &m.Title, &m.Description, &m.Duration, &m.Genre, &m.ReleaseYear,
			&rating, &m.ImagePosterURL, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			log.Printf("❌ Scan movie error: %v", err)
			return nil, err
		}
		m.Rating = rating
		movies = append(movies, m)
	}
	return movies, nil
}

// ---------------- Get Movie By ID ----------------
func GetMovieByID(id int) (*Movie, error) {
	m := &Movie{}
	var rating *float64 // handle nullable
	err := DB.QueryRow(context.Background(),
		`SELECT id, title, description, duration, genre, release_year, rating, image_poster_url, created_at, updated_at
		 FROM movies WHERE id=$1`, id,
	).Scan(
		&m.ID, &m.Title, &m.Description, &m.Duration, &m.Genre, &m.ReleaseYear,
		&rating, &m.ImagePosterURL, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	m.Rating = rating
	return m, nil
}

// ---------------- Update Movie ----------------
func UpdateMovie(movie *Movie) error {
	_, err := DB.Exec(context.Background(),
		`UPDATE movies SET title=$1, description=$2, duration=$3, genre=$4, release_year=$5, rating=$6, image_poster_url=$7, updated_at=NOW() WHERE id=$8`,
		movie.Title, movie.Description, movie.Duration, movie.Genre, movie.ReleaseYear, movie.Rating, movie.ImagePosterURL, movie.ID)
	if err != nil {
		log.Printf("❌ UpdateMovie error: %v", err)
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
