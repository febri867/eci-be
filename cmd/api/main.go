package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "ECI API is running",
		})
	})

	if err := app.Run(":4000"); err != nil {
		log.Fatal(err)
	}
}
