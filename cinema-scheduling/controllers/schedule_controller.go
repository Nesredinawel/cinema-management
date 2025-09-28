package controllers

import (
	"cinema-scheduling/models"
	"cinema-scheduling/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// ---------------- Add Schedule ----------------
func AddSchedule(c *gin.Context) {
	var req struct {
		MovieID        int     `json:"movie_id" binding:"required"`
		MovieToken     string  `json:"movie_token" binding:"required"` // ✅ require movie token
		HallID         int     `json:"hall_id" binding:"required"`
		HallToken      string  `json:"hall_token" binding:"required"` // ✅ require hall token
		ShowTime       string  `json:"show_time" binding:"required"`
		AvailableSeats int     `json:"available_seats" binding:"required"`
		Price          float64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Verify movie token
	movieEntity, movieID, err := utils.VerifyToken(req.MovieToken)
	if err != nil || movieEntity != "movie" || movieID != req.MovieID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or mismatched movie token"})
		return
	}

	// ✅ Verify hall token
	hallEntity, hallID, err := utils.VerifyToken(req.HallToken)
	if err != nil || hallEntity != "hall" || hallID != req.HallID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or mismatched hall token"})
		return
	}

	// Parse show_time
	showTime, err := time.Parse(time.RFC3339, req.ShowTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid show_time format. Use RFC3339 (e.g. 2025-09-19T15:04:05Z)"})
		return
	}

	schedule := models.Schedule{
		MovieID:        req.MovieID,
		HallID:         req.HallID,
		ShowTime:       showTime,
		AvailableSeats: req.AvailableSeats,
		Price:          req.Price,
	}

	if err := models.CreateSchedule(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule"})
		return
	}

	// ✅ Generate schedule token
	token, err := utils.GenerateToken("schedule", schedule.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate schedule token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Schedule created",
		"schedule": schedule,
		"token":    token,
	})
}

// ---------------- List Schedules ----------------
func ListSchedules(c *gin.Context) {
	schedules, err := models.GetAllSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedules"})
		return
	}

	var result []gin.H
	for _, s := range schedules {
		token, _ := utils.GenerateToken("schedule", s.ID)
		result = append(result, gin.H{
			"schedule": s,
			"token":    token,
		})
	}

	c.JSON(http.StatusOK, gin.H{"schedules": result})
}

// ---------------- Get Schedule by ID ----------------
func GetSchedule(c *gin.Context) {
	scheduleIDStr := c.Param("schedule_id")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	schedule, err := models.GetScheduleByID(scheduleID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule"})
		return
	}

	token, _ := utils.GenerateToken("schedule", schedule.ID)

	c.JSON(http.StatusOK, gin.H{
		"schedule": schedule,
		"token":    token,
	})
}

// ---------------- Update Schedule ----------------
func UpdateSchedule(c *gin.Context) {
	scheduleIDStr := c.Param("schedule_id")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	existingSchedule, err := models.GetScheduleByID(scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule"})
		return
	}
	if existingSchedule == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if movieID, ok := req["movie_id"].(float64); ok {
		existingSchedule.MovieID = int(movieID)
	}
	if hallID, ok := req["hall_id"].(float64); ok {
		existingSchedule.HallID = int(hallID)
	}
	if showTimeStr, ok := req["show_time"].(string); ok && showTimeStr != "" {
		showTime, err := time.Parse(time.RFC3339, showTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid show_time format. Use RFC3339"})
			return
		}
		existingSchedule.ShowTime = showTime
	}
	if availableSeats, ok := req["available_seats"].(float64); ok {
		existingSchedule.AvailableSeats = int(availableSeats)
	}
	if price, ok := req["price"].(float64); ok {
		existingSchedule.Price = price
	}

	if err := models.UpdateSchedule(existingSchedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule"})
		return
	}

	token, _ := utils.GenerateToken("schedule", existingSchedule.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Schedule updated",
		"schedule": existingSchedule,
		"token":    token,
	})
}

// ---------------- Delete Schedule ----------------
func DeleteSchedule(c *gin.Context) {
	scheduleIDStr := c.Param("schedule_id")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	if err := models.DeleteSchedule(scheduleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted"})
}
