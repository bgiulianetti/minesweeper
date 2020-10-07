package server

import (
	"github.com/mercadolibre/minesweeper/controllers"
	"github.com/mercadolibre/minesweeper/dao"
	"github.com/mercadolibre/minesweeper/services"
)

func resolveGameController() *controllers.GameController {

	return &controllers.GameController{
		GameService: &services.GameService{
			InMemoryContainer: *dao.CreateInMemoryContainer(),
			Container:         dao.CreateContainer(),
		},
	}
}
