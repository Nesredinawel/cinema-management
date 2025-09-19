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
		price, err := strconv.ParseFloat(priceStr, 64)
		if err == nil {
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

	c.JSON(http.StatusCreated, gin.H{"message": "Snack created", "snack": snack})
}

// ---------------- List Snacks ----------------
func ListSnacks(c *gin.Context) {
	snacks, err := models.GetAllSnacks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch snacks"})
		return
	}

	for _, s := range snacks {
		s.SnackImageURL = utils.ConvertSnackImageToPublicURL(s.SnackImageURL, "http://localhost:8082/")
	}

	c.JSON(http.StatusOK, gin.H{"snacks": snacks})
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

	c.JSON(http.StatusOK, gin.H{"snack": snack})
}

// ---------------- Update Snack ----------------
func UpdateSnack(c *gin.Context) {
	snackIDStr := c.Param("snack_id")
	snackID, err := strconv.Atoi(snackIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid snack ID"})
		return
	}

	// Fetch existing snack from DB
	existingSnack, err := models.GetSnackByID(snackID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch snack"})
		return
	}
	if existingSnack == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Snack not found"})
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
		existingSnack.Name = name
	}
	if desc, ok := req["description"].(string); ok {
		existingSnack.Description = utils.StrPtr(desc)
	}
	if cat, ok := req["category"].(string); ok {
		existingSnack.Category = utils.StrPtr(cat)
	}
	if price, ok := req["price"].(float64); ok {
		existingSnack.Price = price
	}
	if poster, ok := req["image_poster_url"].(string); ok {
		existingSnack.SnackImageURL = &poster

		// convert to public URL if needed
		if !strings.HasPrefix(poster, "http://") && !strings.HasPrefix(poster, "https://") {
			publicURL := utils.ConvertPosterToPublicURL(existingSnack.SnackImageURL, "http://localhost:8082/")
			existingSnack.SnackImageURL = &publicURL
		}
	}

	// Save updated snack
	if err := models.UpdateSnack(existingSnack); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update snack"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Snack updated", "snack": existingSnack})
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
