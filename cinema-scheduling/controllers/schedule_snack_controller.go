package controllers

import (
	"cinema-scheduling/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ---------------- Add Snack to Schedule ----------------
func AddScheduleSnack(c *gin.Context) {
	scheduleID, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	var req struct {
		SnackID   int  `json:"snack_id" binding:"required"`
		Available bool `json:"available"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ss := &models.ScheduleSnack{
		ScheduleID: scheduleID,
		SnackID:    req.SnackID,
		Available:  req.Available,
	}

	if err := models.AddScheduleSnack(ss); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add snack to schedule"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Snack added to schedule", "schedule_snack": ss})
}

// ---------------- List Snacks for Schedule ----------------
func ListScheduleSnacks(c *gin.Context) {
	scheduleID, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	snacks, err := models.GetScheduleSnacks(scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule snacks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedule_snacks": snacks})
}

// ---------------- Update Schedule Snack ----------------
func UpdateScheduleSnack(c *gin.Context) {
	ssID, err := strconv.Atoi(c.Param("schedule_snack_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule snack ID"})
		return
	}

	// Fetch existing schedule snack
	existingSS, err := models.GetScheduleSnackByID(ssID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule snack"})
		return
	}
	if existingSS == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule snack not found"})
		return
	}

	// Bind incoming JSON as map for partial updates
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only provided fields
	if available, ok := req["available"].(bool); ok {
		existingSS.Available = available
	}
	if snackID, ok := req["snack_id"].(float64); ok { // JSON numbers come as float64
		existingSS.SnackID = int(snackID)
	}
	if scheduleID, ok := req["schedule_id"].(float64); ok {
		existingSS.ScheduleID = int(scheduleID)
	}

	// Save updated schedule snack
	if err := models.UpdateScheduleSnack(existingSS); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule snack"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule snack updated", "schedule_snack": existingSS})
}

// ---------------- Delete Schedule Snack ----------------
func DeleteScheduleSnack(c *gin.Context) {
	ssID, err := strconv.Atoi(c.Param("schedule_snack_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule snack ID"})
		return
	}

	if err := models.DeleteScheduleSnack(ssID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete schedule snack"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule snack deleted"})
}
