package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/minesweeper/src/api/controllers"
)

// Associate URLs with controllers
func mapUrlsToControllers(router *gin.Engine, gameController *controllers.GameController) {

	router.GET("/ping", gameController.Pong)
}
