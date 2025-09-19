package controllers

import (
	"cinema-scheduling/models"
	"cinema-scheduling/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ---------------- Add Movie ----------------
func AddMovie(c *gin.Context) {
	var movie models.Movie
	var genreIDs []int

	// Get form values
	movie.Title = c.PostForm("title")
	movie.Description = c.PostForm("description")

	// Convert numeric fields safely
	if durStr := c.PostForm("duration"); durStr != "" {
		if dur, err := strconv.Atoi(durStr); err == nil {
			movie.Duration = dur
		}
	}
	if yearStr := c.PostForm("release_year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			movie.ReleaseYear = year
		}
	}
	if ratingStr := c.PostForm("rating"); ratingStr != "" {
		if r, err := strconv.ParseFloat(ratingStr, 64); err == nil {
			movie.Rating = &r
		}
	}

	// Parse genre_ids (comma-separated in form-data)
	if genreStr := c.PostForm("genre_ids"); genreStr != "" {
		for _, g := range strings.Split(genreStr, ",") {
			if gid, err := strconv.Atoi(strings.TrimSpace(g)); err == nil {
				genreIDs = append(genreIDs, gid)
			}
		}
	}

	// Handle poster
	file, _ := c.FormFile("image_poster_url")
	posterURL := c.PostForm("image_poster_url")

	savedPoster, err := utils.SaveMoviePoster(file, posterURL, "uploads/posters", c.SaveUploadedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save poster"})
		return
	}
	if savedPoster != nil {
		movie.ImagePosterURL = savedPoster
	}

	// Convert to public URL only if relative path
	if movie.ImagePosterURL != nil {
		pathStr := *movie.ImagePosterURL
		if !strings.HasPrefix(pathStr, "http://") && !strings.HasPrefix(pathStr, "https://") {
			publicURL := utils.ConvertPosterToPublicURL(movie.ImagePosterURL, "http://localhost:8082/")
			movie.ImagePosterURL = &publicURL
		}
	}

	// Create movie with genres
	if err := models.CreateMovie(&movie, genreIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Movie created", "movie": movie})
}

// ---------------- List Movies ----------------
func ListMovies(c *gin.Context) {
	movies, err := models.GetAllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}

	for _, m := range movies {
		if m.ImagePosterURL != nil {
			pathStr := *m.ImagePosterURL
			if !strings.HasPrefix(pathStr, "http://") && !strings.HasPrefix(pathStr, "https://") {
				publicURL := utils.ConvertPosterToPublicURL(m.ImagePosterURL, "http://localhost:8082/")
				m.ImagePosterURL = &publicURL
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"movies": movies})
}

// ---------------- Get Movie by ID ----------------
func GetMovie(c *gin.Context) {
	movieIDStr := c.Param("movie_id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	movie, err := models.GetMovieByID(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movie"})
		return
	}
	if movie == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	if movie.ImagePosterURL != nil {
		pathStr := *movie.ImagePosterURL
		if !strings.HasPrefix(pathStr, "http://") && !strings.HasPrefix(pathStr, "https://") {
			publicURL := utils.ConvertPosterToPublicURL(movie.ImagePosterURL, "http://localhost:8082/")
			movie.ImagePosterURL = &publicURL
		}
	}

	c.JSON(http.StatusOK, gin.H{"movie": movie})
}

// ---------------- Update Movie ----------------
func UpdateMovie(c *gin.Context) {
	movieIDStr := c.Param("movie_id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	// Fetch existing movie
	existingMovie, err := models.GetMovieByID(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movie"})
		return
	}
	if existingMovie == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	// Parse incoming JSON
	var req struct {
		Title          *string  `json:"title"`
		Description    *string  `json:"description"`
		Duration       *int     `json:"duration"`
		ReleaseYear    *int     `json:"release_year"`
		Rating         *float64 `json:"rating"`
		ImagePosterURL *string  `json:"image_poster_url"`
		GenreIDs       []int    `json:"genre_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Apply updates
	if req.Title != nil {
		existingMovie.Title = *req.Title
	}
	if req.Description != nil {
		existingMovie.Description = *req.Description
	}
	if req.Duration != nil {
		existingMovie.Duration = *req.Duration
	}
	if req.ReleaseYear != nil {
		existingMovie.ReleaseYear = *req.ReleaseYear
	}
	if req.Rating != nil {
		existingMovie.Rating = req.Rating
	}
	if req.ImagePosterURL != nil {
		existingMovie.ImagePosterURL = req.ImagePosterURL
		if !strings.HasPrefix(*req.ImagePosterURL, "http://") && !strings.HasPrefix(*req.ImagePosterURL, "https://") {
			publicURL := utils.ConvertPosterToPublicURL(req.ImagePosterURL, "http://localhost:8082/")
			existingMovie.ImagePosterURL = &publicURL
		}
	}

	// Save update with genres
	if err := models.UpdateMovie(existingMovie, req.GenreIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie updated", "movie": existingMovie})
}

// ---------------- Delete Movie ----------------
func DeleteMovie(c *gin.Context) {
	movieIDStr := c.Param("movie_id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	if err := models.DeleteMovie(movieID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted"})
}
