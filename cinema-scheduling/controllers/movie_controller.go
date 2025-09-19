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

	// Get form values
	movie.Title = c.PostForm("title")
	movie.Description = c.PostForm("description")
	movie.Genre = c.PostForm("genre")

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

	// Create movie
	if err := models.CreateMovie(&movie); err != nil {
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

	// Fetch existing movie from DB
	existingMovie, err := models.GetMovieByID(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movie"})
		return
	}
	if existingMovie == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	// Bind incoming JSON into a map
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only fields provided
	if title, ok := req["title"].(string); ok {
		existingMovie.Title = title
	}
	if description, ok := req["description"].(string); ok {
		existingMovie.Description = description
	}
	if genre, ok := req["genre"].(string); ok {
		existingMovie.Genre = genre
	}
	if dur, ok := req["duration"].(float64); ok { // JSON numbers are float64
		existingMovie.Duration = int(dur)
	}
	if year, ok := req["release_year"].(float64); ok {
		existingMovie.ReleaseYear = int(year)
	}
	if rating, ok := req["rating"].(float64); ok {
		existingMovie.Rating = &rating
	}
	if poster, ok := req["image_poster_url"].(string); ok {
		existingMovie.ImagePosterURL = &poster

		// convert to public URL if needed
		if !strings.HasPrefix(poster, "http://") && !strings.HasPrefix(poster, "https://") {
			publicURL := utils.ConvertPosterToPublicURL(existingMovie.ImagePosterURL, "http://localhost:8082/")
			existingMovie.ImagePosterURL = &publicURL
		}
	}

	// Save updated movie
	if err := models.UpdateMovie(existingMovie); err != nil {
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
