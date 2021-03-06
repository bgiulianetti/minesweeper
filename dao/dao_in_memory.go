package dao

import (
	"sync"

	"github.com/mercadolibre/minesweeper/domain"
)

// InMemoryContainer ...
type InMemoryContainer struct {
	userGames []*domain.UserGame
	mutex     *sync.Mutex
}

// CreateInMemoryContainer initialize the container
func CreateInMemoryContainer() *InMemoryContainer {

	container := &InMemoryContainer{
		userGames: make([]*domain.UserGame, 0),
		mutex:     &sync.Mutex{},
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

// Update updates a game
func (imc *InMemoryContainer) Update(userGame *domain.UserGame) error {

	imc.mutex.Lock()
	defer imc.mutex.Unlock()

	for i, user := range imc.userGames {
		if user.UserID == userGame.UserID {
			imc.userGames[i].Games = userGame.Games
		}
	}
	return nil
}

// Insert inserts a new userGame
func (imc *InMemoryContainer) Insert(userGame *domain.UserGame) error {

	imc.mutex.Lock()
	defer imc.mutex.Unlock()

	imc.userGames = append(imc.userGames, userGame)
	return nil
}

// DeleteAll deletes all the games
func (imc *InMemoryContainer) DeleteAll() error {
	imc.mutex.Lock()
	defer imc.mutex.Unlock()

	imc.userGames = nil
	return nil
}
