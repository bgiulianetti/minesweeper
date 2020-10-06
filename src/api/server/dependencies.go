package server

import (
	"github.com/mercadolibre/minesweeper/src/api/controllers"
	"github.com/mercadolibre/minesweeper/src/api/dao"
	"github.com/mercadolibre/minesweeper/src/api/services"
)

func resolveGameController() *controllers.GameController {

	return &controllers.GameController{
		GameService: &services.GameService{
			Container: *dao.CreateInMemoryContainer(),
		},
	}
}
