package models

import (
	"context"
	"log"
	"time"
)

// Genre struct
type Genre struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ---------------- Create Genre ----------------
func CreateGenre(genre *Genre) error {
	err := DB.QueryRow(context.Background(),
		`INSERT INTO genres (name, created_at, updated_at)
		 VALUES ($1, NOW(), NOW()) RETURNING id, created_at, updated_at`,
		genre.Name,
	).Scan(&genre.ID, &genre.CreatedAt, &genre.UpdatedAt)

	if err != nil {
		log.Printf("❌ CreateGenre error: %v", err)
		return err
	}
	return nil
}

// ---------------- List Genres ----------------
func GetAllGenres() ([]*Genre, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT id, name, created_at, updated_at FROM genres ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []*Genre
	for rows.Next() {
		g := &Genre{}
		if err := rows.Scan(&g.ID, &g.Name, &g.CreatedAt, &g.UpdatedAt); err != nil {
			log.Printf("❌ Scan genre error: %v", err)
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, nil
}

// ---------------- Get Genre By ID ----------------
func GetGenreByID(id int) (*Genre, error) {
	g := &Genre{}
	err := DB.QueryRow(context.Background(),
		`SELECT id, name, created_at, updated_at FROM genres WHERE id=$1`, id,
	).Scan(&g.ID, &g.Name, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// ---------------- Get Genre Names By IDs ----------------
func GetGenreNamesByIDs(ids []int) ([]string, error) {
	if len(ids) == 0 {
		return []string{}, nil
	}

	query := `SELECT name FROM genres WHERE id = ANY($1)`
	rows, err := DB.Query(context.Background(), query, ids)
	if err != nil {
		log.Printf("❌ GetGenreNamesByIDs error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Printf("❌ Scan genre name error: %v", err)
			return nil, err
		}
		names = append(names, name)
	}

	return names, nil
}

// ---------------- Update Genre ----------------
func UpdateGenre(genre *Genre) error {
	_, err := DB.Exec(context.Background(),
		`UPDATE genres SET name=$1, updated_at=NOW() WHERE id=$2`,
		genre.Name, genre.ID,
	)
	if err != nil {
		log.Printf("❌ UpdateGenre error: %v", err)
		return err
	}
	return nil
}

// ---------------- Delete Genre ----------------
func DeleteGenre(id int) error {
	_, err := DB.Exec(context.Background(),
		`DELETE FROM genres WHERE id=$1`, id,
	)
	if err != nil {
		log.Printf("❌ DeleteGenre error: %v", err)
		return err
	}
	return nil
}
