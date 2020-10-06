package dao

import "github.com/mercadolibre/minesweeper/domain"

// InMemoryContainer ...
type InMemoryContainer struct {
	userGames []*domain.UserGame
}

// CreateInMemoryContainer ...
func CreateInMemoryContainer() *InMemoryContainer {

	container := &InMemoryContainer{
		userGames: make([]*domain.UserGame, 0),
	}
	return container
}

// Get ...
func (imc *InMemoryContainer) Get(userID string) (*domain.UserGame, error) {
	for _, userGame := range imc.userGames {
		if userGame.UserID == userID {
			return userGame, nil
		}
	}
	return nil, nil
}

// Upsert ...
func (imc *InMemoryContainer) Upsert(userGame *domain.UserGame) error {
	userFound := false
	for i, user := range imc.userGames {
		if user.UserID == userGame.UserID {
			imc.userGames[i].Games = userGame.Games
			userFound = true
		}
	}

	if !userFound {
		imc.userGames = append(imc.userGames, userGame)
	}
	return nil
}
