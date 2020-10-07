package dao

import "github.com/mercadolibre/minesweeper/domain"

// InMemoryContainer ...
type InMemoryContainer struct {
	userGames []*domain.UserGame
}

// CreateInMemoryContainer initialize the container
func CreateInMemoryContainer() *InMemoryContainer {

	container := &InMemoryContainer{
		userGames: make([]*domain.UserGame, 0),
	}
	return container
}

// Get gets a game from a userID
func (imc *InMemoryContainer) Get(userID string) (*domain.UserGame, error) {
	for _, userGame := range imc.userGames {
		if userGame.UserID == userID {
			return userGame, nil
		}
	}
	return nil, nil
}

// GetAll gets all games from a userID
func (imc *InMemoryContainer) GetAll() ([]*domain.UserGame, error) {
	return imc.userGames, nil
}

// Upsert inserts or updates a game
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

// DeleteAll deletes all the games
func (imc *InMemoryContainer) DeleteAll() error {
	imc.userGames = nil
	return nil
}
