package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GameController ...
type GameController struct{}

// Pong allows validation that the API is responding
func (gc GameController) Pong(c *gin.Context) {
	c.Set("skip", true)
	c.JSON(http.StatusOK, "Pong from minesweeper")
}
