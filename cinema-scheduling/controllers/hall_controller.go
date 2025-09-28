package controllers

import (
	"cinema-scheduling/models"
	"cinema-scheduling/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ---------------- Add Hall ----------------
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

	// ✅ Generate token after creation
	token, _ := utils.GenerateToken("hall", hall.ID)

	c.JSON(http.StatusCreated, gin.H{"message": "Hall created", "hall": hall, "token": token})
}

// ---------------- List Halls ----------------
func ListHalls(c *gin.Context) {
	halls, err := models.GetAllHalls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch halls"})
		return
	}

	var result []gin.H
	for _, h := range halls {
		// ✅ attach token for each hall
		token, _ := utils.GenerateToken("hall", h.ID)
		result = append(result, gin.H{
			"hall":  h,
			"token": token,
		})
	}

	c.JSON(http.StatusOK, gin.H{"halls": result})
}

// ---------------- Get Hall by ID ----------------
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

	// ✅ generate token
	token, _ := utils.GenerateToken("hall", hall.ID)

	c.JSON(http.StatusOK, gin.H{"hall": hall, "token": token})
}

// ---------------- Update Hall ----------------
func UpdateHall(c *gin.Context) {
	hallIDStr := c.Param("hall_id")
	hallID, err := strconv.Atoi(hallIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hall ID"})
		return
	}

	existingHall, err := models.GetHallByID(hallID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hall"})
		return
	}
	if existingHall == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hall not found"})
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Capacity *int    `json:"capacity"`
		Location *string `json:"location"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != nil {
		existingHall.Name = *req.Name
	}
	if req.Capacity != nil {
		existingHall.Capacity = *req.Capacity
	}
	if req.Location != nil {
		existingHall.Location = req.Location
	}

	if err := models.UpdateHall(existingHall); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hall"})
		return
	}

	// ✅ regenerate token
	token, _ := utils.GenerateToken("hall", existingHall.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Hall updated", "hall": existingHall, "token": token})
}

// ---------------- Delete Hall ----------------
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
