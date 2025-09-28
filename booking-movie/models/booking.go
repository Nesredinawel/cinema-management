package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// ---------------- Booking Structs ----------------
type Booking struct {
	ID               int            `json:"id"`
	UserID           int            `json:"user_id"`
	ScheduleID       int            `json:"schedule_id"`
	TotalAmount      float64        `json:"total_amount"`
	Status           string         `json:"status"`
	PaymentReference *string        `json:"payment_reference,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	Seats            []BookingSeat  `json:"seats,omitempty"`
	Snacks           []BookingSnack `json:"snacks,omitempty"`
}

type BookingSeat struct {
	ID         int       `json:"id"`
	BookingID  int       `json:"booking_id"`
	SeatNumber string    `json:"seat_number"` // matches controller
	CreatedAt  time.Time `json:"created_at"`
}

type BookingSnack struct {
	ID              int       `json:"id"`
	BookingID       int       `json:"booking_id"`
	ScheduleSnackID int       `json:"schedule_snack_id"`
	Quantity        int       `json:"quantity"`
	Price           float64   `json:"price"`
	CreatedAt       time.Time `json:"created_at"`
}

// ---------------- Create Booking ----------------
func CreateBookingTx(tx pgx.Tx, b *Booking) error {
	return tx.QueryRow(context.Background(),
		`INSERT INTO bookings (user_id, schedule_id, total_amount, status, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,NOW(),NOW()) RETURNING id, created_at, updated_at`,
		b.UserID, b.ScheduleID, b.TotalAmount, b.Status,
	).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

// ---------------- Insert Booking Seat ----------------
func InsertBookingSeatTx(tx pgx.Tx, bs *BookingSeat) error {
	return tx.QueryRow(context.Background(),
		`INSERT INTO booking_seats (booking_id, seat_number, created_at)
		 VALUES ($1,$2,NOW()) RETURNING id`,
		bs.BookingID, bs.SeatNumber,
	).Scan(&bs.ID)
}

// ---------------- Insert Booking Snack ----------------
func InsertBookingSnackTx(tx pgx.Tx, bs *BookingSnack) error {
	return tx.QueryRow(context.Background(),
		`INSERT INTO booking_snacks (booking_id, schedule_snack_id, quantity, price, created_at)
		 VALUES ($1,$2,$3,$4,NOW()) RETURNING id`,
		bs.BookingID, bs.ScheduleSnackID, bs.Quantity, bs.Price,
	).Scan(&bs.ID)
}

// ---------------- Get Booking by ID ----------------
func GetBookingByID(id int) (*Booking, error) {
	ctx := context.Background()
	b := &Booking{}
	err := DB.QueryRow(ctx,
		`SELECT id, user_id, schedule_id, total_amount, status, payment_reference, created_at, updated_at
		 FROM bookings WHERE id=$1`, id,
	).Scan(&b.ID, &b.UserID, &b.ScheduleID, &b.TotalAmount, &b.Status, &b.PaymentReference, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Load seats
	rows, _ := DB.Query(ctx,
		`SELECT id, booking_id, seat_number, created_at FROM booking_seats WHERE booking_id=$1`, b.ID)
	defer rows.Close()
	for rows.Next() {
		var s BookingSeat
		if err := rows.Scan(&s.ID, &s.BookingID, &s.SeatNumber, &s.CreatedAt); err == nil {
			b.Seats = append(b.Seats, s)
		}
	}

	// Load snacks
	rows2, _ := DB.Query(ctx,
		`SELECT id, booking_id, schedule_snack_id, quantity, price, created_at FROM booking_snacks WHERE booking_id=$1`, b.ID)
	defer rows2.Close()
	for rows2.Next() {
		var s BookingSnack
		if err := rows2.Scan(&s.ID, &s.BookingID, &s.ScheduleSnackID, &s.Quantity, &s.Price, &s.CreatedAt); err == nil {
			b.Snacks = append(b.Snacks, s)
		}
	}

	return b, nil
}

// ---------------- Get All Bookings ----------------
func GetAllBookings() ([]Booking, error) {
	ctx := context.Background()
	rows, err := DB.Query(ctx,
		`SELECT id, user_id, schedule_id, total_amount, status, payment_reference, created_at, updated_at
		 FROM bookings ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.ScheduleID, &b.TotalAmount, &b.Status, &b.PaymentReference, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		// Load seats
		seatRows, _ := DB.Query(ctx, `SELECT id, booking_id, seat_number, created_at FROM booking_seats WHERE booking_id=$1`, b.ID)
		for seatRows.Next() {
			var s BookingSeat
			seatRows.Scan(&s.ID, &s.BookingID, &s.SeatNumber, &s.CreatedAt)
			b.Seats = append(b.Seats, s)
		}
		seatRows.Close()
		// Load snacks
		snackRows, _ := DB.Query(ctx, `SELECT id, booking_id, schedule_snack_id, quantity, price, created_at FROM booking_snacks WHERE booking_id=$1`, b.ID)
		for snackRows.Next() {
			var s BookingSnack
			snackRows.Scan(&s.ID, &s.BookingID, &s.ScheduleSnackID, &s.Quantity, &s.Price, &s.CreatedAt)
			b.Snacks = append(b.Snacks, s)
		}
		snackRows.Close()

		bookings = append(bookings, b)
	}

	return bookings, nil
}

// ---------------- Update Booking ----------------
func UpdateBooking(b *Booking) error {
	_, err := DB.Exec(context.Background(),
		`UPDATE bookings SET total_amount=$1, status=$2, updated_at=NOW() WHERE id=$3`,
		b.TotalAmount, b.Status, b.ID,
	)
	return err
}

// ---------------- Delete Booking ----------------
func DeleteBooking(id int) error {
	_, err := DB.Exec(context.Background(),
		`DELETE FROM bookings WHERE id=$1`, id)
	return err
}
