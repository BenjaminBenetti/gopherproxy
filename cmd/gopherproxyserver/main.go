package main

import (
	"github.com/CanadianCommander/gopherproxy/cmd/gopherproxyserver/api"
	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logging.CreateLogger(zap.InfoLevel)
	var gin = gin.Default()
	var apiGroup = gin.Group("/api")

	api.CreateApi(apiGroup)

	gin.Run("0.0.0.0:8080")
}
