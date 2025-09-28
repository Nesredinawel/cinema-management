package controllers

import (
	"cinema-scheduling/models"
	"cinema-scheduling/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ---------------- Add Schedule Snack ----------------
func AddScheduleSnack(c *gin.Context) {
	var req struct {
		ScheduleID    int    `json:"schedule_id" binding:"required"`
		SnackID       int    `json:"snack_id" binding:"required"`
		Available     bool   `json:"available"`
		ScheduleToken string `json:"schedule_token" binding:"required"`
		SnackToken    string `json:"snack_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify schedule token
	sType, sID, err := utils.VerifyToken(req.ScheduleToken)
	if err != nil || sType != "schedule" || sID != req.ScheduleID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or mismatched schedule token"})
		return
	}

	// Verify snack token
	snType, snID, err := utils.VerifyToken(req.SnackToken)
	if err != nil || snType != "snack" || snID != req.SnackID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or mismatched snack token"})
		return
	}

	ss := &models.ScheduleSnack{
		ScheduleID: req.ScheduleID,
		SnackID:    req.SnackID,
		Available:  req.Available,
	}

	if err := models.AddScheduleSnack(ss); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add snack to schedule"})
		return
	}

	// Generate token for schedule snack
	token, _ := utils.GenerateToken("schedule_snack", ss.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":        "Snack added to schedule",
		"schedule_snack": ss,
		"token":          token,
	})
}

// ---------------- List Schedule Snacks ----------------
func ListScheduleSnacks(c *gin.Context) {
	scheduleIDStr := c.Param("schedule_id")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	snacks, err := models.GetScheduleSnacks(scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule snacks"})
		return
	}

	// Attach token to each schedule snack
	var result []gin.H
	for _, ss := range snacks {
		token, _ := utils.GenerateToken("schedule_snack", ss.ID)
		result = append(result, gin.H{
			"schedule_snack": ss,
			"token":          token,
		})
	}

	c.JSON(http.StatusOK, gin.H{"schedule_snacks": result})
}

// ---------------- Get Schedule Snack for a specific schedule ----------------
func GetScheduleSnack(c *gin.Context) {
	scheduleIDStr := c.Param("schedule_id")
	snackIDStr := c.Param("snack_id")

	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	snackID, err := strconv.Atoi(snackIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid snack ID"})
		return
	}

	// Fetch the schedule snack by schedule_id and snack_id
	scheduleSnack, err := models.GetScheduleSnackByScheduleAndSnack(scheduleID, snackID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule snack"})
		return
	}
	if scheduleSnack == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule snack not found for this schedule"})
		return
	}

	// Generate token
	token, _ := utils.GenerateToken("schedule_snack", scheduleSnack.ID)

	c.JSON(http.StatusOK, gin.H{
		"schedule_snack": scheduleSnack,
		"token":          token,
	})
}

// ---------------- Update Schedule Snack ----------------
func UpdateScheduleSnack(c *gin.Context) {
	ssIDStr := c.Param("schedule_snack_id")
	ssID, err := strconv.Atoi(ssIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule snack ID"})
		return
	}

	existingSS, err := models.GetScheduleSnackByID(ssID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule snack"})
		return
	}
	if existingSS == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule snack not found"})
		return
	}

	var req struct {
		Available *bool `json:"available"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Available != nil {
		existingSS.Available = *req.Available
	}

	if err := models.UpdateScheduleSnack(existingSS); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule snack"})
		return
	}

	// Generate token
	token, _ := utils.GenerateToken("schedule_snack", existingSS.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":        "Schedule snack updated",
		"schedule_snack": existingSS,
		"token":          token,
	})
}

// ---------------- Delete Schedule Snack ----------------
func DeleteScheduleSnack(c *gin.Context) {
	ssIDStr := c.Param("schedule_snack_id")
	ssID, err := strconv.Atoi(ssIDStr)
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
