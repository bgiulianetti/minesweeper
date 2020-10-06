package services

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/mercadolibre/minesweeper/src/api/constants"
	"github.com/mercadolibre/minesweeper/src/api/dao"
	"github.com/mercadolibre/minesweeper/src/api/domain"
	"github.com/mercadolibre/minesweeper/src/api/errors"
)

// GameService ...
type GameService struct {
	Container dao.InMemoryContainer
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

// GetGameByGameID returns all the games by a user
func (gs *GameService) GetGameByGameID(userID string, gameID int64) (*domain.Game, error) {
	userGames, err := gs.Container.Get(userID)
	if err != nil {
		return nil, err
	}

	if userGames == nil {
		return nil, nil
	}

	for _, game := range userGames.Games {
		if game.GameID == gameID {
			return game, nil
		}
	}
	return nil, nil
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

// ShowSolution shows the solution's solution
func (gs *GameService) ShowSolution(userID string, gameID int64) (string, error) {

	game, err := gs.GetGameByGameID(userID, gameID)
	if err != nil {
		return "", err
	}
	if game == nil {
		return "", nil
	}
	return boardSolutionToString(game.Board), nil
}

// FlagCell flags a cell on the board
func (gs *GameService) FlagCell(flagRequest *domain.FlagCellRequest) (*domain.Game, error) {
	userGame, err := gs.GetGamesByUserID(flagRequest.UserID)
	if err != nil {
		return nil, err
	}

	gameIndex := getGameIndex(flagRequest.GameID, userGame)
	if gameIndex == -1 {
		return nil, nil
	}

	if flagRequest.Column > userGame.Games[gameIndex].Columns {
		return nil, &errors.ApiError{
			Message:  "flag out of boundries (columns exceeded)",
			ErrorStr: "out_of_boundries",
		}
	}

	if flagRequest.Row > userGame.Games[gameIndex].Rows {
		return nil, &errors.ApiError{
			Message:  "flag out of boundries (rows exceeded)",
			ErrorStr: "out_of_boundries",
		}
	}

	if !userGame.Games[gameIndex].Board[flagRequest.Column][flagRequest.Row].IsRevealed {
		userGame.Games[gameIndex].Board[flagRequest.Column][flagRequest.Row].Flag = flagRequest.Flag
		upsertErr := gs.Container.Upsert(userGame)
		if err != nil {
			return nil, upsertErr
		}
	}
	return userGame.Games[gameIndex], nil
}

// ShowStatus shows the solution's solution
func (gs *GameService) ShowStatus(userID string, gameID int64) (string, error) {

	game, err := gs.GetGameByGameID(userID, gameID)
	if err != nil {
		return "", err
	}
	if game == nil {
		return "", nil
	}
	return boardToString(game.Board), nil
}

// RevealCell reveals a cell
func (gs *GameService) RevealCell(revealCellRequest *domain.RevealCellRequest) (*domain.Game, error) {
	userGame, err := gs.GetGamesByUserID(revealCellRequest.UserID)
	if err != nil {
		return nil, err
	}

	gameIndex := getGameIndex(revealCellRequest.GameID, userGame)
	if gameIndex == -1 {
		return nil, nil
	}

	if revealCellRequest.Column > userGame.Games[gameIndex].Columns {
		return nil, &errors.ApiError{
			Message:  "flag out of boundries (columns exceeded)",
			ErrorStr: "out_of_boundries",
		}
	}

	if revealCellRequest.Row > userGame.Games[gameIndex].Rows {
		return nil, &errors.ApiError{
			Message:  "flag out of boundries (rows exceeded)",
			ErrorStr: "out_of_boundries",
		}
	}

	if userGame.Games[gameIndex].Board[revealCellRequest.Column][revealCellRequest.Row].HasMine {
		userGame.Games[gameIndex].Status = constants.GameStatusLose
		userGame.Games[gameIndex].Finish = time.Now()
	} else {
		revealCell(userGame.Games[gameIndex].Board, revealCellRequest.Column, revealCellRequest.Row, userGame.Games[gameIndex].Columns, userGame.Games[gameIndex].Rows)
		if checkIfWon(userGame.Games[gameIndex].Board, userGame.Games[gameIndex].Columns, userGame.Games[gameIndex].Rows) {
			userGame.Games[gameIndex].Status = constants.GameResultWon
			userGame.Games[gameIndex].Finish = time.Now()
		}
	}
	upsertErr := gs.Container.Upsert(userGame)
	if err != nil {
		return nil, upsertErr
	}
	return userGame.Games[gameIndex], nil
}

func checkIfWon(board [][]domain.Cell, columns, rows int) bool {

	for i := 0; i < columns; i++ {
		for j := 0; j < rows; j++ {
			if board[i][j].HasMine && board[i][j].Flag != constants.FlagRedFlag {
				return false
			}
			if !board[i][j].HasMine && !board[i][j].IsRevealed {
				return false
			}
		}
	}
	return true
}

func getAllRedFlaggedAndRevealedCellsFromBoard(board [][]domain.Cell, columns, rows int) (int, int) {

	redFlaggedCellsCount := 0
	revealedCellsCount := 0
	for i := 0; i < columns; i++ {
		for j := 0; j < rows; j++ {
			if board[i][j].Flag == constants.FlagRedFlag {
				redFlaggedCellsCount++
			} else if board[i][j].IsRevealed {
				revealedCellsCount++
			}
		}
	}
	return redFlaggedCellsCount, revealedCellsCount
}

func revealCell(board [][]domain.Cell, column int, row int, totalColumns int, totalRows int) {
	board[column][row].IsRevealed = true
	if board[column][row].SourroundedBy == 0 {
		for xOffset := -1; xOffset <= 1; xOffset++ {
			for yOffset := -1; yOffset <= 1; yOffset++ {
				i := column + xOffset
				j := row + yOffset
				if i > -1 && i < totalColumns && j > -1 && j < totalRows {
					neighbour := board[i][j]
					if !neighbour.HasMine && !neighbour.IsRevealed && neighbour.Flag != constants.FlagRedFlag && neighbour.Flag != constants.FlagQuestionMark {
						revealCell(board, i, j, totalColumns, totalRows)
					}
				}
			}
		}
	}
}

func createNewGameFromRequest(gameRequest *domain.NewGameConditionsRequest) (*domain.Game, error) {

	gameID := generateUniqueID()
	newGame := &domain.Game{
		Start:   time.Now(),
		Columns: gameRequest.Columns,
		Rows:    gameRequest.Rows,
		Status:  constants.GameStatusOnGoing,
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

func boardSolutionToString(board [][]domain.Cell) string {

	stringBoard := ""
	for i := range board {
		for _, cell := range board[i] {
			if cell.HasMine {
				stringBoard += " * "
			} else if cell.SourroundedBy == 0 {
				stringBoard += " _ "
			} else {
				stringBoard += " " + strconv.Itoa(cell.SourroundedBy) + " "
			}
		}
		stringBoard += "|"
	}
	return stringBoard
}

func boardToString(board [][]domain.Cell) string {

	stringBoard := ""
	for i := range board {
		for _, cell := range board[i] {
			if cell.IsRevealed {
				if cell.HasMine {
					stringBoard += " * "
				} else if cell.SourroundedBy == 0 {
					stringBoard += " _ "
				} else {
					stringBoard += " " + strconv.Itoa(cell.SourroundedBy) + " "
				}
			} else if cell.Flag == constants.FlagRedFlag {
				stringBoard += " F "
			} else if cell.Flag == constants.FlagQuestionMark {
				stringBoard += " ? "
			} else {
				stringBoard += " □ "
			}
		}
		stringBoard += "|"
	}
	return stringBoard
}

func getGameIndex(gameID int64, userGame *domain.UserGame) int {
	for i, game := range userGame.Games {
		if game.GameID == gameID {
			return i
		}
		i++
	}
	return -1
}
