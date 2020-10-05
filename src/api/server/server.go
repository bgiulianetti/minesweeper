package server

import "github.com/gin-gonic/gin"

// New creates a new router
func New() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	bootstrap(router)

	return router
}
