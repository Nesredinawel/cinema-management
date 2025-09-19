package controllers

import (
	"cinema-scheduling/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Add Hall
func AddHall(c *gin.Context) {
	var hall models.Hall
	if err := c.ShouldBindJSON(&hall); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.CreateHall(&hall); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hall"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Hall created", "hall": hall})
}

// List Halls
func ListHalls(c *gin.Context) {
	halls, err := models.GetAllHalls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch halls"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"halls": halls})
}

// Get Hall by ID
func GetHall(c *gin.Context) {
	hallIDStr := c.Param("hall_id")
	hallID, err := strconv.Atoi(hallIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hall ID"})
		return
	}

	hall, err := models.GetHallByID(hallID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hall"})
		return
	}

	if hall == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hall not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hall": hall})
}

// Update Hall
// ---------------- Update Hall ----------------
func UpdateHall(c *gin.Context) {
	hallIDStr := c.Param("hall_id")
	hallID, err := strconv.Atoi(hallIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hall ID"})
		return
	}

	// Fetch existing hall from DB
	existingHall, err := models.GetHallByID(hallID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hall"})
		return
	}
	if existingHall == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hall not found"})
		return
	}

	// Bind incoming JSON into a map
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only fields provided
	if name, ok := req["name"].(string); ok {
		existingHall.Name = name
	}
	if capacity, ok := req["capacity"].(float64); ok { // JSON numbers are float64
		existingHall.Capacity = int(capacity)
	}
	if location, ok := req["location"].(string); ok {
		existingHall.Location = &location
	}

	// Save updated hall
	if err := models.UpdateHall(existingHall); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hall"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hall updated", "hall": existingHall})
}

// Delete Hall
func DeleteHall(c *gin.Context) {
	hallIDStr := c.Param("hall_id")
	hallID, err := strconv.Atoi(hallIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hall ID"})
		return
	}
	if err := models.DeleteHall(hallID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete hall"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Hall deleted"})
}
