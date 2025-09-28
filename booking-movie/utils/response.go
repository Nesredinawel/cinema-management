package utils

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"data": data})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}
