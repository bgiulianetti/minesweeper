package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/minesweeper/src/api/controllers"
	"github.com/mercadolibre/minesweeper/src/api/middlewares"
)

// Associate URLs with controllers
func mapUrlsToControllers(router *gin.Engine, gameController *controllers.GameController) {

	router.GET("/ping", gameController.Pong)

	router.GET("minesweeper/:user_id/games",
		middlewares.AdaptHandler(gameController.ValidateGetGamesByUserID),
		middlewares.AdaptHandler(gameController.GetGamesByUserID),
	)
}
