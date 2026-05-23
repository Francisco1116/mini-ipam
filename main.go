package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		v1.POST("/networks", func(c*gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Create network endpoint hit!",
			})
		})

		v1.GET("/networks/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{
				"message": "Get network info for ID: " + id,
			})
		})

		v1.POST("/networks/:id/subnets", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{
				"message": "Allocate subnet for network ID: " + id,
			})
		})
	}

	r.Run(":8080")
}