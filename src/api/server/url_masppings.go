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

	router.GET("minesweeper/:user_id/games/:game_id",
		middlewares.AdaptHandler(gameController.ValidateGetGameByGameID),
		middlewares.AdaptHandler(gameController.GetGameByGameID),
	)

	router.GET("minesweeper/:user_id/games/:game_id/solution",
		middlewares.AdaptHandler(gameController.ValidateGetGameByGameID),
		middlewares.AdaptHandler(gameController.ShowSolution),
	)

	router.GET("minesweeper/:user_id/games/:game_id/status",
		middlewares.AdaptHandler(gameController.ValidateGetGameByGameID),
		middlewares.AdaptHandler(gameController.ShowStatus),
	)

	router.POST("minesweeper/:user_id/game",
		middlewares.AdaptHandler(gameController.ValidatePost),
		middlewares.AdaptHandler(gameController.CreateNewGame),
	)

	router.POST("minesweeper/:user_id/games/:game_id/flag",
		middlewares.AdaptHandler(gameController.ValidateFlag),
		middlewares.AdaptHandler(gameController.FlagCell),
	)

	router.POST("minesweeper/:user_id/games/:game_id/reveal",
		middlewares.AdaptHandler(gameController.ValidateReveal),
		middlewares.AdaptHandler(gameController.RevealCell),
	)
}
