package domain

import "time"

// Cell ..
type Cell struct {
	IsRevealed    bool   `json:"is_revealed"`
	HasMine       bool   `json:"has_mine"`
	SourroundedBy int    `json:"sourrounded_by"`
	Flag          string `json:"flag"`
}

// Game models tha minesweeper game properties
type Game struct {
	GameID        int64     `json:"game_id"`
	Rows          int       `json:"rows"`
	Columns       int       `json:"columns"`
	Mines         int       `json:"mines"`
	Start         time.Time `json:"start_time"`
	Finish        time.Time `json:"finish_time"`
	CellsRevealed int       `json:"cells_revealed"`
	Status        string    `json:"status"`
	Board         [][]Cell  `json:"board,omitempty"`
}

// UserGame models wich games owns wich user
type UserGame struct {
	Games  []*Game `json:"games"`
	UserID string  `json:"user_id"`
}

// NewGameConditionsRequest ...
type NewGameConditionsRequest struct {
	UserID  string `json:"user_id"`
	Rows    int    `json:"rows"`
	Columns int    `json:"columns"`
	Mines   int    `json:"mines"`
}

// FlagCellRequest ...
type FlagCellRequest struct {
	UserID string `json:"user_id"`
	GameID int64  `json:"game_id"`
	Row    int    `json:"row"`
	Column int    `json:"column"`
	Flag   string `json:"flag"`
}

// RevealCellRequest ...
type RevealCellRequest struct {
	UserID string `json:"user_id"`
	GameID int64  `json:"game_id"`
	Row    int    `json:"row"`
	Column int    `json:"column"`
}