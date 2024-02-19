package main

import (
	"github.com/CanadianCommander/gopherproxy/internal/api"
	"github.com/gin-gonic/gin"
)

func main() {
	var gin = gin.Default()
	var apiGroup = gin.Group("/api")

	api.CreateApi(apiGroup)

	gin.Run("0.0.0.0:8080")
}
