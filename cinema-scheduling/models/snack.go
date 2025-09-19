package models

import (
	"context"
	"log"
	"time"
)

type Snack struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	Description   *string   `json:"description"`     // nullable
	Category      *string   `json:"category"`        // nullable
	SnackImageURL *string   `json:"snack_image_url"` // nullable
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ---------------- Add Snack ----------------
func CreateSnack(snack *Snack) error {
	err := DB.QueryRow(context.Background(),
		`INSERT INTO snacks (name, price, description, category, snack_image_url, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,NOW(),NOW()) RETURNING id`,
		snack.Name, snack.Price, snack.Description, snack.Category, snack.SnackImageURL,
	).Scan(&snack.ID)
	if err != nil {
		log.Printf("❌ CreateSnack error: %v", err)
		return err
	}
	return nil
}

// ---------------- List Snacks ----------------
func GetAllSnacks() ([]*Snack, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT id, name, price, description, category, snack_image_url, created_at, updated_at 
		 FROM snacks ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snacks []*Snack
	for rows.Next() {
		s := &Snack{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Price, &s.Description, &s.Category, &s.SnackImageURL, &s.CreatedAt, &s.UpdatedAt); err != nil {
			log.Printf("❌ Scan snack error: %v", err)
			return nil, err
		}
		snacks = append(snacks, s)
	}
	return snacks, nil
}

// ---------------- Get Snack By ID ----------------
func GetSnackByID(id int) (*Snack, error) {
	s := &Snack{}
	err := DB.QueryRow(context.Background(),
		`SELECT id, name, price, description, category, snack_image_url, created_at, updated_at 
		 FROM snacks WHERE id=$1`, id,
	).Scan(&s.ID, &s.Name, &s.Price, &s.Description, &s.Category, &s.SnackImageURL, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// ---------------- Update Snack ----------------
func UpdateSnack(snack *Snack) error {
	_, err := DB.Exec(context.Background(),
		`UPDATE snacks SET name=$1, price=$2, description=$3, category=$4, snack_image_url=$5, updated_at=NOW() WHERE id=$6`,
		snack.Name, snack.Price, snack.Description, snack.Category, snack.SnackImageURL, snack.ID)
	if err != nil {
		log.Printf("❌ UpdateSnack error: %v", err)
		return err
	}
	return nil
}

// ---------------- Delete Snack ----------------
func DeleteSnack(id int) error {
	_, err := DB.Exec(context.Background(),
		`DELETE FROM snacks WHERE id=$1`, id)
	if err != nil {
		log.Printf("❌ DeleteSnack error: %v", err)
		return err
	}
	return nil
}
