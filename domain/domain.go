package domain

import "time"

// Cell ..
type Cell struct {
	IsRevealed    bool   `json:"is_revealed" bson:"is_revealed"`
	HasMine       bool   `json:"has_mine" bson:"has_mine"`
	SourroundedBy int    `json:"sourrounded_by" bson:"sourrounded_by"`
	Flag          string `json:"flag" bson:"flag"`
}

// Game models tha minesweeper game properties
type Game struct {
	GameID        int64     `json:"game_id" bson:"game_id"`
	Rows          int       `json:"rows" bson:"rows"`
	Columns       int       `json:"columns" bson:"columns"`
	Mines         int       `json:"mines" bson:"mines"`
	Start         time.Time `json:"start_time" bson:"start_time"`
	Finish        time.Time `json:"finish_time" bson:"finish_time"`
	CellsRevealed int       `json:"cells_revealed" bson:"cells_revealed"`
	Status        string    `json:"status" bson:"status"`
	Board         [][]Cell  `json:"board" bson:"board"`
}

// UserGame models wich games owns wich user
type UserGame struct {
	Games  []*Game `json:"games" bson:"games"`
	UserID string  `json:"user_id" bson:"user_id"`
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
