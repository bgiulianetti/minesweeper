package middlewares

import (
	"github.com/gin-gonic/gin"
)

// AdaptHandler ...
func AdaptHandler(handler func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {

		err := handler(c)

		if err != nil {
			c.Error(err)
			c.Abort()
		}
	}
}
