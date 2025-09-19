package controllers

import (
	"cinema-scheduling/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ---------------- Add Genre ----------------
func AddGenre(c *gin.Context) {
	var genre models.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if genre.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Genre name is required"})
		return
	}

	if err := models.CreateGenre(&genre); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create genre"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Genre created", "genre": genre})
}

// ---------------- List Genres ----------------
func ListGenres(c *gin.Context) {
	genres, err := models.GetAllGenres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch genres"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"genres": genres})
}

// ---------------- Get Genre by ID ----------------
func GetGenre(c *gin.Context) {
	genreIDStr := c.Param("genre_id")
	genreID, err := strconv.Atoi(genreIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	genre, err := models.GetGenreByID(genreID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch genre"})
		return
	}
	if genre == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"genre": genre})
}

// ---------------- Update Genre ----------------
func UpdateGenre(c *gin.Context) {
	genreIDStr := c.Param("genre_id")
	genreID, err := strconv.Atoi(genreIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	existingGenre, err := models.GetGenreByID(genreID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch genre"})
		return
	}
	if existingGenre == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if name, ok := req["name"].(string); ok && name != "" {
		existingGenre.Name = name
	}

	if err := models.UpdateGenre(existingGenre); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update genre"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Genre updated", "genre": existingGenre})
}

// ---------------- Delete Genre ----------------
func DeleteGenre(c *gin.Context) {
	genreIDStr := c.Param("genre_id")
	genreID, err := strconv.Atoi(genreIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	if err := models.DeleteGenre(genreID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete genre"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Genre deleted"})
}
