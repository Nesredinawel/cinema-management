package controllers

import (
	"booking-movie/models"
	"booking-movie/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ---------------- BookingRequest ----------------
type BookingRequest struct {
	UserID      int                   `json:"user_id,omitempty"`
	ScheduleID  int                   `json:"schedule_id" binding:"required"`
	Seats       []string              `json:"seats" binding:"required"`
	Snacks      []models.BookingSnack `json:"snacks,omitempty"`
	TotalAmount float64               `json:"total_amount" binding:"required"`
}

// ---------------- Create Booking ----------------
func CreateBookingHandler(c *gin.Context) {
	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user_id from JWT
	jwtUserIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}
	jwtUserID, ok := jwtUserIDRaw.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id type in token"})
		return
	}

	roleRaw, _ := c.Get("role")
	role, _ := roleRaw.(string)
	if req.UserID == 0 || role == "admin" || role == "staff" {
		req.UserID = jwtUserID
	}

	// ---------------- Start transaction ----------------
	tx, err := models.DB.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start transaction"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(context.Background())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("internal error: %v", r)})
		}
	}()

	// ---------------- Create booking ----------------
	booking := &models.Booking{
		UserID:      req.UserID,
		ScheduleID:  req.ScheduleID,
		TotalAmount: req.TotalAmount,
		Status:      "pending",
	}

	if err := models.CreateBookingTx(tx, booking); err != nil {
		_ = tx.Rollback(context.Background())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create booking"})
		return
	}

	// Generate token for booking
	bookingToken, _ := utils.GenerateToken("booking", booking.ID)

	// ---------------- Insert seats ----------------
	for _, seat := range req.Seats {
		bs := &models.BookingSeat{
			BookingID:  booking.ID,
			SeatNumber: seat,
		}
		if err := models.InsertBookingSeatTx(tx, bs); err != nil {
			_ = tx.Rollback(context.Background())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to insert seat %s", seat)})
			return
		}
		booking.Seats = append(booking.Seats, *bs)
	}

	// ---------------- Insert snacks (optional) ----------------
	var snacksWithToken []gin.H
	for _, snack := range req.Snacks {
		bsnack := &models.BookingSnack{
			BookingID:       booking.ID,
			ScheduleSnackID: snack.ScheduleSnackID,
			Quantity:        snack.Quantity,
			Price:           snack.Price,
		}
		if err := models.InsertBookingSnackTx(tx, bsnack); err != nil {
			_ = tx.Rollback(context.Background())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert booking snack"})
			return
		}
		booking.Snacks = append(booking.Snacks, *bsnack)

		// Generate token for each snack
		token, _ := utils.GenerateToken("booking_snack", bsnack.ID)
		snacksWithToken = append(snacksWithToken, gin.H{
			"booking_snack": bsnack,
			"token":         token,
		})
	}

	// ---------------- Commit transaction ----------------
	if err := tx.Commit(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit booking transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Booking created successfully",
		"booking": gin.H{
			"data":  booking,
			"token": bookingToken,
		},
		"snacks": snacksWithToken,
	})
}

// ---------------- List Bookings ----------------
func ListBookings(c *gin.Context) {
	bookings, err := models.GetAllBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}

	var result []gin.H
	for _, b := range bookings {
		token, _ := utils.GenerateToken("booking", b.ID)
		result = append(result, gin.H{
			"booking": b,
			"token":   token,
		})
	}

	c.JSON(http.StatusOK, gin.H{"bookings": result})
}

// ---------------- Get Booking ----------------
func GetBooking(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking ID"})
		return
	}

	booking, err := models.GetBookingByID(bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch booking"})
		return
	}
	if booking == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	token, _ := utils.GenerateToken("booking", booking.ID)

	c.JSON(http.StatusOK, gin.H{
		"booking": booking,
		"token":   token,
	})
}

// ---------------- Update Booking ----------------
func UpdateBooking(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	existingBooking, err := models.GetBookingByID(bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch booking"})
		return
	}
	if existingBooking == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	var req struct {
		Status      *string  `json:"status"`
		TotalAmount *float64 `json:"total_amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status != nil {
		existingBooking.Status = *req.Status
	}
	if req.TotalAmount != nil {
		existingBooking.TotalAmount = *req.TotalAmount
	}

	if err := models.UpdateBooking(existingBooking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}

	token, _ := utils.GenerateToken("booking", existingBooking.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking updated",
		"booking": existingBooking,
		"token":   token,
	})
}

// ---------------- Delete Booking ----------------
func DeleteBooking(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	if err := models.DeleteBooking(bookingID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted"})
}
