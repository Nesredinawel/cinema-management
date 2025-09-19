package models

import (
	"context"
	"log"
	"time"
)

type ScheduleSnack struct {
	ID         int       `json:"id"`
	ScheduleID int       `json:"schedule_id"`
	SnackID    int       `json:"snack_id"`
	Available  bool      `json:"available"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ---------------- Add ScheduleSnack ----------------
func AddScheduleSnack(ss *ScheduleSnack) error {
	err := DB.QueryRow(context.Background(),
		`INSERT INTO schedule_snacks (schedule_id, snack_id, available, created_at, updated_at)
		 VALUES ($1,$2,$3,NOW(),NOW()) RETURNING id`,
		ss.ScheduleID, ss.SnackID, ss.Available,
	).Scan(&ss.ID)
	if err != nil {
		log.Printf("❌ AddScheduleSnack error: %v", err)
		return err
	}
	return nil
}

// ---------------- List Snacks for a Schedule ----------------
func GetScheduleSnacks(scheduleID int) ([]*ScheduleSnack, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT id, schedule_id, snack_id, available, created_at, updated_at
		 FROM schedule_snacks WHERE schedule_id=$1 ORDER BY created_at DESC`, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scheduleSnacks []*ScheduleSnack
	for rows.Next() {
		ss := &ScheduleSnack{}
		if err := rows.Scan(&ss.ID, &ss.ScheduleID, &ss.SnackID, &ss.Available, &ss.CreatedAt, &ss.UpdatedAt); err != nil {
			log.Printf("❌ Scan schedule_snack error: %v", err)
			return nil, err
		}
		scheduleSnacks = append(scheduleSnacks, ss)
	}
	return scheduleSnacks, nil
}

// ---------------- Get ScheduleSnack by ID ----------------
func GetScheduleSnackByID(id int) (*ScheduleSnack, error) {
	ss := &ScheduleSnack{}
	err := DB.QueryRow(context.Background(),
		`SELECT id, schedule_id, snack_id, available, created_at, updated_at 
		 FROM schedule_snacks WHERE id=$1`, id,
	).Scan(&ss.ID, &ss.ScheduleID, &ss.SnackID, &ss.Available, &ss.CreatedAt, &ss.UpdatedAt)

	if err != nil {
		log.Printf("❌ GetScheduleSnackByID error: %v", err)
		return nil, err
	}

	return ss, nil
}

// ---------------- Update ScheduleSnack ----------------
func UpdateScheduleSnack(ss *ScheduleSnack) error {
	_, err := DB.Exec(context.Background(),
		`UPDATE schedule_snacks SET available=$1, updated_at=NOW() WHERE id=$2`,
		ss.Available, ss.ID)
	if err != nil {
		log.Printf("❌ UpdateScheduleSnack error: %v", err)
		return err
	}
	return nil
}

// ---------------- Delete ScheduleSnack ----------------
func DeleteScheduleSnack(id int) error {
	_, err := DB.Exec(context.Background(),
		`DELETE FROM schedule_snacks WHERE id=$1`, id)
	if err != nil {
		log.Printf("❌ DeleteScheduleSnack error: %v", err)
		return err
	}
	return nil
}
