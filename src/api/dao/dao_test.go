package dao

import (
	"testing"

	"github.com/mercadolibre/minesweeper/src/api/constants"
	"github.com/mercadolibre/minesweeper/src/api/domain"
	"github.com/stretchr/testify/assert"
)

func TestGameUpsert(t *testing.T) {
	cases := []struct {
		name           string
		userGame       *domain.UserGame
		userID         string
		expectedStatus string
	}{
		{
			name: "OK/INSERT",
			userGame: &domain.UserGame{
				UserID: "test_user_1",
				Games: []*domain.Game{
					{
						GameID:  1,
						Rows:    2,
						Columns: 2,
						Status:  constants.GameStatusOnGoing,
					},
				},
			},
			expectedStatus: constants.GameStatusOnGoing,
			userID:         "test_user_1",
		},
		{
			name: "OK/UPDATE",
			userGame: &domain.UserGame{
				UserID: "test_user_1",
				Games: []*domain.Game{
					{
						GameID:  2,
						Rows:    2,
						Columns: 2,
						Status:  constants.GameResultWon,
					},
				},
			},
			expectedStatus: constants.GameResultWon,
			userID:         "test_user_1",
		},
		{
			name: "FAIL/UPSERT_WRONG_USER",
			userGame: &domain.UserGame{
				UserID: "test_user_2",
				Games: []*domain.Game{
					{
						GameID:  1,
						Rows:    2,
						Columns: 2,
					},
				},
			},
			userID: "wrong_test_user",
		},
	}

	container := CreateInMemoryContainer()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			container.Upsert(c.userGame)
			gameFromContainer, _ := container.Get(c.userID)
			if c.name == "OK/INSERT" || c.name == "OK/UPDATE" {
				assert.Equal(t, gameFromContainer.UserID, c.userGame.UserID)
				assert.Equal(t, gameFromContainer.Games[0].GameID, c.userGame.Games[0].GameID)
				assert.Equal(t, gameFromContainer.Games[0].Status, c.expectedStatus)
			} else {
				assert.NotEqual(t, gameFromContainer, c.userGame)
			}
		})
	}
}
