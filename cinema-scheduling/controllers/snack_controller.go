package controllers

import (
	"cinema-scheduling/models"
	"cinema-scheduling/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ---------------- Add Snack ----------------
func AddSnack(c *gin.Context) {
	var snack models.Snack

	// Handle form values
	snack.Name = c.PostForm("name")
	snack.Description = utils.StrPtr(c.PostForm("description"))
	snack.Category = utils.StrPtr(c.PostForm("category"))

	priceStr := c.PostForm("price")
	if priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			snack.Price = price
		}
	}

	// Handle image: file upload OR URL
	file, _ := c.FormFile("snack_image_url")
	imageURL := c.PostForm("snack_image_url")

	savedImage, err := utils.SaveSnackImage(file, imageURL, "uploads/snacks", c.SaveUploadedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save snack image"})
		return
	}
	snack.SnackImageURL = savedImage

	// Convert to public URL
	snack.SnackImageURL = utils.ConvertSnackImageToPublicURL(snack.SnackImageURL, "http://localhost:8082/")

	// Save to DB
	if err := models.CreateSnack(&snack); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create snack"})
		return
	}

	// ✅ Generate token
	token, _ := utils.GenerateToken("snack", snack.ID)

	c.JSON(http.StatusCreated, gin.H{"message": "Snack created", "snack": snack, "token": token})
}

// ---------------- List Snacks ----------------
func ListSnacks(c *gin.Context) {
	snacks, err := models.GetAllSnacks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch snacks"})
		return
	}

	var result []gin.H
	for _, s := range snacks {
		s.SnackImageURL = utils.ConvertSnackImageToPublicURL(s.SnackImageURL, "http://localhost:8082/")

		// ✅ attach token
		token, _ := utils.GenerateToken("snack", s.ID)
		result = append(result, gin.H{
			"snack": s,
			"token": token,
		})
	}

	c.JSON(http.StatusOK, gin.H{"snacks": result})
}

// ---------------- Get Snack by ID ----------------
func GetSnack(c *gin.Context) {
	snackIDStr := c.Param("snack_id")
	snackID, err := strconv.Atoi(snackIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid snack ID"})
		return
	}

	snack, err := models.GetSnackByID(snackID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch snack"})
		return
	}
	if snack == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Snack not found"})
		return
	}

	snack.SnackImageURL = utils.ConvertSnackImageToPublicURL(snack.SnackImageURL, "http://localhost:8082/")

	// ✅ generate token
	token, _ := utils.GenerateToken("snack", snack.ID)

	c.JSON(http.StatusOK, gin.H{"snack": snack, "token": token})
}

// ---------------- Update Snack ----------------
func UpdateSnack(c *gin.Context) {
	snackIDStr := c.Param("snack_id")
	snackID, err := strconv.Atoi(snackIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid snack ID"})
		return
	}

	existingSnack, err := models.GetSnackByID(snackID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch snack"})
		return
	}
	if existingSnack == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Snack not found"})
		return
	}

	var req struct {
		Name          *string  `json:"name"`
		Description   *string  `json:"description"`
		Category      *string  `json:"category"`
		Price         *float64 `json:"price"`
		SnackImageURL *string  `json:"snack_image_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != nil {
		existingSnack.Name = *req.Name
	}
	if req.Description != nil {
		existingSnack.Description = req.Description
	}
	if req.Category != nil {
		existingSnack.Category = req.Category
	}
	if req.Price != nil {
		existingSnack.Price = *req.Price
	}
	if req.SnackImageURL != nil {
		existingSnack.SnackImageURL = req.SnackImageURL
		if !strings.HasPrefix(*req.SnackImageURL, "http://") && !strings.HasPrefix(*req.SnackImageURL, "https://") {
			publicURL := utils.ConvertPosterToPublicURL(req.SnackImageURL, "http://localhost:8082/")
			existingSnack.SnackImageURL = &publicURL
		}
	}

	if err := models.UpdateSnack(existingSnack); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update snack"})
		return
	}

	// ✅ regenerate token
	token, _ := utils.GenerateToken("snack", existingSnack.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Snack updated", "snack": existingSnack, "token": token})
}

// ---------------- Delete Snack ----------------
func DeleteSnack(c *gin.Context) {
	snackIDStr := c.Param("snack_id")
	snackID, err := strconv.Atoi(snackIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid snack ID"})
		return
	}

	if err := models.DeleteSnack(snackID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete snack"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Snack deleted"})
}
