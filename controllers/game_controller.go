package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/minesweeper/constants"
	"github.com/mercadolibre/minesweeper/domain"
	"github.com/mercadolibre/minesweeper/errors"
)

// GameService ...
type GameService interface {
	GetGamesByUserID(userID string) (*domain.UserGame, error)
	GetGameByGameID(userID string, gameID int64) (*domain.Game, error)
	CreateGame(gameRequest *domain.NewGameConditionsRequest) (*domain.Game, error)
	ShowSolution(userID string, gameID int64) (string, error)
	FlagCell(flagRequest *domain.FlagCellRequest) (*domain.Game, error)
	ShowStatus(userID string, gameID int64) (string, error)
	RevealCell(revealCellRequest *domain.RevealCellRequest) (*domain.Game, error)
	DeleteAllGames() error
	GetAllGames() ([]*domain.UserGame, error)
}

// GameController expone los servicios del controller
type GameController struct {
	GameService GameService
}

// Pong allows validation that the API is responding
func (gc GameController) Pong(c *gin.Context) {
	c.Set("skip", true)
	c.JSON(http.StatusOK, "Pong from: minesweeper")
}

// GetGamesByUserID ...
func (gc GameController) GetGamesByUserID(c *gin.Context) error {

	userID, ok := c.Get("userID")
	if !ok {
		return &errors.ApiError{"undefined userID", "internal_server_error", http.StatusInternalServerError}
	}
	games, err := gc.GameService.GetGamesByUserID(fmt.Sprintf("%v", userID))
	if err != nil {
		return &errors.ApiError{
			Message:  err.Error(),
			ErrorStr: "internal_server_errror",
			Status:   http.StatusInternalServerError,
		}
	}

	if games == nil {
		c.JSON(http.StatusNotFound, &errors.ApiError{
			Message:  "user not found",
			ErrorStr: "not_found",
			Status:   http.StatusNotFound,
		})
	} else {
		c.JSON(http.StatusOK, games)
	}
	return nil
}

// GetGameByGameID ...
func (gc GameController) GetGameByGameID(c *gin.Context) error {

	userID, ok := c.Get("userID")
	if !ok {
		return &errors.ApiError{"undefined userID", "internal_server_error", http.StatusInternalServerError}
	}

	gameID, ok := c.Get("gameID")
	if !ok {
		return &errors.ApiError{"undefined userID", "internal_server_error", http.StatusInternalServerError}
	}

	game, err := gc.GameService.GetGameByGameID((fmt.Sprintf("%v", userID)), gameID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errors.ApiError{
			Message:  err.Error(),
			ErrorStr: "internal_server_errror",
			Status:   http.StatusInternalServerError,
		})
	}

	if game == nil {
		c.JSON(http.StatusNotFound, &errors.ApiError{
			Message:  "game not found",
			ErrorStr: "not_found",
			Status:   http.StatusNotFound,
		})
	} else {
		c.JSON(http.StatusOK, game)
	}
	return nil
}

// CreateNewGame creates a new game
func (gc GameController) CreateNewGame(c *gin.Context) error {

	boundBody, ok := c.Get("boundBody")
	if !ok {
		return &errors.ApiError{"undefined boundBody", "internal_server_errror", http.StatusInternalServerError}
	}

	body := boundBody.(*domain.NewGameConditionsRequest)
	game, err := gc.GameService.CreateGame(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errors.ApiError{
			Message:  err.Error(),
			ErrorStr: "internal_server_errror",
			Status:   http.StatusInternalServerError,
		})
		return nil
	}
	c.JSON(http.StatusCreated, game)
	return nil
}

// ShowSolution shows the solution of a game
func (gc GameController) ShowSolution(c *gin.Context) error {

	userID, ok := c.Get("userID")
	if !ok {
		return &errors.ApiError{"undefined userID", "internal_server_error", http.StatusInternalServerError}
	}

	gameID, ok := c.Get("gameID")
	if !ok {
		return &errors.ApiError{"undefined userID", "internal_server_error", http.StatusInternalServerError}
	}

	solution, err := gc.GameService.ShowSolution((fmt.Sprintf("%v", userID)), gameID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errors.ApiError{
			Message:  err.Error(),
			ErrorStr: "internal_server_errror",
			Status:   http.StatusInternalServerError,
		})
	}

	if solution == "" {
		c.JSON(http.StatusInternalServerError, &errors.ApiError{
			Message:  "game_not_found",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
	} else {
		c.JSON(http.StatusOK, solution)
	}
	return nil
}

// ShowStatus shows the current status of a game
func (gc GameController) ShowStatus(c *gin.Context) error {

	userID, ok := c.Get("userID")
	if !ok {
		return &errors.ApiError{"undefined userID", "internal_server_error", http.StatusInternalServerError}
	}

	gameID, ok := c.Get("gameID")
	if !ok {
		return &errors.ApiError{"undefined userID", "internal_server_error", http.StatusInternalServerError}
	}

	solution, err := gc.GameService.ShowStatus((fmt.Sprintf("%v", userID)), gameID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errors.ApiError{
			Message:  err.Error(),
			ErrorStr: "internal_server_errror",
			Status:   http.StatusInternalServerError,
		})
	}

	if solution == "" {
		c.JSON(http.StatusInternalServerError, &errors.ApiError{
			Message:  "game_not_found",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
	} else {
		c.JSON(http.StatusOK, solution)
	}
	return nil
}

// FlagCell flags a cell of the board
func (gc GameController) FlagCell(c *gin.Context) error {

	boundBody, ok := c.Get("boundBody")
	if !ok {
		return &errors.ApiError{"undefined boundBody", "internal_server_errror", http.StatusInternalServerError}
	}

	body := boundBody.(*domain.FlagCellRequest)
	game, err := gc.GameService.FlagCell(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return nil
	}

	if game == nil {
		c.JSON(http.StatusBadRequest, game)
		return nil
	}
	c.JSON(http.StatusOK, game)
	return nil
}

// RevealCell reveals a cell of the board
func (gc GameController) RevealCell(c *gin.Context) error {

	boundBody, ok := c.Get("boundBody")
	if !ok {
		return &errors.ApiError{"undefined boundBody", "internal_server_errror", http.StatusInternalServerError}
	}

	body := boundBody.(*domain.RevealCellRequest)
	game, err := gc.GameService.RevealCell(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  err.Error(),
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
		return nil
	}

	if game == nil {
		c.JSON(http.StatusBadRequest, nil)
		return nil
	}

	c.JSON(http.StatusOK, game)
	return nil
}

// DeleteAllGames deletes all games
func (gc GameController) DeleteAllGames(c *gin.Context) error {

	err := gc.GameService.DeleteAllGames()
	if err != nil {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  err.Error(),
			ErrorStr: "internal_server_error",
			Status:   http.StatusInternalServerError,
		})
		return nil
	}
	c.JSON(http.StatusOK, nil)
	return nil
}

// GetllGames gets all games
func (gc GameController) GetllGames(c *gin.Context) error {

	games, err := gc.GameService.GetAllGames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errors.ApiError{
			Message:  err.Error(),
			ErrorStr: "internal_server_error",
			Status:   http.StatusInternalServerError,
		})
		return nil
	}

	c.JSON(http.StatusOK, games)
	return nil
}

// ValidateGetGamesByUserID valida el request para obtener una prediccion de usuario
func (gc GameController) ValidateGetGamesByUserID(c *gin.Context) error {

	// Validate user_id
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, &errors.ApiError{"invalid user_id", "bad_request", http.StatusBadRequest})
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
		c.JSON(http.StatusBadRequest, &errors.ApiError{"invalid user_id", "bad_request", http.StatusBadRequest})
		return nil
	}

	gameID := c.Param("game_id")
	intGameID, err := strconv.ParseInt(gameID, 10, 64)
	if err != nil || intGameID < 0 {
		c.JSON(http.StatusBadRequest, &errors.ApiError{"invalid game_id: " + err.Error(), "bad_request", http.StatusBadRequest})
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
		})
		return nil
	}

	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  "user_id is mandatory",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
		return nil
	}

	if boundBody.Columns <= 0 || boundBody.Columns > 30 {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  "columns must be greater than 0 and less or equal than 30",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
		return nil
	}

	if boundBody.Rows <= 0 || boundBody.Rows > 30 {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  "rows must be greater than 0 and less or equal than 30",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
		return nil
	}

	if boundBody.Rows != boundBody.Columns {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  "rows and columns must be equals",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
		return nil
	}

	if boundBody.Mines <= 0 || boundBody.Mines > boundBody.Columns*boundBody.Rows {
		minesError := &errors.ApiError{
			Message:  "the number of mines must be at least one, and less or equal than total of cells in the game",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, minesError)
		return nil
	}

	boundBody.UserID = userID
	c.Set("boundBody", boundBody)

	return nil
}

// ValidateFlag validates the flag request
func (gc GameController) ValidateFlag(c *gin.Context) error {

	boundBody := &domain.FlagCellRequest{}
	err := c.BindJSON(boundBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  err.Error(),
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
		return nil
	}

	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  "user_id is mandatory",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
		return nil
	}

	gameID := c.Param("game_id")
	intGameID, err := strconv.ParseInt(gameID, 10, 64)
	if err != nil || intGameID < 0 {
		c.JSON(http.StatusBadRequest, &errors.ApiError{"invalid game_id: " + err.Error(), "bad_request", http.StatusBadRequest})
		return nil
	}

	if boundBody.Column < 0 {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  "column must be grater than 0",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
		return nil
	}

	if boundBody.Row < 0 {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  "row must be grater than 0",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
		return nil
	}

	if boundBody.Flag != constants.FlagQuestionMark && boundBody.Flag != constants.FlagRedFlag {
		minesError := &errors.ApiError{
			Message:  "Available flag options: [" + constants.FlagQuestionMark + ", [" + constants.FlagRedFlag + "]",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, minesError)
		return nil
	}

	boundBody.UserID = userID
	boundBody.GameID = intGameID
	c.Set("boundBody", boundBody)

	return nil
}

// ValidateReveal validates the reveal request
func (gc GameController) ValidateReveal(c *gin.Context) error {

	boundBody := &domain.RevealCellRequest{}
	err := c.BindJSON(boundBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  err.Error(),
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
		return nil
	}

	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  "user_id is mandatory",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
		return nil
	}

	gameID := c.Param("game_id")
	intGameID, err := strconv.ParseInt(gameID, 10, 64)
	if err != nil || intGameID < 0 {
		c.JSON(http.StatusBadRequest, &errors.ApiError{"invalid game_id: " + err.Error(), "bad_request", http.StatusBadRequest})
		return nil
	}

	if boundBody.Column < 0 {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  "columns must be grater than 0",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
		return nil
	}

	if boundBody.Row < 0 {
		c.JSON(http.StatusBadRequest, &errors.ApiError{
			Message:  "rows must be grater than 0",
			ErrorStr: "bad_request",
			Status:   http.StatusBadRequest,
		})
		return nil
	}

	boundBody.UserID = userID
	boundBody.GameID = intGameID
	c.Set("boundBody", boundBody)

	return nil
}
