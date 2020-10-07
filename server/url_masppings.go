package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/minesweeper/controllers"
	"github.com/mercadolibre/minesweeper/middlewares"
)

// Associate URLs with controllers
func mapUrlsToControllers(router *gin.Engine, gameController *controllers.GameController) {
	router.GET("/ping", gameController.Pong)

	router.GET("minesweeper/users/:user_id/games",
		middlewares.AdaptHandler(gameController.ValidateGetGamesByUserID),
		middlewares.AdaptHandler(gameController.GetGamesByUserID),
	)

	router.GET("minesweeper/users/:user_id/games/:game_id",
		middlewares.AdaptHandler(gameController.ValidateGetGameByGameID),
		middlewares.AdaptHandler(gameController.GetGameByGameID),
	)

	router.GET("minesweeper/users/:user_id/games/:game_id/solution",
		middlewares.AdaptHandler(gameController.ValidateGetGameByGameID),
		middlewares.AdaptHandler(gameController.ShowSolution),
	)

	router.GET("minesweeper/users/:user_id/games/:game_id/status",
		middlewares.AdaptHandler(gameController.ValidateGetGameByGameID),
		middlewares.AdaptHandler(gameController.ShowStatus),
	)

	router.POST("minesweeper/users/:user_id/games",
		middlewares.AdaptHandler(gameController.ValidatePost),
		middlewares.AdaptHandler(gameController.CreateNewGame),
	)

	router.POST("minesweeper/users/:user_id/games/:game_id/flag",
		middlewares.AdaptHandler(gameController.ValidateFlag),
		middlewares.AdaptHandler(gameController.FlagCell),
	)

	router.POST("minesweeper/users/:user_id/games/:game_id/reveal",
		middlewares.AdaptHandler(gameController.ValidateReveal),
		middlewares.AdaptHandler(gameController.RevealCell),
	)

	router.DELETE("minesweeper/games",
		middlewares.AdaptHandler(gameController.DeleteAllGames),
	)

	router.GET("minesweeper/games",
		middlewares.AdaptHandler(gameController.GetllGames),
	)
}
