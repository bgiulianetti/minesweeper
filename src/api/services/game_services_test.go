package services

import (
	"testing"

	"github.com/mercadolibre/minesweeper/src/api/constants"
	"github.com/mercadolibre/minesweeper/src/api/dao"
	"github.com/mercadolibre/minesweeper/src/api/domain"
	"github.com/stretchr/testify/assert"
)

func TestGameCreation(t *testing.T) {
	cases := []struct {
		name                  string
		gameConditionsRequest *domain.NewGameConditionsRequest
		expectedGame          *domain.Game
	}{
		{
			name: "OK",
			gameConditionsRequest: &domain.NewGameConditionsRequest{
				UserID:  "test_user",
				Rows:    5,
				Columns: 5,
				Mines:   3,
			},
			expectedGame: &domain.Game{
				Mines: 5,
			},
		},
		{
			name: "FAIL/TOO_MANY_MINES",
			gameConditionsRequest: &domain.NewGameConditionsRequest{
				UserID:  "test_user",
				Rows:    2,
				Columns: 2,
				Mines:   5,
			},
			expectedGame: nil,
		},
		{
			name: "FAIL/NOT_ENOUGH_MINES",
			gameConditionsRequest: &domain.NewGameConditionsRequest{
				UserID:  "test_user",
				Rows:    10,
				Columns: 10,
				Mines:   0,
			},
			expectedGame: nil,
		},
	}

	gameService := &GameService{
		Container: *dao.CreateInMemoryContainer(),
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			newGame, _ := gameService.CreateGame(c.gameConditionsRequest)

			if c.name == "OK" {
				assert.Equal(t, newGame.Mines, c.gameConditionsRequest.Mines)
			} else {
				assert.Equal(t, newGame, c.expectedGame)
			}

		})
	}
}
func TestMinesCount(t *testing.T) {
	cases := []struct {
		name                  string
		gameConditionsRequest *domain.NewGameConditionsRequest
	}{
		{
			name: "OK",
			gameConditionsRequest: &domain.NewGameConditionsRequest{
				Rows:    5,
				Columns: 5,
				Mines:   3,
			},
		},
		{
			name: "FAIL/WITHOUT_MINES",
			gameConditionsRequest: &domain.NewGameConditionsRequest{
				Rows:    10,
				Columns: 10,
				Mines:   5,
			},
		},
	}

	gameService := &GameService{
		Container: *dao.CreateInMemoryContainer(),
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			minesCount := 0
			newGame, _ := gameService.CreateGame(c.gameConditionsRequest)
			for i := 0; i < c.gameConditionsRequest.Columns; i++ {
				for j := 0; j < c.gameConditionsRequest.Rows; j++ {
					if newGame.Board[i][j].HasMine {
						minesCount++
					}
				}
			}
			assert.Equal(t, minesCount, c.gameConditionsRequest.Mines)
		})
	}
}
func TestRevealCell(t *testing.T) {
	cases := []struct {
		name                  string
		gameConditionsRequest *domain.NewGameConditionsRequest
		revealCellRequest     *domain.RevealCellRequest
		expectedGame          *domain.Game
	}{
		{
			name: "OK",
			gameConditionsRequest: &domain.NewGameConditionsRequest{
				UserID:  "test_user",
				Rows:    5,
				Columns: 5,
				Mines:   5,
			},
			revealCellRequest: &domain.RevealCellRequest{
				Row:    0,
				Column: 0,
			},
			expectedGame: &domain.Game{
				Mines: 5,
			},
		},
		{
			name: "FAIL/OUT_OF_BOUNDRIES",
			gameConditionsRequest: &domain.NewGameConditionsRequest{
				UserID:  "test_user",
				Rows:    5,
				Columns: 5,
				Mines:   5,
			},
			revealCellRequest: &domain.RevealCellRequest{
				Row:    5,
				Column: 2,
			},
			expectedGame: nil,
		},
	}

	gameService := &GameService{
		Container: *dao.CreateInMemoryContainer(),
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			newGame, _ := gameService.CreateGame(c.gameConditionsRequest)
			c.revealCellRequest.UserID = c.gameConditionsRequest.UserID
			c.revealCellRequest.GameID = newGame.GameID
			gameWithCellRevealed, _ := gameService.RevealCell(c.revealCellRequest)

			if c.name == "OK" {
				assert.Equal(t, gameWithCellRevealed.Board[c.revealCellRequest.Column][c.revealCellRequest.Row].IsRevealed, true)
			} else {
				assert.Equal(t, gameWithCellRevealed, c.expectedGame)
			}
		})
	}
}
func TestFlagCell(t *testing.T) {
	cases := []struct {
		name                  string
		gameConditionsRequest *domain.NewGameConditionsRequest
		flagCellRequest       *domain.FlagCellRequest
		expectedGame          *domain.Game
	}{
		{
			name: "OK",
			gameConditionsRequest: &domain.NewGameConditionsRequest{
				UserID:  "test_user",
				Rows:    5,
				Columns: 5,
				Mines:   3,
			},
			flagCellRequest: &domain.FlagCellRequest{
				Row:    0,
				Column: 0,
				Flag:   constants.FlagRedFlag,
			},
			expectedGame: &domain.Game{
				Mines: 3,
			},
		},
		{
			name: "FAIL/OUT_OF_BOUNDRIES_ROW",
			gameConditionsRequest: &domain.NewGameConditionsRequest{
				UserID:  "test_user",
				Rows:    5,
				Columns: 5,
				Mines:   5,
			},
			flagCellRequest: &domain.FlagCellRequest{
				Row:    5,
				Column: 2,
				Flag:   constants.FlagRedFlag,
			},
			expectedGame: nil,
		},
		{
			name: "FAIL/OUT_OF_BOUNDRIES_COL",
			gameConditionsRequest: &domain.NewGameConditionsRequest{
				UserID:  "test_user",
				Rows:    5,
				Columns: 5,
				Mines:   5,
			},
			flagCellRequest: &domain.FlagCellRequest{
				Row:    2,
				Column: 5,
				Flag:   constants.FlagRedFlag,
			},
			expectedGame: nil,
		},
		{
			name: "FAIL/FLAG_TWICE",
			gameConditionsRequest: &domain.NewGameConditionsRequest{
				UserID:  "test_user",
				Rows:    5,
				Columns: 5,
				Mines:   3,
			},
			flagCellRequest: &domain.FlagCellRequest{
				Row:    4,
				Column: 2,
				Flag:   constants.FlagRedFlag,
			},
			expectedGame: nil,
		},
		{
			name: "FAIL/FLAG_WITH_QUESTION_MARK_RED_FLAGGED_CELL",
			gameConditionsRequest: &domain.NewGameConditionsRequest{
				UserID:  "test_user",
				Rows:    5,
				Columns: 5,
				Mines:   3,
			},
			flagCellRequest: &domain.FlagCellRequest{
				Row:    4,
				Column: 2,
				Flag:   constants.FlagQuestionMark,
			},
			expectedGame: &domain.Game{
				Mines: 3,
			},
		},
	}

	gameService := &GameService{
		Container: *dao.CreateInMemoryContainer(),
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			newGame, _ := gameService.CreateGame(c.gameConditionsRequest)
			c.flagCellRequest.UserID = c.gameConditionsRequest.UserID
			c.flagCellRequest.GameID = newGame.GameID
			if c.name == "FAIL/FLAG_TWICE" {
				gameService.FlagCell(c.flagCellRequest)
			}
			if c.name == "FAIL/FLAG_WITH_QUESTION_MARK_RED_FLAGGED_CELL" {
				gameService.FlagCell(&domain.FlagCellRequest{
					Row:    c.flagCellRequest.Row,
					Column: c.flagCellRequest.Column,
					Flag:   constants.FlagRedFlag,
					UserID: c.flagCellRequest.UserID,
					GameID: c.flagCellRequest.GameID,
				})
			}
			gameWithFlaggedCell, _ := gameService.FlagCell(c.flagCellRequest)

			if c.name == "OK" || c.name == "FAIL/FLAG_WITH_QUESTION_MARK_RED_FLAGGED_CELL" {
				assert.Equal(t, gameWithFlaggedCell.Board[c.flagCellRequest.Column][c.flagCellRequest.Row].Flag, c.flagCellRequest.Flag)
			} else if c.name == "FAIL/FLAG_TWICE" {
				assert.Equal(t, gameWithFlaggedCell.Board[c.flagCellRequest.Column][c.flagCellRequest.Row].Flag, "")
			} else {
				assert.Equal(t, gameWithFlaggedCell, c.expectedGame)
			}
		})
	}
}
func TestGetGamesByUserID(t *testing.T) {
	cases := []struct {
		name                  string
		gameConditionsRequest *domain.NewGameConditionsRequest
		expectedGame          *domain.Game
	}{
		{
			name: "OK",
			gameConditionsRequest: &domain.NewGameConditionsRequest{
				UserID:  "test_user",
				Rows:    5,
				Columns: 5,
				Mines:   3,
			},
			expectedGame: &domain.Game{
				Mines: 5,
			},
		},
	}

	gameService := &GameService{
		Container: *dao.CreateInMemoryContainer(),
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gameService.CreateGame(c.gameConditionsRequest)
			gameService.CreateGame(c.gameConditionsRequest)

			userGames, _ := gameService.GetGamesByUserID(c.gameConditionsRequest.UserID)

			assert.Equal(t, len(userGames.Games), 2)
			assert.NotEqual(t, userGames.Games[0].GameID, userGames.Games[1].GameID)

		})
	}
}
func TestGameWinOrLose(t *testing.T) {
	cases := []struct {
		name           string
		game           *domain.Game
		expectedStatus string
		x              int
		y              int
	}{
		{
			name: "OK/LOSE",
			x:    0,
			y:    0,
			game: &domain.Game{
				GameID:  1,
				Rows:    2,
				Columns: 2,
				Mines:   1,
				Status:  constants.GameStatusOnGoing,
				Board: [][]domain.Cell{
					{
						domain.Cell{
							HasMine:       true,
							IsRevealed:    false,
							Flag:          "",
							SourroundedBy: 0,
						},
						domain.Cell{
							HasMine:       false,
							IsRevealed:    false,
							Flag:          "",
							SourroundedBy: 1,
						},
					},
					{
						domain.Cell{
							HasMine:       false,
							IsRevealed:    false,
							Flag:          "",
							SourroundedBy: 1,
						},
						domain.Cell{
							HasMine:       false,
							IsRevealed:    false,
							Flag:          "",
							SourroundedBy: 1,
						},
					},
				},
			},
			expectedStatus: constants.GameStatusLose,
		},
		{
			name: "OK/WIN",
			x:    0,
			y:    0,
			game: &domain.Game{
				GameID:  1,
				Rows:    2,
				Columns: 2,
				Mines:   1,
				Status:  constants.GameStatusOnGoing,
				Board: [][]domain.Cell{
					{
						domain.Cell{
							HasMine:       false,
							IsRevealed:    false,
							Flag:          "",
							SourroundedBy: 0,
						},
						domain.Cell{
							HasMine:       true,
							IsRevealed:    true,
							Flag:          "",
							SourroundedBy: 1,
						},
					},
					{
						domain.Cell{
							HasMine:       true,
							IsRevealed:    true,
							Flag:          "",
							SourroundedBy: 1,
						},
						domain.Cell{
							HasMine:       true,
							IsRevealed:    true,
							Flag:          "",
							SourroundedBy: 1,
						},
					},
				},
			},
			expectedStatus: constants.GameResultWon,
		},
	}

	gameService := &GameService{
		Container: *dao.CreateInMemoryContainer(),
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			newGame, _ := gameService.RevealCellFloodFill(c.game, c.x, c.y)
			assert.Equal(t, newGame.Status, c.expectedStatus)
		})
	}
}
