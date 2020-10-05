package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/minesweeper/src/api/domain"
	"github.com/mercadolibre/minesweeper/src/api/errors"
)

// GameService ...
type GameService interface {
	GetGamesByUserID(userID string) (*domain.UserGame, error)
	CreateGame(gameRequest *domain.NewGameConditionsRequest) (*domain.Game, error)
}

// GameController expone los servicios del controller
type GameController struct {
	GameService GameService
}

// Pong allows validation that the API is responding
func (gc GameController) Pong(c *gin.Context) {
	c.Set("skip", true)
	c.JSON(http.StatusOK, "Pong")
}

// GetGamesByUserID ...
func (gc GameController) GetGamesByUserID(c *gin.Context) error {

	userID, ok := c.Get("userID")
	if !ok {
		return &errors.ApiError{"undefined userID", "internal_server_error", http.StatusInternalServerError, ""}
	}
	games, err := gc.GameService.GetGamesByUserID(fmt.Sprintf("%v", userID))
	if err != nil {
		return &errors.ApiError{
			Message:  err.Error(),
			ErrorStr: "internal_server_errror",
			Status:   http.StatusInternalServerError,
			Cause:    "",
		}
	}

	if games == nil {
		c.JSON(http.StatusNotFound, &errors.ApiError{
			Message:  "user not found",
			ErrorStr: "not_found",
			Status:   http.StatusNotFound,
			Cause:    "",
		})
	} else {
		c.JSON(http.StatusOK, games)
	}
	return nil
}

// CreateNewGame creates a new game
func (gc GameController) CreateNewGame(c *gin.Context) error {

	boundBody, ok := c.Get("boundBody")
	if !ok {
		return &errors.ApiError{"undefined boundBody", "internal_server_errror", http.StatusInternalServerError, ""}
	}

	body := boundBody.(*domain.NewGameConditionsRequest)
	game, err := gc.GameService.CreateGame(body)
	if err != nil {
		return &errors.ApiError{
			Message:  err.Error(),
			ErrorStr: "internal_server_errror",
			Status:   http.StatusInternalServerError,
			Cause:    "",
		}
	}
	c.JSON(http.StatusOK, game)
	return nil
}

// ValidateGetGamesByUserID valida el request para obtener una prediccion de usuario
func (gc GameController) ValidateGetGamesByUserID(c *gin.Context) error {

	// Validate user_id
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, &errors.ApiError{"invalid user_id", "bad_request", http.StatusBadRequest, ""})
		return nil
	}
	c.Set("userID", userID)
	return nil
}

// ValidateGetGameByGameID valida el request para obtener una prediccion de usuario
func (gc GameController) ValidateGetGameByGameID(c *gin.Context) error {

	// Validate user_id
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, &errors.ApiError{"invalid user_id", "bad_request", http.StatusBadRequest, ""})
		return nil
	}

	gameID := c.Param("game_id")
	intGameID, err := strconv.ParseInt(gameID, 10, 64)
	if err != nil || intGameID < 0 {
		c.JSON(http.StatusBadRequest, &errors.ApiError{"invalid game_id: " + err.Error(), "bad_request", http.StatusBadRequest, ""})
		return nil
	}

	c.Set("userID", userID)
	c.Set("gameID", intGameID)
	return nil
}

// ValidatePost validates the post request
func (gc GameController) ValidatePost(c *gin.Context) error {

	boundBody := &domain.NewGameConditionsRequest{}
	err := c.BindJSON(boundBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  err.Error(),
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
			Cause:    "",
		})
		return nil
	}

	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  "user_id is mandatory",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
			Cause:    "",
		})
		return nil
	}

	if boundBody.Columns <= 0 {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  "columns must be grater than 0",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
			Cause:    "",
		})
		return nil
	}

	if boundBody.Rows <= 0 {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  "rows must be grater than 0",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
			Cause:    "",
		})
		return nil
	}

	if boundBody.Mines <= 0 || boundBody.Mines >= boundBody.Columns*boundBody.Rows {
		minesError := &errors.ApiError{
			Message:  "the number of mines must be at least one, and lower than total of cells in the game",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
			Cause:    "",
		}
		c.JSON(http.StatusBadRequest, minesError)
		return nil
	}

	boundBody.UserID = userID
	c.Set("boundBody", boundBody)

	return nil
}
