package api

import (
	"github.com/gin-gonic/gin"
)

func CreateApi(routeBuilder *gin.RouterGroup) *gin.RouterGroup {
	routeBuilder.GET(ConnectionListenerRoute, ConnectionListen)

	return routeBuilder
}
