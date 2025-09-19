package controllers

import (
	"cinema-scheduling/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ---------------- Add Schedule ----------------
func AddSchedule(c *gin.Context) {
	var req struct {
		MovieID        int    `json:"movie_id" binding:"required"`
		HallID         int    `json:"hall_id" binding:"required"`
		ShowTime       string `json:"show_time" binding:"required"` // ISO format (RFC3339)
		AvailableSeats int    `json:"available_seats" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	}

	if err := models.CreateSchedule(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Schedule created", "schedule": schedule})
}

// ---------------- List Schedules ----------------
func ListSchedules(c *gin.Context) {
	schedules, err := models.GetAllSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedules": schedules})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule"})
		return
	}

	if schedule == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedule": schedule})
}

// ---------------- Update Schedule ----------------
func UpdateSchedule(c *gin.Context) {
	scheduleIDStr := c.Param("schedule_id")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	// Fetch existing schedule from DB
	existingSchedule, err := models.GetScheduleByID(scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule"})
		return
	}
	if existingSchedule == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	// Bind incoming JSON into a map
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only provided fields
	if movieID, ok := req["movie_id"].(float64); ok { // JSON numbers are float64
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

	// Save updated schedule
	if err := models.UpdateSchedule(existingSchedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule updated", "schedule": existingSchedule})
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
