package models

import (
	"context"
	"log"
	"time"
)

type Hall struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	Location  *string   `json:"location"` // nullable
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ---------------- Create Hall ----------------
func CreateHall(h *Hall) error {
	err := DB.QueryRow(context.Background(),
		`INSERT INTO halls (name, capacity, location, created_at, updated_at) 
		 VALUES ($1,$2,$3,NOW(),NOW())
         RETURNING id, created_at, updated_at`,
		h.Name, h.Capacity, h.Location).
		Scan(&h.ID, &h.CreatedAt, &h.UpdatedAt)

	if err != nil {
		log.Printf("❌ CreateHall error: %v", err)
		return err
	}
	return nil
}

// ---------------- Get All Halls ----------------
func GetAllHalls() ([]*Hall, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT id, name, capacity, location, created_at, updated_at 
         FROM halls ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var halls []*Hall
	for rows.Next() {
		h := &Hall{}
		if err := rows.Scan(&h.ID, &h.Name, &h.Capacity, &h.Location, &h.CreatedAt, &h.UpdatedAt); err != nil {
			log.Printf("❌ Scan hall error: %v", err)
			return nil, err
		}
		halls = append(halls, h)
	}
	return halls, nil
}

// ---------------- Get Hall By ID ----------------
func GetHallByID(id int) (*Hall, error) {
	h := &Hall{}
	err := DB.QueryRow(context.Background(),
		`SELECT id, name, capacity, location, created_at, updated_at 
         FROM halls WHERE id = $1`, id).
		Scan(&h.ID, &h.Name, &h.Capacity, &h.Location, &h.CreatedAt, &h.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return h, nil
}

// ---------------- Update Hall ----------------
func UpdateHall(h *Hall) error {
	_, err := DB.Exec(context.Background(),
		`UPDATE halls 
         SET name = $1, capacity = $2, updated_at = NOW() 
         WHERE id = $3`,
		h.Name, h.Capacity, &h.Location, h.ID)

	if err != nil {
		log.Printf("❌ UpdateHall error: %v", err)
		return err
	}
	return nil
}

// ---------------- Delete Hall ----------------
func DeleteHall(id int) error {
	_, err := DB.Exec(context.Background(),
		`DELETE FROM halls WHERE id = $1`, id)

	if err != nil {
		log.Printf("❌ DeleteHall error: %v", err)
		return err
	}
	return nil
}
