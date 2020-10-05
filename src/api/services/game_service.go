package services

import (
	"math/rand"
	"time"

	"github.com/mercadolibre/minesweeper/src/api/dao"
	"github.com/mercadolibre/minesweeper/src/api/domain"
)

// GameService ...
type GameService struct {
	Container dao.InMemoryContainer
}

// CreateGame creates a new game
func (gs *GameService) CreateGame(gameRequest *domain.NewGameConditionsRequest) (*domain.Game, error) {

	newGame, err := createNewGameFromRequest(gameRequest)
	if err != nil {
		return nil, err
	}

	userGames, getErr := gs.Container.Get(gameRequest.UserID)
	if getErr != nil {
		return nil, getErr
	}

	if userGames == nil {
		userGames = &domain.UserGame{
			UserID: gameRequest.UserID,
			Games:  make([]*domain.Game, 0),
		}
	}
	userGames.Games = append(userGames.Games, newGame)
	upsertErr := gs.Container.Upsert(userGames)
	if upsertErr != nil {
		return nil, upsertErr
	}
	return newGame, nil
}

// GetGamesByUserID returns all the games by a user
func (gs *GameService) GetGamesByUserID(userID string) (*domain.UserGame, error) {
	userGame, err := gs.Container.Get(userID)
	if err != nil {
		return nil, err
	}
	if userGame == nil {
		return nil, nil
	}
	return userGame, nil
}

func createNewGameFromRequest(gameRequest *domain.NewGameConditionsRequest) (*domain.Game, error) {

	gameID := generateUniqueID()
	newGame := &domain.Game{
		Start:   time.Now(),
		Columns: gameRequest.Columns,
		Rows:    gameRequest.Rows,
		Status:  "on_going",
		Board:   initializeBoard(gameRequest.Columns, gameRequest.Rows, gameRequest.Mines, gameID),
		GameID:  gameID,
	}
	return newGame, nil
}

func initializeBoard(columns int, rows int, mines int, seed int64) [][]domain.Cell {

	boardWithMines := placeMines(columns, rows, mines, seed)
	boardComplete := setNeighgoursCount(boardWithMines, columns, rows)
	return boardComplete
}

func placeMines(columns int, rows int, mines int, seed int64) [][]domain.Cell {

	newBoard := make([][]domain.Cell, columns)
	for col := range newBoard {
		newBoard[col] = make([]domain.Cell, rows)
	}

	rand.Seed(seed)
	minesLeftToPlace := mines
	for minesLeftToPlace > 0 {
		columnNumber := rand.Intn(columns)
		rowNumber := rand.Intn(rows)
		if !newBoard[columnNumber][rowNumber].HasMine {
			newBoard[columnNumber][rowNumber].HasMine = true
			minesLeftToPlace--
		}
	}
	return newBoard
}

func setNeighgoursCount(board [][]domain.Cell, columns int, rows int) [][]domain.Cell {

	for i := 0; i < columns; i++ {
		for j := 0; j < rows; j++ {
			board[i][j].SourroundedBy = countNeighbours(i, j, board, columns, rows)
		}
	}
	return board
}

func countNeighbours(x int, y int, board [][]domain.Cell, columns int, rows int) int {

	totalNeighbours := 0
	for xOffset := -1; xOffset <= 1; xOffset++ {
		for yOffset := -1; yOffset <= 1; yOffset++ {
			i := x + xOffset
			j := y + yOffset
			if i > -1 && i < columns && j > -1 && j < rows {
				neighbour := board[i][j]
				if neighbour.HasMine {
					totalNeighbours++
				}
			}
		}
	}
	return totalNeighbours
}

func generateUniqueID() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}
