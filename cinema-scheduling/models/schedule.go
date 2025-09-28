package models

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type Schedule struct {
	ID             int       `json:"id"`
	MovieID        int       `json:"movie_id"`
	HallID         int       `json:"hall_id"` // ✅ link to hall table
	ShowTime       time.Time `json:"show_time"`
	AvailableSeats int       `json:"available_seats"`
	Price          float64   `json:"price"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ---------------- Create Schedule ----------------
func CreateSchedule(s *Schedule) error {
	err := DB.QueryRow(context.Background(),
		`INSERT INTO schedules (movie_id, hall_id, show_time, available_seats, price, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,NOW(),NOW()) RETURNING id`,
		s.MovieID, s.HallID, s.ShowTime, s.AvailableSeats, s.Price,
	).Scan(&s.ID)
	if err != nil {
		log.Printf("❌ CreateSchedule error: %v", err)
		return err
	}
	return nil
}

// ---------------- List Schedules ----------------
func GetAllSchedules() ([]*Schedule, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT id, movie_id, hall_id, show_time, available_seats, price, created_at, updated_at
		 FROM schedules ORDER BY show_time ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*Schedule
	for rows.Next() {
		s := &Schedule{}
		if err := rows.Scan(&s.ID, &s.MovieID, &s.HallID, &s.ShowTime, &s.AvailableSeats, &s.Price, &s.CreatedAt, &s.UpdatedAt); err != nil {
			log.Printf("❌ Scan schedule error: %v", err)
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

// ---------------- Get Schedule By ID ----------------
func GetScheduleByID(id int) (*Schedule, error) {
	var schedule Schedule
	err := DB.QueryRow(context.Background(),
		`SELECT id, movie_id, hall_id, show_time, available_seats, price
         FROM schedules WHERE id=$1`, id,
	).Scan(&schedule.ID, &schedule.MovieID, &schedule.HallID,
		&schedule.ShowTime, &schedule.AvailableSeats, &schedule.Price)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return &schedule, nil
}

// ---------------- Update Schedule ----------------
func UpdateSchedule(s *Schedule) error {
	_, err := DB.Exec(context.Background(),
		`UPDATE schedules SET movie_id=$1, hall_id=$2, show_time=$3, available_seats=$4, price=$5 updated_at=NOW() WHERE id=$5`,
		s.MovieID, s.HallID, s.ShowTime, s.AvailableSeats, s.Price, s.ID)
	if err != nil {
		log.Printf("❌ UpdateSchedule error: %v", err)
		return err
	}
	return nil
}

// ---------------- Delete Schedule ----------------
func DeleteSchedule(id int) error {
	_, err := DB.Exec(context.Background(),
		`DELETE FROM schedules WHERE id=$1`, id)
	if err != nil {
		log.Printf("❌ DeleteSchedule error: %v", err)
		return err
	}
	return nil
}
